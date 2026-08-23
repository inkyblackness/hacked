#include <gtest/gtest.h>

#include "hacked/core/media/Text.h"

namespace
{
class TextTest : public ::testing::Test
{
protected:
   void TearDown() override
   {
      if (instance != NULL)
      {
         textDecoderRelease(&instance);
      }

      Test::TearDown();
   }

   TextDecoder *instance = NULL;
};

TEST_F(TextTest, textDecoderReleaseIgnoresNullInput)
{
   textDecoderRelease(NULL);
}

TEST_F(TextTest, textDecoderReleaseIgnoresNullVariableInput)
{
   textDecoderRelease(&instance);
}

TEST_F(TextTest, codepage437DecoderCanBeCreated)
{
   instance = textDecoderCreate(TEXT_CHARSET_CODEPAGE437);
   EXPECT_TRUE(instance != NULL);
   textDecoderRelease(&instance);
}

TEST_F(TextTest, codepage437DecoderReleaseSetsVariableToNull)
{
   instance = textDecoderCreate(TEXT_CHARSET_CODEPAGE437);
   EXPECT_TRUE(instance != NULL);
   textDecoderRelease(&instance);
   EXPECT_TRUE(instance == NULL);
}

TEST_F(TextTest, codepage437DecoderMaxOutputPerInputByte)
{
   instance = textDecoderCreate(TEXT_CHARSET_CODEPAGE437);
   EXPECT_TRUE(instance != NULL);
   size_t const result = textDecoderMaxOutputPerInputByte(instance);
   EXPECT_EQ(result, 1);
}

TEST_F(TextTest, codepage437DecoderReset)
{
   instance = textDecoderCreate(TEXT_CHARSET_CODEPAGE437);
   EXPECT_TRUE(instance != NULL);
   textDecoderReset(instance);
}

TEST_F(TextTest, codepage437FlushBehaviour)
{
   instance = textDecoderCreate(TEXT_CHARSET_CODEPAGE437);
   EXPECT_TRUE(instance != NULL);
   Codepoint out = 0;
   size_t outSize = 1;
   TextEncodingResult const result = textDecoderFlush(instance, &out, &outSize);
   EXPECT_EQ(outSize, 0);
   EXPECT_EQ(result, TEXT_ENCODING_DONE);
}

TEST_F(TextTest, codepage437Decode)
{
   instance = textDecoderCreate(TEXT_CHARSET_CODEPAGE437);
   EXPECT_TRUE(instance != NULL);
   uint8_t const input[5] = {0x31, 0x32, 0x33, 144, 225};
   size_t inSize = sizeof(input);
   Codepoint output[5] = {};
   size_t outSize = sizeof(output) / sizeof(output[0]);
   TextEncodingResult const result = textDecoderDecode(instance, output, &outSize, input, &inSize);
   EXPECT_EQ(result, TEXT_ENCODING_DONE);
   EXPECT_EQ(outSize, 5);
   EXPECT_EQ(inSize, 5);
   Codepoint const expected[5] = {0x31, 0x32, 0x33, 201, 223};
   for (size_t i = 0; i < 5; i++)
   {
      EXPECT_EQ(output[i], expected[i]) << "failed for i=" << std::to_string(i);
   }
}

TEST_F(TextTest, codepage437DecodeNeedsMoreOutput)
{
   instance = textDecoderCreate(TEXT_CHARSET_CODEPAGE437);
   EXPECT_TRUE(instance != NULL);
   uint8_t const input[5] = {0x31, 0x32, 0x33, 144, 225};
   size_t inSize = sizeof(input);
   Codepoint output[2] = {};
   size_t outSize = sizeof(output) / sizeof(output[0]);
   TextEncodingResult const result = textDecoderDecode(instance, output, &outSize, input, &inSize);
   EXPECT_EQ(result, TEXT_ENCODING_MORE_OUTPUT_NEEDED);
   EXPECT_EQ(outSize, 2);
   EXPECT_EQ(inSize, 2);
   Codepoint const expected[2] = {0x31, 0x32};
   for (size_t i = 0; i < 2; i++)
   {
      EXPECT_EQ(output[i], expected[i]) << "failed for i=" << std::to_string(i);
   }
}

}
