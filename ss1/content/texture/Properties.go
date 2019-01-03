package texture

import "github.com/inkyblackness/hacked/ss1/serial"

const (
	// PropertiesSize specifies, in bytes, the length a properties structure has.
	PropertiesSize = 11
)

// Properties describe additional information about a texture.
type Properties struct {
	_                byte
	_                byte
	_                int16
	DistanceModifier int16
	Climbable        byte
	_                byte

	TransparencyControl TransparencyControl

	AnimationGroup byte
	AnimationIndex byte
}

// PropertiesList implements serial.Codable to serialize a set of texture properties.
type PropertiesList []Properties

// Code serializes the list with the provided coder.
func (list PropertiesList) Code(coder serial.Coder) {
	version := propertiesFileVersion
	coder.Code(&version)
	for i := 0; i < len(list); i++ {
		coder.Code(&list[i])
	}
}
