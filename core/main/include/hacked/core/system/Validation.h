#pragma once

#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

#define VALIDATION_RESULT_MAX_MESSAGES 10

typedef struct
{
   /**
    * Message stack from a validation. By default, if no message is set, then validation was successful.
    * The stack grows from index 0 with the innermost detail to the last available element.
    * In case more messages are added than possible in the array, the last message pointer will be overwritten.
    * This way the uppermost layer will always be visible (the scope), plus the most nested detail that caused the initial validation to fail.
    * A NULL entry terminates the list.
    *
    * Strings that these pointers refer to must outlive the lifetime of the structure, and thus be static.
    */
   char const *messages[VALIDATION_RESULT_MAX_MESSAGES];
} ValidationResult;

/**
 * Returns true if the validation result indicates issues.
 *
 * @param result the result to check
 * @return true if there are issues, false otherwise
 */
extern bool validationResultHasFailure(ValidationResult const *result);

/**
 * Add a validation message to the result.
 *
 * @param result the result to extend, it will be modified
 * @param message the message to add
 * @return a copy of the result, for quick returns in case no further validation is needed
 */
extern ValidationResult validationResultAddMessage(ValidationResult *result, char const *message);

/**
 * Add a validation message to the result IF the condition is false.
 *
 * @param result the result to extend, it will be modified
 * @param message the message to add
 * @param isValid the condition under which to add given message
 * @return a copy of the result, for quick returns in case no further validation is needed
 */
extern ValidationResult validateResultAddConditional(ValidationResult *result, bool isValid, char const *message);

/**
 * Merges one validation result with a target. It is a convenience function to add the results
 * of another validation with one's own, plus adding a message IF the other validation result indicates a failure.
 *
 * @param result the result to extend, it will be modified if other indicates a failure
 * @param other the other result to add to result
 * @param message a message to add if other indicates a failure
 * @return a copy of the result, for quick returns in case no further validation is needed
 */
extern ValidationResult validationResultMerge(ValidationResult *result, ValidationResult other, char const *message);

/**
 * A function that will abort the program in case the result indicates a failure.
 *
 * @param result the result to verify.
 */
extern void validationResultAssert(ValidationResult result);

typedef void (*ValidationResultAssertFunc)(ValidationResult result);

/**
 * Register an assertion function that is getting called when validationResultAssert determined
 * a result to indicate a failure.
 * The registered function should somehow report the messages and must then terminate the application.
 * It is not expected that this registered function returns.
 * If no function is registered or it returns, nothing will happen and the application will
 * most likely crash either silently, or spectacularly.
 *
 * @param func the function to register
 */
extern void validationResultSetAssertFunc(ValidationResultAssertFunc func);

#ifndef NDEBUG
/**
 * Use this macro for performance-intensive places or places where failures should be determined during development.
 * It will resolve to nothing when in release build.
 *
 * @param result the ValidationResult instance to assert on
 */
#define validationResultAssertDebug(result) validationResultAssert(result)
#else
#define validationResultAssertDebug(result)                                                                                                                    \
   do                                                                                                                                                          \
   {                                                                                                                                                           \
   } while (0)
#endif

#ifdef __cplusplus
}
#endif
