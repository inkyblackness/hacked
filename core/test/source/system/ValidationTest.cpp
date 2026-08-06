#include <vector>

#include <gtest/gtest.h>

#include "hacked/core/system/Validation.h"

static void assertResultMessages(ValidationResult const &result, char const *const name, std::vector<char const *> const &messages)
{
   for (size_t i = 0; i < messages.size(); i++)
   {
      EXPECT_EQ(result.messages[i], messages.at(i)) << "failed for " << std::to_string(i) << " in " << name;
   }
   for (size_t i = messages.size(); i < VALIDATION_RESULT_MAX_MESSAGES; i++)
   {
      EXPECT_TRUE(result.messages[i] == nullptr) << "tail not null at " << std::to_string(i);
   }
}

TEST(ValidationTest, initialObjectIsClean)
{
   ValidationResult constexpr result = {};
   assertResultMessages(result, "result", {});
}

TEST(ValidationTest, hasFailureReturnsFalseOnNewObject)
{
   ValidationResult constexpr result = {};
   EXPECT_FALSE(validationResultHasFailure(&result));
}

TEST(ValidationTest, hasFailureReturnsTrueWhenMessageIsPresent)
{
   ValidationResult result = {};
   result.messages[0] = "message";
   EXPECT_TRUE(validationResultHasFailure(&result));
}

TEST(ValidationTest, addMessageAddsSingleMessage)
{
   char const *const message = "test";
   ValidationResult result = {};
   ValidationResult const returnValue = validationResultAddMessage(&result, message);
   assertResultMessages(result, "result", {message});
   assertResultMessages(returnValue, "returnValue", {message});
}

TEST(ValidationTest, addMessageAddsMessageAtEnd)
{
   char const *const message = "test";
   ValidationResult result = {};
   validationResultAddMessage(&result, "message1");
   validationResultAddMessage(&result, "message2");
   ValidationResult const returnValue = validationResultAddMessage(&result, message);
   EXPECT_EQ(result.messages[2], message);
   EXPECT_EQ(returnValue.messages[2], message);
}

TEST(ValidationTest, addMessageOverridesLast)
{
   char const *const message = "test";
   ValidationResult result = {};
   for (size_t i = 0; i < VALIDATION_RESULT_MAX_MESSAGES; i++)
   {
      validationResultAddMessage(&result, "filler");
   }
   ValidationResult const returnValue = validationResultAddMessage(&result, message);
   EXPECT_EQ(result.messages[VALIDATION_RESULT_MAX_MESSAGES - 1], message);
   EXPECT_EQ(returnValue.messages[VALIDATION_RESULT_MAX_MESSAGES - 1], message);
}

TEST(ValidationTest, mergeAddsAllMessages)
{
   char const *const message1 = "message1";
   char const *const message2 = "message2";
   char const *const message3 = "message3";
   char const *const messageOther1 = "other1";
   char const *const messageOther2 = "other2";
   ValidationResult result = {};
   validationResultAddMessage(&result, message1);
   validationResultAddMessage(&result, message2);
   ValidationResult other = {};
   validationResultAddMessage(&other, messageOther1);
   validationResultAddMessage(&other, messageOther2);
   ValidationResult const returnValue = validationResultMerge(&result, other, message3);
   assertResultMessages(result, "result", {message1, message2, messageOther1, messageOther2, message3});
   assertResultMessages(returnValue, "returnValue", {message1, message2, messageOther1, messageOther2, message3});
}

TEST(ValidationTest, mergeOverwritesWithLastMessage)
{
   char const *const filler = "filler";
   char const *const message2 = "message2";
   char const *const messageOther1 = "other1";
   char const *const messageOther2 = "other2";
   ValidationResult result = {};
   std::vector<char const *> expected;
   for (size_t i = 0; i < (VALIDATION_RESULT_MAX_MESSAGES - 2); i++)
   {
      validationResultAddMessage(&result, filler);
      expected.push_back(filler);
   }
   ValidationResult other = {};
   validationResultAddMessage(&other, messageOther1);
   validationResultAddMessage(&other, messageOther2);
   ValidationResult const returnValue = validationResultMerge(&result, other, message2);
   expected.push_back(messageOther1);
   expected.push_back(message2);
   assertResultMessages(result, "result", expected);
   assertResultMessages(returnValue, "returnValue", expected);
}

TEST(ValidationTest, assertDoesNothingIfOk)
{
   ValidationResult constexpr result = {};
   validationResultAssert(result);
}

static bool assertFuncCalled = false;

static void testAssertFunc(ValidationResult const result)
{
   (void)result;
   assertFuncCalled = true;
}

TEST(ValidationTest, assertCallsFunction)
{
   ValidationResult result = {};
   validationResultAddMessage(&result, "test");
   assertFuncCalled = false;
   validationResultSetAssertFunc(testAssertFunc);
   validationResultAssert(result);
   EXPECT_TRUE(assertFuncCalled);
   validationResultSetAssertFunc(NULL);
   assertFuncCalled = false;
}

TEST(ValidationTest, assertDebugDoesNothingIfOk)
{
   ValidationResult constexpr result = {};
   (void)result;
   validationResultAssertDebug(result);
}

TEST(ValidationTest, assertWorksAccordingToBuildType)
{
   ValidationResult result = {};
   validationResultAddMessage(&result, "test");
   assertFuncCalled = false;
   validationResultSetAssertFunc(testAssertFunc);
   validationResultAssertDebug(result);
#ifdef NDEBUG
   EXPECT_FALSE(assertFuncCalled);
#else
   EXPECT_TRUE(assertFuncCalled);
#endif
   validationResultSetAssertFunc(NULL);
   assertFuncCalled = false;
}
