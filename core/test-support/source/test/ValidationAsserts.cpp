#include <cstddef>

#include <gtest/gtest.h>

#include "hacked/core/test/ValidationAsserts.h"

void assertResultMessages(ValidationResult const &result, std::vector<std::string> const &messages)
{
   for (size_t i = 0; i < messages.size(); i++)
   {
      if (result.messages[i] != nullptr)
      {
         EXPECT_EQ(std::string(result.messages[i]), messages.at(i)) << "failed at " << std::to_string(i);
      }
      else
      {
         EXPECT_TRUE(false) << "NULL at " << std::to_string(i) << ", expecting '" << messages[i] << "'";
      }
   }
   for (size_t i = messages.size(); i < VALIDATION_RESULT_MAX_MESSAGES; i++)
   {
      EXPECT_TRUE(result.messages[i] == nullptr) << "tail not null at " << std::to_string(i) << ", has '" << result.messages[i] << "'";
   }
}
