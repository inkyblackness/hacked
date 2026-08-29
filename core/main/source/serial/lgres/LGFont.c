#include <stdlib.h>

#include "hacked/core/media/Text.h"
#include "hacked/core/serial/io/Primitives.h"
#include "hacked/core/serial/lgres/LGFont.h"

typedef enum LGFontHeaderConstants
{
   HEADER_FONT_TYPE = 0x0000,
   HEADER_FIRST_CHARACTER_INDEX = 0x0024,
   HEADER_LAST_CHARACTER_INDEX = 0x0026,
   HEADER_X_OFFSETS = 0x0048,
   HEADER_BITMAP = 0x004C,
   HEADER_WIDTH = 0x0050,
   HEADER_HEIGHT = 0x0052,
   HEADER_SIZE = 0x0054,

   HEADER_FONT_TYPE_COLOR = 0xCCCC
} LGFontHeaderConstants;

typedef struct LGFont LGFont;

typedef uint8_t (*ReadFontPixelFunc)(LGFont const *lgFont, size_t index, int16_t x, int16_t y);

typedef struct LGFont
{
   bool color;
   ReadFontPixelFunc readPixel;
   int16_t firstCharacterIndex;
   int16_t lastCharacterIndex;
   int16_t const *xOffsets;
   // Bitmap information is deliberately not a media/Bitmap, since the size parameters could be outside limits.
   // LGRes has the font serialized in a single line.
   uint8_t const *bitmapData;
   int16_t bitmapWidth;
   int16_t bitmapHeight;
} LGFont;

[[nodiscard]] static uint8_t readFontPixelColor(LGFont const *const lgFont, size_t const index, int16_t const x, int16_t const y)
{
   return lgFont->bitmapData[(size_t)lgFont->bitmapWidth * (size_t)y + (size_t)lgFont->xOffsets[index] + (size_t)x];
}

[[nodiscard]] static uint8_t readFontPixelMonochrome(LGFont const *const lgFont, size_t const index, int16_t const x, int16_t const y)
{
   size_t const byteOffset = (lgFont->xOffsets[index] + x) / 8;
   uint8_t const b = lgFont->bitmapData[(size_t)lgFont->bitmapWidth * (size_t)y + byteOffset];
   size_t const bitOffset = (lgFont->xOffsets[index] + x) % 8;
   return ((b & (0x80 >> bitOffset)) != 0) ? 0x01 : 0x00;
}

static bool readHeader(LGFont *lgFont, uint8_t const *const data, size_t const dataSize)
{
   if (dataSize < HEADER_SIZE)
   {
      return false;
   }
   lgFont->color = serialReadU16LittleEndian(data + HEADER_FONT_TYPE) != 0;
   lgFont->readPixel = lgFont->color ? readFontPixelColor : readFontPixelMonochrome;
   lgFont->firstCharacterIndex = serialReadS16LittleEndian(data + HEADER_FIRST_CHARACTER_INDEX);
   lgFont->lastCharacterIndex = serialReadS16LittleEndian(data + HEADER_LAST_CHARACTER_INDEX);
   if ((lgFont->firstCharacterIndex < 0) || (lgFont->lastCharacterIndex < lgFont->firstCharacterIndex))
   {
      return false;
   }
   size_t const xOffsets = serialReadU32LittleEndian(data + HEADER_X_OFFSETS);
   if (xOffsets >= dataSize)
   {
      return false;
   }
   lgFont->xOffsets = (int16_t *)(data + xOffsets);
   // TODO: check if all offsets are non-negative and their size is within pixel space limits
   size_t const bitmap = serialReadU32LittleEndian(data + HEADER_BITMAP);
   lgFont->bitmapWidth = serialReadS16LittleEndian(data + HEADER_WIDTH);
   lgFont->bitmapHeight = serialReadS16LittleEndian(data + HEADER_HEIGHT);
   if ((lgFont->bitmapWidth <= 0) || (lgFont->bitmapHeight <= 0) || (bitmap >= dataSize))
   {
      return false;
   }
   size_t bitmapSize = (size_t)lgFont->bitmapWidth * (size_t)lgFont->bitmapHeight;
   if ((bitmap + bitmapSize) > dataSize)
   {
      return false;
   }
   lgFont->bitmapData = data + bitmap;
   return true;
}

