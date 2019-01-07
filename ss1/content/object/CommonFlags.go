package object

import "fmt"

// CommonFlag describe flags applicable to all object classes.
type CommonFlag int

// String returns the textual representation of the value.
func (flag CommonFlag) String() string {
	if int(flag) >= len(commonFlagNames) {
		return fmt.Sprintf("Unknown 0x%04X", int(flag))
	}
	return commonFlagNames[flag]
}

func (flag CommonFlag) mask() CommonFlagField {
	return CommonFlagField(1<<uint(flag)) & commonFlagFieldMask
}

// Common flags.
const (
	CommonFlagUseful                CommonFlag = 0
	CommonFlagSolid                 CommonFlag = 1
	CommonFlagUsableInventoryObject CommonFlag = 4
	CommonFlagBlockRendering        CommonFlag = 5
	CommonFlagSolidIfClosed         CommonFlag = 8
	CommonFlagFlatSolid             CommonFlag = 9
	CommonFlagDoubledBitmapSize     CommonFlag = 10
	CommonFlagDestroyOnContact      CommonFlag = 11
	CommonFlagClassSpecific1        CommonFlag = 12
	CommonFlagClassSpecific2        CommonFlag = 13
	CommonFlagClassSpecific3        CommonFlag = 14
	CommonFlagUseless               CommonFlag = 15
)

var commonFlagNames = []string{
	"Useful",
	"Solid",
	"Not a Flag: 0x0004",
	"Not a Flag: 0x0008",
	"UsableInventoryObject",
	"BlockRendering",
	"Not a Flag: 0x0040",
	"Not a Flag: 0x0080",
	"SolidIfClosed",
	"FlatSolid",
	"DoubleBitmapSize",
	"DestroyOnContact",
	"ClassSpecific1",
	"ClassSpecific2",
	"ClassSpecific3",
	"Useless",
}

// CommonFlags returns all known constants.
func CommonFlags() []CommonFlag {
	return []CommonFlag{
		CommonFlagUseful, CommonFlagSolid, CommonFlagUsableInventoryObject, CommonFlagBlockRendering,
		CommonFlagSolidIfClosed, CommonFlagFlatSolid, CommonFlagDoubledBitmapSize, CommonFlagDestroyOnContact,
		CommonFlagClassSpecific1, CommonFlagClassSpecific2, CommonFlagClassSpecific3, CommonFlagUseless,
	}
}

const (
	commonFlagFieldMask = 0xFF33

	commonFlagLightTypeShift = 6
	commonFlagLightTypeMask  = 0x00C0
	commonFlagUseModeShift   = 2
	commonFlagUseModeMask    = 0x000C
)

// CommonFlagField describes simple features applicable for all object classes.
type CommonFlagField uint16

// Has returns whether the field has given flag set.
func (field CommonFlagField) Has(flag CommonFlag) bool {
	return (field & flag.mask()) != 0
}

// With returns a new field instance with given flag set.
func (field CommonFlagField) With(flag CommonFlag) CommonFlagField {
	return field | flag.mask()
}

// Without returns a new field instance without the given flag set.
func (field CommonFlagField) Without(flag CommonFlag) CommonFlagField {
	return field & ^flag.mask()
}

// LightType returns the stored type enumeration value.
func (field CommonFlagField) LightType() LightType {
	return LightType((field & commonFlagLightTypeMask) >> commonFlagLightTypeShift)
}

// WithLightType returns a new flag field value with given light type set.
func (field CommonFlagField) WithLightType(value LightType) CommonFlagField {
	return CommonFlagField((uint16(field) & ^uint16(commonFlagLightTypeMask)) | ((uint16(value) << commonFlagLightTypeShift) & commonFlagLightTypeMask))
}

// UseMode returns the stored enumeration value.
func (field CommonFlagField) UseMode() UseMode {
	return UseMode((field & commonFlagUseModeMask) >> commonFlagUseModeShift)
}

// WithUseMode returns a new flag field value with given use mode set.
func (field CommonFlagField) WithUseMode(value UseMode) CommonFlagField {
	return CommonFlagField((uint16(field) & ^uint16(commonFlagUseModeMask)) | ((uint16(value) << commonFlagUseModeShift) & commonFlagUseModeMask))
}

// LightType describes how an object is affected by lighting.
type LightType byte

// String returns the textual representation of the value.
func (mode LightType) String() string {
	if int(mode) >= len(lightTypeNames) {
		return fmt.Sprintf("Unknown 0x%02X", int(mode))
	}
	return lightTypeNames[mode]
}

// LightType constants.
const (
	LightTypeSimple   LightType = 0x00
	LightTypeComplex  LightType = 0x01
	LightTypeIgnore   LightType = 0x02
	LightTypeDeferred LightType = 0x03
)

var lightTypeNames = []string{"Simple", "Complex", "Ignore", "Deferred"}

// LightTypes returns all known constants.
func LightTypes() []LightType {
	return []LightType{LightTypeSimple, LightTypeComplex, LightTypeIgnore, LightTypeDeferred}
}

// UseMode describes how an object can be interacted with.
type UseMode byte

// String returns the textual representation of the value.
func (mode UseMode) String() string {
	if int(mode) >= len(useModeNames) {
		return fmt.Sprintf("Unknown 0x%02X", int(mode))
	}
	return useModeNames[mode]
}

// UseMode constants.
const (
	UseModePickup UseMode = 0x00
	UseModeUse    UseMode = 0x01
	UseModeNull   UseMode = 0x03
)

var useModeNames = []string{"Pickup", "Use", "Unknown 0x02", "Null"}

// UseModes returns all known constants.
func UseModes() []UseMode {
	return []UseMode{UseModePickup, UseModeUse, UseModeNull}
}
