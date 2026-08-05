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

TEST(FontTest, fontFindCodepointReturnsEntry)
{
   static size_t constexpr codepointCount = 20;
   Font instance = {};
   FontCodepointEntry entries[codepointCount] = {};
   for (size_t i = 0; i < codepointCount; i++)
   {
      entries[i].codepoint = U'A' + i;
   }
   instance.codepointCount = codepointCount;
   instance.codepoints = entries;

   for (auto const &expected : entries)
   {
      EXPECT_TRUE(fontFindCodepointEntry(&instance, expected.codepoint) == &expected) << "failed for " + std::to_string(expected.codepoint);
   }
}
