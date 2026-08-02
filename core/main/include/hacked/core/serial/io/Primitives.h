#pragma once

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

extern uint16_t serialReadU16LittleEndian(void const *addr);
extern int16_t serialReadS16LittleEndian(void const *addr);
extern uint32_t serialReadU32LittleEndian(void const *addr);
extern int32_t serialReadS32LittleEndian(void const *addr);

#ifdef __cplusplus
}
#endif
