#pragma once

#include "hacked/core/media/Font.h"

#ifdef __cplusplus
extern "C" {
#endif

[[nodiscard]] extern Font const *lgresDecodeFont(uint8_t const *data, size_t dataSize);

#ifdef __cplusplus
}
#endif
