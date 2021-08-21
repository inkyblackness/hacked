package object

import "fmt"

// DamageType classifies damage.
type DamageType byte

// String returns the textual representation of the value.
func (damageType DamageType) String() string {
	if int(damageType) >= len(damageTypeNames) {
		return fmt.Sprintf("Unknown 0x%02X", int(damageType))
	}
	return damageTypeNames[damageType]
}

func (damageType DamageType) mask() DamageTypeMask {
	return 1 << damageType
}

// DamageType constants.
const (
	DamageTypeExplosion    DamageType = 0
	DamageTypeEnergy       DamageType = 1
	DamageTypeMagnetic     DamageType = 2
	DamageTypeRadiation    DamageType = 3
	DamageTypeGas          DamageType = 4
	DamageTypeTranquilizer DamageType = 5
	DamageTypeNeedle       DamageType = 6
	DamageTypeBio          DamageType = 7
)

var damageTypeNames = []string{
	"Explosion", "Energy", "Magnetic", "Radiation",
	"Gas", "Tranquilizer", "Needle", "Bio",
}

// DamageTypes returns all known constants.
func DamageTypes() []DamageType {
	return []DamageType{
		DamageTypeExplosion, DamageTypeEnergy, DamageTypeMagnetic, DamageTypeRadiation,
		DamageTypeGas, DamageTypeTranquilizer, DamageTypeNeedle, DamageTypeBio,
	}
}

// DamageTypeMask combines a set of damage types.
type DamageTypeMask byte

// Has returns whether the mask contains the specified value.
func (mask DamageTypeMask) Has(dmg DamageType) bool {
	return (mask & dmg.mask()) != 0
}

// With returns a new mask that specifies the combination of this and the given damage type.
func (mask DamageTypeMask) With(dmg DamageType) DamageTypeMask {
	return mask | dmg.mask()
}

// Without returns a new mask that specifies the remainder of this, not including the given type.
func (mask DamageTypeMask) Without(dmg DamageType) DamageTypeMask {
	return mask & ^dmg.mask()
}

// SpecialDamageType is a combination of "primary" (double) and "super" (quadruple) damage potential.
// The identifier are "freeform" enumerations without public constants.
type SpecialDamageType byte

const (
	specialDamageTypePrimaryShift = 0
	specialDamageTypePrimaryMask  = 0x0F
	specialDamageTypeSuperShift   = 4
	specialDamageTypeSuperMask    = 0xF0

	// SpecialDamageTypeLimit identifies the maximum value of special damage type.
	SpecialDamageTypeLimit = 0x0F
)

// PrimaryValue returns the double damage type.
func (dmg SpecialDamageType) PrimaryValue() int {
	return int((dmg & specialDamageTypePrimaryMask) >> specialDamageTypePrimaryShift)
}

// WithPrimaryValue returns a new type instance with the given primary value set.
func (dmg SpecialDamageType) WithPrimaryValue(value int) SpecialDamageType {
	result := dmg
	if (value >= 0) && (value <= SpecialDamageTypeLimit) {
		result = SpecialDamageType((byte(result) & ^byte(specialDamageTypePrimaryMask)) | byte(value<<specialDamageTypePrimaryShift))
	}
	return result
}

// SuperValue returns the quadruple damage type.
func (dmg SpecialDamageType) SuperValue() int {
	return int((dmg & specialDamageTypeSuperMask) >> specialDamageTypeSuperShift)
}

// WithSuperValue returns a new type instance with the given super value set.
func (dmg SpecialDamageType) WithSuperValue(value int) SpecialDamageType {
	result := dmg
	if (value >= 0) && (value <= SpecialDamageTypeLimit) {
		result = SpecialDamageType((byte(result) & ^byte(specialDamageTypeSuperMask)) | byte(value<<specialDamageTypeSuperShift))
	}
	return result
}
