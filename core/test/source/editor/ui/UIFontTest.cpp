#include <gtest/gtest.h>

#include "hacked/core/editor/ui/UIFont.h"

namespace
{
TEST(UIFontTest, fontIsAvailable)
{
   Font const *font = uiFont();
   EXPECT_TRUE(font != nullptr);
#if 0
   if (font != nullptr)
   {
      for (size_t y = 0; y < font->atlas.size.height; y++)
      {
         for (size_t x = 0; x < font->atlas.size.width; x++)
         {
            uint8_t in = font->atlas.data[font->atlas.stride * y + (x)];
            printf("%c", in ? '#' : '.');
         }
         printf("\n");
      }

      for (size_t i = 0; i < font->codepointCount; i++)
      {
         printf("====== 0x%08X -> '%lc'\n", font->codepoints[i].codepoint, (wint_t)font->codepoints[i].codepoint);
         for (size_t y = 0; y < font->codepoints[i].rect.size.height; y++)
         {
            for (size_t x = 0; x < font->codepoints[i].rect.size.width; x++)
            {
               uint8_t in = font->atlas.data[font->atlas.stride * (font->codepoints[i].rect.topLeft.y + y) + (font->codepoints[i].rect.topLeft.x + x)];
               printf("%c", in ? '#' : '.');
            }
            printf("\n");
         }
      }
   }
#endif
   fontRelease(&font);
}

}