#include "hacked/core/serial/io/Primitives.h"

uint16_t serialReadU16LittleEndian(void const *addr)
{
   uint8_t const *data = addr;
   return ((uint16_t)data[1] << 8) | (uint16_t)data[0];
}

int16_t serialReadS16LittleEndian(void const *addr)
{
   return (int16_t)serialReadU16LittleEndian(addr);
}

uint32_t serialReadU32LittleEndian(void const *addr)
{
   uint8_t const *data = addr;
   return ((uint32_t)data[3] << 24) | ((uint32_t)data[2] << 16) | ((uint32_t)data[1] << 8) | (uint32_t)data[0];
}

int32_t serialReadS32LittleEndian(void const *addr)
{
   return (int32_t)serialReadU32LittleEndian(addr);
}
