#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>

[[nodiscard]] extern uint16_t serialReadU16LittleEndian(void const *addr);
[[nodiscard]] extern int16_t serialReadS16LittleEndian(void const *addr);
[[nodiscard]] extern uint32_t serialReadU32LittleEndian(void const *addr);
[[nodiscard]] extern int32_t serialReadS32LittleEndian(void const *addr);

#ifdef __cplusplus
}
#endif