[[nodiscard]] static FontCodepointEntry *allocateCodepoints(LGFont const *const lgFont, size_t *codepointCount)
{
   *codepointCount = lgFont->lastCharacterIndex - lgFont->firstCharacterIndex + 1;
   FontCodepointEntry *const codepoints = calloc(*codepointCount, sizeof(FontCodepointEntry));
   TextDecoder *decoder = textDecoderCreate(TEXT_CHARSET_CODEPAGE437);
   for (size_t i = 0; i < *codepointCount; i++)
   {
      size_t outSize = 1;
      uint8_t in = (uint8_t)(lgFont->firstCharacterIndex + i);
      size_t inSize = 1;
      (void)textDecoderDecode(decoder, &codepoints[i].codepoint, &outSize, &in, &inSize);
   }
   textDecoderRelease(&decoder);
   return codepoints;
}

[[nodiscard]] static bool isEmptyFontGlyph(LGFont const *const lgFont, size_t const index)
{
   for (int16_t y = 0; y < lgFont->bitmapHeight; y++)
   {
      if (lgFont->readPixel(lgFont, index, 0, y) != 0x00)
      {
         return false;
      }
   }
   return true;
}

static PixelAxisSize const emptyMarkerHeight = 1;

[[nodiscard]] static bool attemptAtlasLayoutForWidth(LGFont const *const lgFont, FontCodepointEntry *const codepoints, size_t const codepointCount,
   PixelAxisSize const width, PixelAxisSize *const height, PixelAxisSize glyphPadding)
{
   PixelAxisSize const fullGlyphPadding = glyphPadding * 2; // apply padding on both sides
   PixelAxisOffset x = 1 + fullGlyphPadding; // start with one to cover for the empty codepoint.
   PixelAxisOffset y = glyphPadding;
   *height = lgFont->bitmapHeight + fullGlyphPadding;
   for (size_t i = 0; i < codepointCount; i++)
   {
      if (codepoints[i].rect.size.height == emptyMarkerHeight)
      {
         // skip empty codepoint entries, they are placed at the very beginning, automatically.
         continue;
      }
      PixelAxisOffset const right = x + codepoints[i].rect.size.width + fullGlyphPadding;
      if (right <= width)
      {
         codepoints[i].rect.topLeft.x = (PixelAxisPosition)(x + glyphPadding);
         codepoints[i].rect.topLeft.y = (PixelAxisPosition)y;
         x = right;
      }
      else if ((*height + lgFont->bitmapHeight) <= width)
      {
         y += lgFont->bitmapHeight + fullGlyphPadding;
         codepoints[i].rect.topLeft.x = (PixelAxisPosition)glyphPadding;
         codepoints[i].rect.topLeft.y = (PixelAxisPosition)y;
         x = codepoints[i].rect.size.width + glyphPadding;
         *height = y + lgFont->bitmapHeight + glyphPadding;
      }
      else
      {
         return false;
      }
   }
   return true;
}

[[nodiscard]] static bool determineAtlasSize(
   LGFont const *const lgFont, FontCodepointEntry *const codepoints, size_t const codepointCount, PixelSize *const size, PixelAxisSize const glyphPadding)
{
   PixelSpaceLimits const limits = pixelSpaceLimits();
   PixelAxisSize left = limits.minSize;
   PixelAxisSize right = limits.maxSize + 1;
   size->width = 0;
   while (left <= right)
   {
      PixelSize attempt = {.width = left + (right - left) / 2};
      if (attemptAtlasLayoutForWidth(lgFont, codepoints, codepointCount, attempt.width, &attempt.height, glyphPadding))
      {
         *size = attempt;
         right = attempt.width;
      }
      else if ((attempt.width > left) && (attempt.width < right))
      {
         left = attempt.width;
      }
      else
      {
         break;
      }
   }
   // recalculate all positions again, because a previous attempt might have modified the positions of the codepoints.
   return attemptAtlasLayoutForWidth(lgFont, codepoints, codepointCount, size->width, &size->height, glyphPadding);
}

/*
 * This function figures out how to semi-optimally create an atlas bitmap that contains all the glyphs.
 * Although LGRes files store the font in a single line, this code here attempts to create a square bitmap with multiple rows.
 * Main reason for this is to avoid exceeding the width of the pixel space limits.
 * There are probably several further optimizations that could be done, yet this function keeps it somewhat simple:
 * For any "empty" entry, re-use the same space in the bitmap. Use a binary-search approach to determine a somewhat square bitmap.
 * There are not many characters to consider, so a brute-force approach is feasible. A naive approach could have been taking
 * the number of characters and square-root this number.
 */
