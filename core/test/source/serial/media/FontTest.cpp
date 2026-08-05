#include <gtest/gtest.h>

#include "hacked/core/media/Font.h"

TEST(FontTest, fontReleaseIgnoresNullArgument)
{
   fontRelease(nullptr);
}

TEST(FontTest, fontReleaseIgnoresNullPointer)
{
   Font *font = nullptr;
   fontRelease(&font);
}

static void testReleaseFont(Font *const font)
{
   ASSERT_TRUE(font != nullptr) << "font not provided";
   ASSERT_TRUE(font->userData != nullptr) << "userdata not provided";
   auto const called = static_cast<bool *>(font->userData);
   *called = true;
}

TEST(FontTest, fontReleaseCallsUserFunction)
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

TEST(FontTest, fontReleaseResetsPointerEvenWithoutRelease)
{
   Font instance = {};
   Font *font = &instance;
   fontRelease(&font);
   EXPECT_TRUE(font == nullptr) << "font should be nullptr";
}

TEST(FontTest, fontFindCodepointReturnsNullForEmptyList)
{
   Font constexpr instance = {};
   auto const result = fontFindCodepointEntry(&instance, U'A');
   EXPECT_TRUE(result == nullptr);
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

TEST(FontTest, fontFindCodepointReturnsEntryEvenNumberOfCodepoints)
{
   verifyFontFindCodepointReturnsEntry(true);
}

TEST(FontTest, fontFindCodepointReturnsEntryOddNumberOfCodepoints)
{
   verifyFontFindCodepointReturnsEntry(false);
}
