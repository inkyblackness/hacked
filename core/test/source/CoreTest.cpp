#include <gtest/gtest.h>

#include "hacked/core/Core.h"

TEST(CoreTest, sample)
{
   EXPECT_EQ(12, sampleReturn(2)) << "failure";
}
