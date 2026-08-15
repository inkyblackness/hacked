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
   validationResultAddConditional(&result, font->height != 0, "height is zero");
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

static void swapCodepointEntries(FontCodepointEntry *const entries, size_t const a, size_t const b)
{
   FontCodepointEntry temp = entries[a];
   entries[a] = entries[b];
   entries[b] = temp;
}

static size_t quicksortPartition(FontCodepointEntry *const entries, size_t const lo, size_t const hi)
{
   FontCodepointEntry const pivot = entries[hi];
   size_t i = lo;
   for (size_t j = lo; j < hi; j++)
   {
      if (entries[j].codepoint <= pivot.codepoint)
      {
         swapCodepointEntries(entries, i, j);
         i++;
      }
   }
   swapCodepointEntries(entries, i, hi);
   return i;
}

static void quicksortCodepoints(FontCodepointEntry *const entries, size_t const lo, size_t const hi)
{
   if (lo >= hi)
   {
      return;
   }
   size_t const partition = quicksortPartition(entries, lo, hi);
   if (partition > 0)
   {
      quicksortCodepoints(entries, lo, partition - 1);
   }
   quicksortCodepoints(entries, partition + 1, hi);
}

void fontSortCodepoints(FontCodepointEntry *const entries, size_t const count)
{
   if ((entries == NULL) || (count < 2))
   {
      return;
   }
   quicksortCodepoints(entries, 0, count - 1);
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
