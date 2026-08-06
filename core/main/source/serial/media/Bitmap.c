#include "hacked/core/media/Bitmap.h"

ValidationResult bitmapValidate(Bitmap const *const bitmap)
{
   ValidationResult result = {};
   if (bitmap == NULL)
   {
      return validationResultAddMessage(&result, "bitmap is NULL");
   }
   validateResultAddConditional(&result, bitmap->data != NULL, "data is NULL");
   validateResultAddConditional(&result, bitmap->size.width > 0, "width is zero");
   validateResultAddConditional(&result, bitmap->size.height > 0, "height is zero");
   validateResultAddConditional(&result, bitmap->stride > 0, "stride is zero");
   validateResultAddConditional(&result, bitmap->size.width <= bitmap->stride, "width is higher than stride");
   validateResultAddConditional(&result, ((size_t)bitmap->size.height * bitmap->stride) <= bitmap->dataLength, "stride times height is more than dataLength");
   return result;
}
