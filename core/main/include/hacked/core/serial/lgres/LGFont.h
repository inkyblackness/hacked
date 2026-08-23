#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include "hacked/core/media/Font.h"

[[nodiscard]] extern Font const *lgresDecodeFont(uint8_t const *data, size_t dataSize, PixelAxisSize glyphPadding);

#ifdef __cplusplus
}
#endif
