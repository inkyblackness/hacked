package object_test

import (
	"fmt"
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/object"

	"github.com/stretchr/testify/assert"
)

func TestDamageTypes(t *testing.T) {
	assert.Equal(t, 8, len(object.DamageTypes()))
}

func TestDamageTypeString(t *testing.T) {
	tt := []struct {
		damageType object.DamageType
		expected   string
	}{
		{damageType: object.DamageTypeExplosion, expected: "Explosion"},
		{damageType: object.DamageTypeEnergy, expected: "Energy"},
		{damageType: object.DamageTypeMagnetic, expected: "Magnetic"},
		{damageType: object.DamageTypeRadiation, expected: "Radiation"},
		{damageType: object.DamageTypeGas, expected: "Gas"},
		{damageType: object.DamageTypeTranquilizer, expected: "Tranquilizer"},
		{damageType: object.DamageTypeNeedle, expected: "Needle"},
		{damageType: object.DamageTypeBio, expected: "Bio"},

		{damageType: object.DamageType(255), expected: "Unknown 0xFF"},
	}

	for _, tc := range tt {
		result := tc.damageType.String()
		assert.Equal(t, tc.expected, result, fmt.Sprintf("Failed for 0x%02X", int(tc.damageType)))
	}
}

func TestDamageTypeMask(t *testing.T) {
	allTypes := object.DamageTypes()
	noDmg := object.DamageTypeMask(0)
	allDmg := ^noDmg
	for _, dmg := range allTypes {
		assert.False(t, noDmg.Has(dmg), fmt.Sprintf("empty mask shouldn't contain damage %v", dmg))
		assert.True(t, allDmg.Has(dmg), fmt.Sprintf("full mask should contain damage %v", dmg))
		assert.True(t, noDmg.With(dmg).Has(dmg), fmt.Sprintf("With() should return mask including damage %v", dmg))
		assert.False(t, allDmg.Without(dmg).Has(dmg), fmt.Sprintf("Without() should return mask excluding damage %v", dmg))
	}
}

func TestSpecialDamageType(t *testing.T) {
	for super := 0; super <= object.SpecialDamageTypeLimit; super++ {
		for primary := 0; primary <= object.SpecialDamageTypeLimit; primary++ {
			key := fmt.Sprintf("%v:%v", super, primary)
			result1 := object.SpecialDamageType(0).WithSuperValue(super).WithPrimaryValue(primary)
			result2 := object.SpecialDamageType(0).WithPrimaryValue(primary).WithSuperValue(super)
			assert.Equal(t, result1, result2, fmt.Sprintf("value should be equal regardless of modifier ordering in %v", key))
			assert.Equal(t, super, result1.SuperValue(), fmt.Sprintf("super value not found in %v", key))
			assert.Equal(t, primary, result1.PrimaryValue(), fmt.Sprintf("primary value not found in %v", key))
		}
	}
}

func TestSpecialDamageTypeLimits(t *testing.T) {
	base := object.SpecialDamageType(0).WithSuperValue(object.SpecialDamageTypeLimit / 2).WithPrimaryValue(object.SpecialDamageTypeLimit/2 - 1)

	assert.Equal(t, base, base.WithSuperValue(-1), "negative value not allowed for super")
	assert.Equal(t, base, base.WithSuperValue(object.SpecialDamageTypeLimit+1), "beyond limit value not allowed for super")
	assert.Equal(t, base, base.WithPrimaryValue(-1), "negative value not allowed for primary")
	assert.Equal(t, base, base.WithPrimaryValue(object.SpecialDamageTypeLimit+1), "beyond limit value not allowed for primary")
}
