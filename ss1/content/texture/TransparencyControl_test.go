package texture_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/texture"

	"github.com/stretchr/testify/assert"
)

func TestTransparencyControlString(t *testing.T) {
	tt := []struct {
		ctrl     texture.TransparencyControl
		expected string
	}{
		{ctrl: texture.TransparencyControlRegular, expected: "Regular"},
		{ctrl: texture.TransparencyControlSpace, expected: "Space"},
		{ctrl: texture.TransparencyControlSpaceBackground, expected: "SpaceBackground"},
		{ctrl: 0x80, expected: "Unknown80"},
	}

	for _, tc := range tt {
		result := tc.ctrl.String()
		assert.Equal(t, tc.expected, result)
	}
}

func TestTransparencyControls(t *testing.T) {
	controls := texture.TransparencyControls()
	assert.Equal(t, 3, len(controls))
}
