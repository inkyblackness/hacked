#include <gtest/gtest.h>

#include "hacked/core/media/PixelSpace.h"

TEST(PixelSpaceTest, offsetDataTypeCanHoldPositionDataType)
{
   EXPECT_GE(sizeof(PixelAxisOffset), sizeof(PixelAxisPosition));
}

TEST(PixelSpaceTest, limitsAreFeasible)
{
   PixelSpaceLimits limits = pixelSpaceLimits();
   // Size is the base property; It needs to be able to hold at least
   // high-res cutscenes of the original game, which were 600x300 in size, taking the higher value.
   EXPECT_EQ(limits.minSize, 1);
   EXPECT_GE(limits.maxSize, 600) << "size can not hold largest known bitmap of original game";
   // Positions need to be able to address any pixel within size.
   // Furthermore, in case a bitmap is to be shifted, it should be possible to shift it up to
   // being not visible (barely outside of typically reachable space).
   EXPECT_LE(limits.minPosition, -limits.maxSize);
   EXPECT_GE(limits.maxPosition, limits.maxSize);
   // Offsets need to allow for accumulating addition/subtraction, so they need to be higher.
   EXPECT_LT(limits.minOffset, limits.minPosition);
   EXPECT_GT(limits.maxOffset, limits.maxPosition);
}
