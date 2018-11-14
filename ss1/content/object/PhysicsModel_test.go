package object_test

import (
	"fmt"
	"github.com/inkyblackness/hacked/ss1/content/object"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPhysicsModels(t *testing.T) {
	assert.Equal(t, 3, len(object.PhysicsModels()))
}

func TestPhysicsModelString(t *testing.T) {
	tt := []struct {
		model    object.PhysicsModel
		expected string
	}{
		{model: object.PhysicsModelInsubstantial, expected: "Insubstantial"},
		{model: object.PhysicsModelRegular, expected: "Regular"},
		{model: object.PhysicsModelStrange, expected: "Strange"},

		{model: object.PhysicsModel(255), expected: "Unknown 0xFF"},
	}

	for _, tc := range tt {
		result := tc.model.String()
		assert.Equal(t, tc.expected, result, fmt.Sprintf("Failed for 0x%02X", int(tc.model)))
	}
}
