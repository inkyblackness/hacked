#include "hacked/core/media/PixelSpace.h"

// The limits were set based on a best guess, what is feasible and probably never needs to be extended.
#define PIXEL_SPACE_AXIS_LIMIT 2048
// The factor allows offset calculation to accumulate, if needed;
// Yet the result still must be within the position/size limits.
#define PIXEL_SPACE_OFFSET_LIMIT_FACTOR 8

PixelSpaceLimits pixelSpaceLimits()
{
   static PixelSpaceLimits const limits = {
      .minPosition = -PIXEL_SPACE_AXIS_LIMIT,
      .maxPosition = PIXEL_SPACE_AXIS_LIMIT,
      .minOffset = PIXEL_SPACE_AXIS_LIMIT * -PIXEL_SPACE_OFFSET_LIMIT_FACTOR,
      .maxOffset = PIXEL_SPACE_AXIS_LIMIT * PIXEL_SPACE_OFFSET_LIMIT_FACTOR,
      .minSize = 1,
      .maxSize = PIXEL_SPACE_AXIS_LIMIT,
   };
   return limits;
}
