package object

// RenderType defines how an object is to be rendered.
type RenderType byte

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
