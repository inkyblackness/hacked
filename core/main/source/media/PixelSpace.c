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

ValidationResult pixelSpaceValidateSize(PixelSize const size)
{
   PixelSpaceLimits const limits = pixelSpaceLimits();
   ValidationResult result = {0};
   validationResultAddConditional(&result, size.width >= limits.minSize, "width below limit");
   validationResultAddConditional(&result, size.width <= limits.maxSize, "width above limit");
   validationResultAddConditional(&result, size.height >= limits.minSize, "height below limit");
   validationResultAddConditional(&result, size.height <= limits.maxSize, "height above limit");
   return result;
}

bool pixelSpaceAreaContainsRect(PixelRect const area, PixelRect const rect)
{
   if ((rect.topLeft.x < area.topLeft.x) || (rect.topLeft.y < area.topLeft.y))
   {
      return false;
   }
   PixelAxisOffset areaRight = (PixelAxisOffset)area.topLeft.x + (PixelAxisOffset)area.size.width;
   PixelAxisOffset rectRight = (PixelAxisOffset)rect.topLeft.x + (PixelAxisOffset)rect.size.width;
   if (areaRight < rectRight)
   {
      return false;
   }
   PixelAxisOffset areaBottom = (PixelAxisOffset)area.topLeft.y + (PixelAxisOffset)area.size.height;
   PixelAxisOffset rectBottom = (PixelAxisOffset)rect.topLeft.y + (PixelAxisOffset)rect.size.height;
   return areaBottom >= rectBottom;
}
