#pragma once

#include <string>
#include <vector>

#include "hacked/core/system/Validation.h"

extern void assertResultMessages(ValidationResult const &result, std::vector<std::string> const &messages);
