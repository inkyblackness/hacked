#include <gtest/gtest.h>

#include "hacked/core/serial/io/Primitives.h"

namespace
{

TEST(PrimitivesTest, serialReadU16LittleEndian)
{
   uint8_t const data[2] = {0x12, 0x34};
   uint16_t const result = serialReadU16LittleEndian(data);
   EXPECT_EQ(result, 0x3412);
}

TEST(PrimitivesTest, serialReadS16LittleEndian)
{
   uint8_t const data[2] = {0x65, 0x87};
   int16_t const result = serialReadS16LittleEndian(data);
   EXPECT_EQ(result, -30875);
}

TEST(PrimitivesTest, serialReadU32LittleEndian)
{
   uint8_t const data[4] = {0x12, 0x34, 0x45, 0x67};
   uint32_t const result = serialReadU32LittleEndian(data);
   EXPECT_EQ(result, 0x67453412);
}

TEST(PrimitivesTest, serialReadS32LittleEndian)
{
   uint8_t const data[4] = {0x32, 0x54, 0x65, 0x87};
   int32_t const result = serialReadS32LittleEndian(data);
   EXPECT_EQ(result, -2023402446);
}

}
