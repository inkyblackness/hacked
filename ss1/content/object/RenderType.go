package object

import "fmt"

// RenderType defines how an object is to be rendered.
type RenderType byte

// String returns the textual representation of the value.
func (renderType RenderType) String() string {
	if int(renderType) >= len(renderTypeNames) {
		return fmt.Sprintf("Unknown 0x%02X", int(renderType))
	}
	return renderTypeNames[renderType]
}

// RenderType constants.
const (
	RenderTypeUnknown   RenderType = 0
	RenderTypeTextPoly  RenderType = 1
	RenderTypeBitmap    RenderType = 2
	RenderTypeTPoly     RenderType = 3
	RenderTypeCritter   RenderType = 4
	RenderTypeAnimPoly  RenderType = 5
	RenderTypeVox       RenderType = 6
	RenderTypeNoObject  RenderType = 7
	RenderTypeTexBitmap RenderType = 8
	RenderTypeFlatPoly  RenderType = 9
	RenderTypeMultiView RenderType = 10
	RenderTypeSpecial   RenderType = 11
	RenderTypeTLPoly    RenderType = 12
)

var renderTypeNames = []string{
	"Unknown 0x00",
	"TextPoly", "Bitmap", "TPoly", "Critter",
	"AnimPoly", "Vox", "NoObject", "TexBitmap",
	"FlatPoly", "MultiView", "Special", "TLPoly",
}

// RenderTypes returns all known constants.
func RenderTypes() []RenderType {
	return []RenderType{
		RenderTypeUnknown,
		RenderTypeTextPoly, RenderTypeBitmap, RenderTypeTPoly, RenderTypeCritter,
		RenderTypeAnimPoly, RenderTypeVox, RenderTypeNoObject, RenderTypeTexBitmap,
		RenderTypeFlatPoly, RenderTypeMultiView, RenderTypeSpecial, RenderTypeTLPoly,
	}
}
