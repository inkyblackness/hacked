#include <gtest/gtest.h>

#include "hacked/core/media/Bitmap.h"

#include "hacked/core/test/ValidationAsserts.h"

namespace
{

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

TEST_F(BitmapTest, validateSizeBelowLimit)
{
   Bitmap bitmap = validBitmap();
   bitmap.size.width = 0;
   bitmap.size.height = 0;
   ValidationResult const result = bitmapValidate(&bitmap);
   assertResultMessages(result, {"width below limit", "height below limit"});
}

TEST_F(BitmapTest, validateSizeAboveLimit)
{
   PixelSpaceLimits const limits = pixelSpaceLimits();
   Bitmap bitmap = validBitmap();
   bitmap.size.width = limits.maxSize + 1;
   bitmap.size.height = limits.maxSize + 1;
   bitmap.stride = bitmap.size.width;
   bitmap.dataLength = bitmap.stride * static_cast<size_t>(bitmap.size.height);
   ValidationResult const result = bitmapValidate(&bitmap);
   assertResultMessages(result, {"width above limit", "height above limit"});
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

}
