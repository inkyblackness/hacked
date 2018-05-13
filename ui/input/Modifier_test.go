package input_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ui/input"

	"github.com/stretchr/testify/assert"
)

func TestHasReturnsTrueForThemselves(t *testing.T) {
	mods := []input.Modifier{input.ModShift, input.ModControl, input.ModAlt, input.ModSuper}

	for _, mod := range mods {
		assert.True(t, mod.Has(mod))
	}
}

func TestHasReturnsFalseForDifferentModifiers(t *testing.T) {
	mods := []input.Modifier{input.ModShift, input.ModControl, input.ModAlt, input.ModSuper}

	for index, mod := range mods {
		nextMod := mods[(index+1)%len(mods)]
		assert.False(t, mod.Has(nextMod))
	}
}

func TestHasReturnsFalseForSetsIncludingOthers(t *testing.T) {
	mods := []input.Modifier{input.ModShift, input.ModControl, input.ModAlt, input.ModSuper}

	for index, mod := range mods {
		nextMod := mods[(index+1)%len(mods)]
		assert.False(t, mod.Has(mod.With(nextMod)))
	}
}

func TestWithoutReturnsReduction(t *testing.T) {
	mod := input.ModShift.With(input.ModAlt).Without(input.ModShift)

	assert.Equal(t, input.ModAlt, mod)
}
