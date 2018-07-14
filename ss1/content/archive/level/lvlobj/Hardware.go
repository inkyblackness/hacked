package lvlobj

import (
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
)

var baseHardware = interpreters.New().
	With("Version", 0, 1).As(interpreters.RangedValue(0, 4))
