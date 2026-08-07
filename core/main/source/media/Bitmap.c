#include "hacked/core/media/Bitmap.h"

ValidationResult bitmapValidate(Bitmap const *const bitmap)
{
   ValidationResult result = {0};
   if (bitmap == NULL)
   {
      return validationResultAddMessage(&result, "bitmap is NULL");
   }
   validationResultAddConditional(&result, bitmap->data != NULL, "data is NULL");
   validationResultAddConditional(&result, bitmap->size.width > 0, "width is zero");
   validationResultAddConditional(&result, bitmap->size.height > 0, "height is zero");
   validationResultAddConditional(&result, bitmap->stride > 0, "stride is zero");
   validationResultAddConditional(&result, bitmap->size.width <= bitmap->stride, "width is higher than stride");
   validationResultAddConditional(&result, ((size_t)bitmap->size.height * bitmap->stride) <= bitmap->dataLength, "stride times height is more than dataLength");
   return result;
}
