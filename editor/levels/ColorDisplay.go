package levels

import (
	"fmt"
)

// ColorDisplay is an enumeration to how a tile should be colored.
type ColorDisplay int

// ColorDisplay constants are listed below.
const (
	ColorDisplayNone    ColorDisplay = 0
	ColorDisplayFloor   ColorDisplay = 1
	ColorDisplayCeiling ColorDisplay = 2
)

// String returns a textual representation.
func (display ColorDisplay) String() string {
	switch display {
	case ColorDisplayNone:
		return "None"
	case ColorDisplayFloor:
		return "Floor"
	case ColorDisplayCeiling:
		return "Ceiling"
	default:
		return fmt.Sprintf("Unknown%d", int(display))
	}
}

// ColorDisplays returns all ColorDisplay constants.
func ColorDisplays() []ColorDisplay {
	return []ColorDisplay{ColorDisplayNone, ColorDisplayFloor, ColorDisplayCeiling}
}
