#include <gtest/gtest.h>

#include "hacked/infrastructure/compliance/Licenses.h"

TEST(LicensesTest, count)
{
   EXPECT_GT(licensesGetLicenseCount(), 0);
}

TEST(LicensesTest, licensesHaveTheirMandatoryInfoFilledOut)
{
   size_t const count = licensesGetLicenseCount();
   for (size_t i = 0; i < count; i++)
   {
      LicenseInfo const *const info = licensesGetLicense(i);
      ASSERT_TRUE(info != NULL) << "License at " << std::to_string(i) << " is wrong";
      EXPECT_TRUE((info->title != NULL) && (info->title[0] != 0x00)) << "License at " << std::to_string(i) << " has no title";
      EXPECT_TRUE((info->text != NULL) && (info->text[0] != 0x00)) << "License at " << std::to_string(i) << " has no text";
   }
}

TEST(LicensesTest, invalidLicenseIsNull)
{
   LicenseInfo const *const info = licensesGetLicense(std::numeric_limits<size_t>::max());
   EXPECT_TRUE(info == NULL) << "something invalid got returned";
}
