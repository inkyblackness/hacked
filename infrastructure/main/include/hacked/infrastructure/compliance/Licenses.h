#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include <stddef.h>

typedef struct
{
   char const *const title;
   char const *const url;
   char const *const text;
} LicenseInfo;

extern size_t licensesGetLicenseCount();
extern LicenseInfo const *licensesGetLicense(size_t index);

#ifdef __cplusplus
}
#endif
