#include "hacked/core/media/Font.h"

void fontRelease(Font const **const fontRef)
{
   if ((fontRef == NULL) || (*fontRef == NULL))
   {
      return;
   }
   Font const *font = *fontRef;
   *fontRef = NULL;
   if (font->release == NULL)
   {
      return;
   }
   font->release(font);
}

ValidationResult fontValidate(Font const *const font)
{
   ValidationResult result = {0};
   if (font == NULL)
   {
      return validationResultAddMessage(&result, "font is NULL");
   }
   validationResultMerge(&result, bitmapValidate(&font->atlas), "font atlas bitmap has issues");
   validationResultAddConditional(&result, font->codepoints != NULL, "codepoints is NULL");
   validationResultAddConditional(&result, font->codepointCount != 0, "codepointCount is zero");
   if ((font->codepoints != NULL) && (font->codepointCount > 0))
   {
      Codepoint lastCodepoint = font->codepoints[0].codepoint;
      PixelRect bitmapRect = {.size = font->atlas.size, .topLeft = {.x = 0, .y = 0}};
      bool sorted = true;
      bool unique = true;
      bool inside = true;
      bool properSize = true;
      for (size_t i = 0; i < font->codepointCount; i++)
      {
         FontCodepointEntry const *const entry = &font->codepoints[i];
         if (i > 0)
         {
            if (entry->codepoint < lastCodepoint)
            {
               sorted = false;
            }
            if (entry->codepoint == lastCodepoint)
            {
               unique = false;
            }
         }
         if ((entry->rect.size.width == 0) || (entry->rect.size.height == 0))
         {
            // An argument could be made that instead pixelSpaceValidateSize is used. Maybe it should.
            properSize = false;
         }
         if (!pixelSpaceAreaContainsRect(bitmapRect, entry->rect))
         {
            inside = false;
         }
         lastCodepoint = entry->codepoint;
      }
      validationResultAddConditional(&result, sorted, "codepoints not sorted");
      validationResultAddConditional(&result, unique, "codepoints not unique");
      validationResultAddConditional(&result, inside, "codepoints outside atlas");
      validationResultAddConditional(&result, properSize, "codepoints with zero size");
   }
   return result;
}

FontCodepointEntry const *fontFindCodepointEntry(Font const *const font, Codepoint const codepoint)
{
   size_t begin = 0;
   size_t end = font->codepointCount;
   while (begin < end)
   {
      size_t const middle = begin + ((end - begin) / 2);
      FontCodepointEntry const *const entry = &font->codepoints[middle];
      if (entry->codepoint < codepoint)
      {
         begin = middle + 1;
      }
      else if (entry->codepoint > codepoint)
      {
         end = middle;
      }
      else
      {
         return entry;
      }
   }
   return NULL;
}
