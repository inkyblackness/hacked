#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include <stddef.h>
#include <stdint.h>

#include "hacked/core/media/PixelSpace.h"
#include "hacked/core/system/Validation.h"

typedef struct
{
   PixelSize size;

   uint8_t *data;
   size_t dataLength;
   size_t stride;
} Bitmap;

[[nodiscard]] extern ValidationResult bitmapValidate(Bitmap const *bitmap);

#ifdef __cplusplus
}
#endif
