#include "hacked/core/media/Font.h"

void fontRelease(Font **fontRef)
{
   if ((fontRef == NULL) || (*fontRef == NULL))
   {
      return;
   }
   Font *font = *fontRef;
   *fontRef = NULL;
   if (font->release == NULL)
   {
      return;
   }
   font->release(font);
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
