package texture

import "fmt"

// TransparencyControl describes how a texture should be displayed.
type TransparencyControl byte

var transparencyControlNames = []string{
	"Regular",
	"Space",
	"SpaceBackground",
}

func (ctrl TransparencyControl) String() string {
	if int(ctrl) < len(transparencyControlNames) {
		return transparencyControlNames[ctrl]
	}
	return fmt.Sprintf("Unknown%02X", int(ctrl))
}

const (
	// TransparencyControlRegular describes a typical texture.
	TransparencyControlRegular TransparencyControl = 0x00
	// TransparencyControlSpace marks space texture, which ignores bitmap data.
	TransparencyControlSpace TransparencyControl = 0x01
	// TransparencyControlSpaceBackground draws space for palette index 0x00.
	TransparencyControlSpaceBackground TransparencyControl = 0x02
)

// TransparencyControls returns all known transparency control values.
func TransparencyControls() []TransparencyControl {
	return []TransparencyControl{TransparencyControlRegular, TransparencyControlSpace, TransparencyControlSpaceBackground}
}
