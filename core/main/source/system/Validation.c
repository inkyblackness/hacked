#include <stddef.h>
#ifndef NDEBUG
#include <assert.h>
#endif

#include "hacked/core/system/Validation.h"

static ValidationResultAssertFunc assertFunc = NULL;

bool validationResultHasFailure(ValidationResult const *const result)
{
   return result->messages[0] != NULL;
}

ValidationResult validationResultAddMessage(ValidationResult *const result, char const *const message)
{
   for (size_t i = 0; i < (VALIDATION_RESULT_MAX_MESSAGES - 1); i++)
   {
      if (result->messages[i] == NULL)
      {
         result->messages[i] = message;
         return *result;
      }
   }
   result->messages[VALIDATION_RESULT_MAX_MESSAGES - 1] = message;
   return *result;
}

ValidationResult validateResultAddConditional(ValidationResult *result, bool isValid, char const *message)
{
   if (isValid)
   {
      return *result;
   }
   return validationResultAddMessage(result, message);
}

ValidationResult validationResultMerge(ValidationResult *const result, ValidationResult const other, char const *const message)
{
   size_t copyIndex = 0;
   for (size_t i = 0; i < (VALIDATION_RESULT_MAX_MESSAGES - 1); i++)
   {
      if (result->messages[i] == NULL)
      {
         result->messages[i] = other.messages[copyIndex];
         copyIndex++;
         if (result->messages[i] == NULL)
         {
            result->messages[i] = message;
            return *result;
         }
      }
   }
   result->messages[VALIDATION_RESULT_MAX_MESSAGES - 1] = message;
   return *result;
}

void validationResultSetAssertFunc(ValidationResultAssertFunc const func)
{
   assertFunc = func;
}

void validationResultAssert(ValidationResult const result)
{
   if ((result.messages[0] != NULL) && (assertFunc != NULL))
   {
      assertFunc(result);
   }
}
