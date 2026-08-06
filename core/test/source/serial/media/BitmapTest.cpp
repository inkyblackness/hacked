#include <gtest/gtest.h>

#include "hacked/core/media/Bitmap.h"

class BitmapTest : public ::testing::Test
{
protected:
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

   static void assertResultMessages(ValidationResult const &result, std::vector<std::string> const &messages)
   {
      for (size_t i = 0; i < messages.size(); i++)
      {
         if (result.messages[i] != nullptr)
         {
            EXPECT_EQ(std::string(result.messages[i]), messages.at(i)) << "failed at " << std::to_string(i);
         }
         else
         {
            EXPECT_TRUE(false) << "NULL at " << std::to_string(i) << ", expecting '" << messages[i] << "'";
         }
      }
      for (size_t i = messages.size(); i < VALIDATION_RESULT_MAX_MESSAGES; i++)
      {
         EXPECT_TRUE(result.messages[i] == nullptr) << "tail not null at " << std::to_string(i) << ", has '" << result.messages[i] << "'";
      }
   }
};

TEST_F(BitmapTest, validateOk)
{
   Bitmap const bitmap = validBitmap();
   ValidationResult const result = bitmapValidate(&bitmap);
   EXPECT_FALSE(validationResultHasFailure(&result));
}

TEST_F(BitmapTest, validatePointerNull)
{
   ValidationResult const result = bitmapValidate(nullptr);
   assertResultMessages(result, {"bitmap is NULL"});
}

TEST_F(BitmapTest, validateDataNull)
{
   Bitmap bitmap = validBitmap();
   bitmap.data = nullptr;
   ValidationResult const result = bitmapValidate(&bitmap);
   assertResultMessages(result, {"data is NULL"});
}

TEST_F(BitmapTest, validateEmptySize)
{
   Bitmap bitmap = validBitmap();
   bitmap.size.width = 0;
   bitmap.size.height = 0;
   ValidationResult const result = bitmapValidate(&bitmap);
   assertResultMessages(result, {"width is zero", "height is zero"});
}

TEST_F(BitmapTest, validateWidthAgainstStride)
{
   Bitmap bitmap = validBitmap();
   bitmap.stride = bitmap.size.width - 1;
   ValidationResult const result = bitmapValidate(&bitmap);
   assertResultMessages(result, {"width is higher than stride"});
}

TEST_F(BitmapTest, validateHeightTimesStrideAgainstDataLength)
{
   Bitmap bitmap = validBitmap();
   bitmap.size.height++;
   ValidationResult const result = bitmapValidate(&bitmap);
   assertResultMessages(result, {"stride times height is more than dataLength"});
}

TEST_F(BitmapTest, validateStrideZero)
{
   // Although this case is covered by width <= stride, to know that stride itself is wrong can help.
   Bitmap bitmap = validBitmap();
   bitmap.stride = 0;
   ValidationResult const result = bitmapValidate(&bitmap);
   assertResultMessages(result, {"stride is zero", "width is higher than stride"});
}
