#include <gtest/gtest.h>

#include "hacked/core/media/Font.h"

#include "hacked/core/test/ValidationAsserts.h"

namespace
{

class FontTest : public ::testing::Test
{
protected:
   static void testReleaseFont(Font *const font)
   {
      ASSERT_TRUE(font != nullptr) << "font not provided";
      ASSERT_TRUE(font->userData != nullptr) << "userdata not provided";
      auto const called = static_cast<bool *>(font->userData);
      *called = true;
   }

   static void verifyFontFindCodepointReturnsEntry(bool const evenNumberOfCodepoints)
   {
      static size_t constexpr codepointCount = 10;
      FontCodepointEntry entries[codepointCount] = {};
      entries[0].codepoint = U'B';
      entries[1].codepoint = U'C';
      entries[2].codepoint = U'E';
      entries[3].codepoint = U'G';
      entries[4].codepoint = U'H';
      entries[5].codepoint = U'I';
      entries[6].codepoint = U'J';
      entries[7].codepoint = U'O';
      entries[8].codepoint = U'P';
      entries[9].codepoint = U'S';
      Font instance = {};
      instance.codepointCount = codepointCount - (evenNumberOfCodepoints ? 0 : 1);
      instance.codepoints = entries;

      for (size_t i = 0; i < instance.codepointCount; i++)
      {
         FontCodepointEntry const *const expected = &entries[i];
         EXPECT_TRUE(fontFindCodepointEntry(&instance, expected->codepoint) == expected) << "failed for " + std::to_string(expected->codepoint);
      }
      EXPECT_TRUE(fontFindCodepointEntry(&instance, U'A') == nullptr) << "should not find unknown codepoint before list";
      EXPECT_TRUE(fontFindCodepointEntry(&instance, U'D') == nullptr) << "should not find unknown codepoint within list 1";
      EXPECT_TRUE(fontFindCodepointEntry(&instance, U'K') == nullptr) << "should not find unknown codepoint within list 2";
      EXPECT_TRUE(fontFindCodepointEntry(&instance, U'L') == nullptr) << "should not find unknown codepoint within list 3";
      EXPECT_TRUE(fontFindCodepointEntry(&instance, U'Z') == nullptr) << "should not find unknown codepoint beyond list";
   }

   static Font validFont()
   {
      static size_t constexpr codepointCount = 3;
      static FontCodepointEntry codepoints[codepointCount] = {};
      for (size_t i = 0; i < codepointCount; i++)
      {
         auto &codepoint = codepoints[i];
         codepoint.codepoint = 'A' + i;
         codepoint.rect.size.width = 2;
         codepoint.rect.size.height = 2;
      }
      Font font = {};
      font.bitmap = validBitmap();
      font.codepointCount = codepointCount;
      font.codepoints = codepoints;
      return font;
   }

private:
   static Bitmap validBitmap()
   {
      static size_t constexpr stride = 4;
      static PixelAxisSize constexpr width = 3;
      static PixelAxisSize constexpr height = 2;
      static size_t constexpr dataLength = stride * height;
      static uint8_t data[dataLength] = {};
      Bitmap bitmap = {};
      bitmap.size.width = width;
      bitmap.size.height = height;
      bitmap.data = data;
      bitmap.dataLength = dataLength;
      bitmap.stride = stride;
      return bitmap;
   }
};

TEST_F(FontTest, fontReleaseIgnoresNullArgument)
{
   fontRelease(nullptr);
}

TEST_F(FontTest, fontReleaseIgnoresNullPointer)
{
   Font *font = nullptr;
   fontRelease(&font);
}

TEST_F(FontTest, fontReleaseCallsUserFunction)
{
   bool called = false;
   Font instance = {};
   instance.userData = &called;
   instance.release = testReleaseFont;
   Font *font = &instance;
   fontRelease(&font);
   EXPECT_TRUE(called) << "release was not called";
   EXPECT_TRUE(font == nullptr) << "font should be nullptr";
}

TEST_F(FontTest, fontReleaseResetsPointerEvenWithoutRelease)
{
   Font instance = {};
   Font *font = &instance;
   fontRelease(&font);
   EXPECT_TRUE(font == nullptr) << "font should be nullptr";
}

TEST_F(FontTest, validatePointerNull)
{
   ValidationResult const result = fontValidate(nullptr);
   assertResultMessages(result, {"font is NULL"});
}

TEST_F(FontTest, validateOk)
{
   Font const font = validFont();
   ValidationResult const result = fontValidate(&font);
   assertResultMessages(result, {});
}

TEST_F(FontTest, validateCodepointsNull)
{
   Font font = validFont();
   font.codepoints = nullptr;
   ValidationResult const result = fontValidate(&font);
   assertResultMessages(result, {"codepoints is NULL"});
}

TEST_F(FontTest, validateCodepointCountZero)
{
   Font font = validFont();
   font.codepointCount = 0;
   ValidationResult const result = fontValidate(&font);
   assertResultMessages(result, {"codepointCount is zero"});
}

TEST_F(FontTest, validateCodepointsNotSorted)
{
   Font font = validFont();
   font.codepoints[0].codepoint += font.codepointCount;
   ValidationResult const result = fontValidate(&font);
   assertResultMessages(result, {"codepoints not sorted"});
}

TEST_F(FontTest, validateCodepointsUnique)
{
   Font font = validFont();
   font.codepoints[0].codepoint = font.codepoints[1].codepoint;
   ValidationResult const result = fontValidate(&font);
   assertResultMessages(result, {"codepoints not unique"});
}

TEST_F(FontTest, validateCodepointsZeroSize)
{
   Font font = validFont();
   font.codepoints[0].rect.size.width = 0;
   ValidationResult const result = fontValidate(&font);
   assertResultMessages(result, {"codepoints with zero size"});
}

TEST_F(FontTest, validateCodepointsOutsideBitmap)
{
   Font font = validFont();
   font.codepoints[0].rect.topLeft.x = font.bitmap.size.width;
   ValidationResult const result = fontValidate(&font);
   assertResultMessages(result, {"codepoints outside bitmap"});
}

TEST_F(FontTest, validateCodepointsTooLargeForBitmap)
{
   Font font = validFont();
   font.codepoints[0].rect.topLeft.y++;
   font.codepoints[0].rect.size.height = font.bitmap.size.height;
   ValidationResult const result = fontValidate(&font);
   assertResultMessages(result, {"codepoints outside bitmap"});
}

TEST_F(FontTest, fontFindCodepointReturnsNullForEmptyList)
{
   Font constexpr instance = {};
   auto const result = fontFindCodepointEntry(&instance, U'A');
   EXPECT_TRUE(result == nullptr);
}

TEST_F(FontTest, fontFindCodepointReturnsEntryEvenNumberOfCodepoints)
{
   verifyFontFindCodepointReturnsEntry(true);
}

TEST_F(FontTest, fontFindCodepointReturnsEntryOddNumberOfCodepoints)
{
   verifyFontFindCodepointReturnsEntry(false);
}

}
