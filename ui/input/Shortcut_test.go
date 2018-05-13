package input_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ui/input"

	"github.com/stretchr/testify/assert"
)

func TestResolveShortcutReturnsValuesOfKnownCombo(t *testing.T) {
	key, knownKey := input.ResolveShortcut("c", input.ModControl)

	assert.True(t, knownKey)
	assert.Equal(t, input.KeyCopy, key)
}

func TestResolveShortcutReturnsFalseForUnknownCombo(t *testing.T) {
	_, knownKey := input.ResolveShortcut("l", input.ModControl.With(input.ModShift))

	assert.False(t, knownKey)
}
