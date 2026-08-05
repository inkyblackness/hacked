#pragma once

#include <stddef.h>
#include <stdint.h>

#include "hacked/core/media/PixelSpace.h"

#ifdef __cplusplus
extern "C" {
#endif

typedef struct
{
   PixelSize size;

   uint8_t *data;
   size_t dataLength;
   size_t stride;
} Bitmap;

#ifdef __cplusplus
}
#endif
