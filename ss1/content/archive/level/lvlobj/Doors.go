package lvlobj

import (
	"fmt"

	"github.com/inkyblackness/hacked/ss1/content/interpreters"
)

var baseDoor = interpreters.New().
	With("LockVariableIndex", 0, 2).As(interpreters.RangedValue(0, 0x1FF)).
	With("LockMessageIndex", 2, 1).As(interpreters.FormattedRangedValue(0, 255,
	func(value int) (result string) {
		return fmt.Sprintf("%d", value+7)
	})).
	With("ForceDoorColor", 3, 1).As(interpreters.FormattedRangedValue(0, 255,
	func(value int) (result string) {
		if colorText, defined := forceColors[value]; defined {
			result = fmt.Sprintf("%s", colorText)
		} else {
			result = ""
		}
		return result
	})).
	With("RequiredAccessLevel", 4, 1).As(interpreters.FormattedRangedValue(0, 255,
	func(value int) (result string) {
		if value == 255 {
			result = fmt.Sprintf("SHODAN")
		} else if accessLevel, known := accessLevelMasks[1<<uint32(value)]; known {
			result = fmt.Sprintf("%s", accessLevel)
		} else {
			result = "Unknown"
		}
		return
	})).
	With("AutoCloseTime", 5, 1).As(interpreters.FormattedRangedValue(0, 255,
	func(value int) string {
		return fmt.Sprintf("%.2f sec", float64(value)*0.5)
	})).
	With("OtherObjectID", 6, 2).As(interpreters.ObjectID())

func initDoors() interpreterRetriever {
	return newInterpreterLeaf(baseDoor)
}
