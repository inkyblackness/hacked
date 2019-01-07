package object_test

import (
	"fmt"
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/object"

	"github.com/stretchr/testify/assert"
)

func TestLightTypes(t *testing.T) {
	assert.Equal(t, 4, len(object.LightTypes()))
}

func TestLightTypeString(t *testing.T) {
	tt := []struct {
		lightType object.LightType
		expected  string
	}{
		{lightType: object.LightTypeSimple, expected: "Simple"},
		{lightType: object.LightTypeComplex, expected: "Complex"},
		{lightType: object.LightTypeIgnore, expected: "Ignore"},
		{lightType: object.LightTypeDeferred, expected: "Deferred"},

		{lightType: object.LightType(255), expected: "Unknown 0xFF"},
	}

	for _, tc := range tt {
		result := tc.lightType.String()
		assert.Equal(t, tc.expected, result, fmt.Sprintf("Failed for 0x%02X", int(tc.lightType)))
	}
}

func TestUseModes(t *testing.T) {
	assert.Equal(t, 3, len(object.UseModes()))
}

func TestUseModeString(t *testing.T) {
	tt := []struct {
		useMode  object.UseMode
		expected string
	}{
		{useMode: object.UseModePickup, expected: "Pickup"},
		{useMode: object.UseModeUse, expected: "Use"},
		{useMode: object.UseModeNull, expected: "Null"},

		{useMode: object.UseMode(255), expected: "Unknown 0xFF"},
	}

	for _, tc := range tt {
		result := tc.useMode.String()
		assert.Equal(t, tc.expected, result, fmt.Sprintf("Failed for 0x%02X", int(tc.useMode)))
	}
}

func TestCommonFlagField(t *testing.T) {
	for _, useMode := range object.UseModes() {
		for _, lightType := range object.LightTypes() {
			for _, flag := range object.CommonFlags() {
				key := fmt.Sprintf("%v:%v:%v", useMode, lightType, flag)
				field := object.CommonFlagField(0).WithUseMode(useMode).WithLightType(lightType).With(flag)

				assert.Equal(t, useMode, field.UseMode(), fmt.Sprintf("UseMode wrong for %v", key))
				assert.Equal(t, lightType, field.LightType(), fmt.Sprintf("LightType wrong for %v", key))
				assert.True(t, field.Has(flag), fmt.Sprintf("flag not set for %v", key))
				assert.False(t, field.Without(flag).Has(flag), fmt.Sprintf("flag still set for %v", key))
			}
		}
	}
}

func TestCommonFlagString(t *testing.T) {
	tt := []struct {
		flag     object.CommonFlag
		expected string
	}{
		{flag: object.CommonFlagUseful, expected: "Useful"},
		{flag: object.CommonFlagSolid, expected: "Solid"},
		{flag: object.CommonFlagUsableInventoryObject, expected: "UsableInventoryObject"},
		{flag: object.CommonFlagBlockRendering, expected: "BlockRendering"},
		{flag: object.CommonFlagSolidIfClosed, expected: "SolidIfClosed"},
		{flag: object.CommonFlagFlatSolid, expected: "FlatSolid"},
		{flag: object.CommonFlagDoubledBitmapSize, expected: "DoubleBitmapSize"},
		{flag: object.CommonFlagDestroyOnContact, expected: "DestroyOnContact"},
		{flag: object.CommonFlagClassSpecific1, expected: "ClassSpecific1"},
		{flag: object.CommonFlagClassSpecific2, expected: "ClassSpecific2"},
		{flag: object.CommonFlagClassSpecific3, expected: "ClassSpecific3"},
		{flag: object.CommonFlagUseless, expected: "Useless"},

		{flag: object.CommonFlag(255), expected: "Unknown 0x00FF"},
	}

	for _, tc := range tt {
		result := tc.flag.String()
		assert.Equal(t, tc.expected, result, fmt.Sprintf("Failed for 0x%02X", int(tc.flag)))
	}
}
