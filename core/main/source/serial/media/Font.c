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
   for (size_t i = 0; i < font->codepointCount; ++i)
   {
      // TODO: implement binary search
      FontCodepointEntry const *const entry = &font->codepoints[i];
      if (entry->codepoint == codepoint)
      {
         return entry;
      }
   }
   return NULL;
}
