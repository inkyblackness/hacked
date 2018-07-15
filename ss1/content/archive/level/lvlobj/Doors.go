package lvlobj

import (
	"fmt"

	"github.com/inkyblackness/hacked/ss1/content/interpreters"
)

var baseDoor = interpreters.New().
	With("LockVariableIndex", 0, 2).As(interpreters.RangedValue(0, 0x1FF)).
	With("LockMessageIndex", 2, 1).As(interpreters.FormattedRangedValue(0, 255,
	func(value int64) (result string) {
		return fmt.Sprintf("%d", value+7)
	})).
	With("ForceDoorColor", 3, 1).As(interpreters.FormattedRangedValue(0, 255,
	func(value int64) (result string) {
		if colorText, defined := forceColors[value]; defined {
			result = fmt.Sprintf("%s  - raw: %d", colorText, value)
		} else {
			result = fmt.Sprintf("%d", value)
		}
		return result
	})).
	With("RequiredAccessLevel", 4, 1).As(interpreters.FormattedRangedValue(0, 255,
	func(value int64) (result string) {
		if value == 255 {
			result = fmt.Sprintf("SHODAN - raw: 255")
		} else if accessLevel, known := accessLevelMasks[1<<uint32(value)]; known {
			result = fmt.Sprintf("%s  - raw: %d", accessLevel, value)
		} else {
			result = fmt.Sprintf("Unknown  - raw: %d", value)
		}
		return
	})).
	With("AutoCloseTime", 5, 1).As(interpreters.FormattedRangedValue(0, 255,
	func(value int64) string {
		return fmt.Sprintf("%.2f sec  - raw: %d", float64(value)*0.5, value)
	})).
	With("OtherObjectID", 6, 2).As(interpreters.ObjectID())

func initDoors() interpreterRetriever {
	return newInterpreterLeaf(baseDoor)
}