[[nodiscard]] static bool determineAtlasLayout(
   LGFont const *const lgFont, FontCodepointEntry *const codepoints, size_t const codepointCount, PixelSize *const size, PixelAxisSize const glyphPadding)
{
   PixelRect const emptyCodepointRect = {.size = {.width = 1, .height = lgFont->bitmapHeight}, .topLeft = {.x = 0, .y = 0}};
   // first pass: initialize target with per codepoint, and determine whether it would be an "empty" codepoint.
   // Mark the codepoints that are empty, using the yet unused height field.
   for (size_t i = 0; i < codepointCount; i++)
   {
      size_t const currentOffset = lgFont->xOffsets[i];
      size_t const nextOffset = lgFont->xOffsets[i + 1];
      PixelAxisSize const width = (PixelAxisSize)(nextOffset - currentOffset);
      if ((width == 0) || ((width == 1) && isEmptyFontGlyph(lgFont, i)))
      {
         // There are several characters in the original files that have only a single column of empty pixels.
         // There is no known case of a zero-width character, yet handle it the same. If that were to exist, it would be exported differently.
         codepoints[i].rect.size.width = emptyCodepointRect.size.width;
         codepoints[i].rect.size.height = emptyMarkerHeight;
      }
      else
      {
         codepoints[i].rect.size.width = width;
      }
   }
   // now determine whether any size of the atlas will make room for all characters.
   if (!determineAtlasSize(lgFont, codepoints, codepointCount, size, glyphPadding))
   {
      return false;
   }

   // last pass: set the height of all codepoints to that of the bitmap.
   for (size_t i = 0; i < codepointCount; i++)
   {
      codepoints[i].rect.size.height = lgFont->bitmapHeight;
   }
   return true;
}

static void copyFontPixel(Font const *const font, LGFont const *const lgFont)
{
   for (size_t i = 0; i < font->codepointCount; i++)
   {
      for (int16_t y = 0; y < font->codepoints[i].rect.size.height; y++)
      {
         uint8_t *out = font->atlas.data + (font->atlas.stride * ((size_t)font->codepoints[i].rect.topLeft.y + (size_t)y)) + font->codepoints[i].rect.topLeft.x;
         for (int16_t x = 0; x < font->codepoints[i].rect.size.width; x++)
         {
            *out = lgFont->readPixel(lgFont, i, x, y);
            out++;
         }
      }
   }
}

static void freeFont(Font const *font)
{
   free(font->codepoints);
   free(font->userData);
}

static Font *allocateFont(FontCodepointEntry *const codepoints, size_t const codepointCount, PixelSize const size)
{
   static size_t const alignment = 32;
   Bitmap bitmap = {0};
   bitmap.size = size;
   bitmap.stride = (((size_t)bitmap.size.width + (alignment - 1)) / alignment) * alignment;
   bitmap.dataLength = bitmap.stride * (size_t)bitmap.size.height;
   void *const memory = calloc(bitmap.dataLength + (alignment - 1) + sizeof(Font), 1);
   bitmap.data = (uint8_t *)((((uintptr_t)(memory) + (alignment - 1)) / alignment) * alignment);
   Font *const font = (Font *)(bitmap.data + bitmap.dataLength);
   font->userData = memory;
   font->release = freeFont;
   font->codepoints = codepoints;
   font->codepointCount = codepointCount;
   font->atlas = bitmap;
   return font;
}

[[nodiscard]] Font const *lgresDecodeFont(uint8_t const *const data, size_t const dataSize, PixelAxisSize const glyphPadding)
{
   LGFont lgFont = {0};
   if (!readHeader(&lgFont, data, dataSize))
   {
      return NULL;
   }
   size_t codepointCount = 0;
   FontCodepointEntry *const codepoints = allocateCodepoints(&lgFont, &codepointCount);
   PixelSize bitmapSize = {0};
   if (!determineAtlasLayout(&lgFont, codepoints, codepointCount, &bitmapSize, glyphPadding))
   {
      free(codepoints);
      return NULL;
   }
   Font *const font = allocateFont(codepoints, codepointCount, bitmapSize);
   font->color = lgFont.color;
   font->height = lgFont.bitmapHeight;
   copyFontPixel(font, &lgFont);
   fontSortCodepoints(font->codepoints, font->codepointCount);
   return font;
}
