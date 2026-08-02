#pragma once

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

/**
 * Position is an absolute location based on an origin (0).
 */
typedef int16_t PixelAxisPosition;
/**
 * Offset is the difference between two positions.
 * It is one value type larger than the position type, to allow for easier checks regarding overflow.
 */
typedef int32_t PixelAxisOffset;
/**
 * Size describes the extent of a region, which must be non-negative.
 */
typedef uint16_t PixelAxisSize;

/**
 * Pixel space is defined and based on the assumption that bitmaps and their presentation
 * is in a rather small space. Because the majority of fields regarding pixels are 16bit
 * integers in the original engine, that is reflected here as well.
 * Given that the largest texture in the original game has a size of 128x128 pixels,
 * only slightly topped by the high-res cutscene movies, the actually needed highest number
 * is still less than 1000.
 *
 * Use the limits expressed in this structure to filter external data for validity.
 *
 * Note that this does not relate to the topics of projection (rendering),
 * where resolutions went up to 1024x768 in the original game,
 * and screen resolutions of 4K UHD and more at the time of writing this.
 */
typedef struct
{
   PixelAxisPosition minPosition;
   PixelAxisPosition maxPosition;
   PixelAxisOffset minOffset;
   PixelAxisOffset maxOffset;
   PixelAxisSize minSize;
   PixelAxisSize maxSize;
} PixelSpaceLimits;

typedef struct
{
   PixelAxisPosition x;
   PixelAxisPosition y;
} PixelPosition;

typedef struct
{
   PixelAxisOffset x;
   PixelAxisOffset y;
} PixelOffset;

typedef struct
{
   PixelAxisSize width;
   PixelAxisSize height;
} PixelSize;

extern PixelSpaceLimits pixelSpaceLimits();

#ifdef __cplusplus
}
#endif
