package world_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/world"

	"github.com/stretchr/testify/assert"
)

func TestDefaultCyberspaceLevels(t *testing.T) {
	tt := []bool{
		false, false, false, false,
		false, false, false, false,
		false, false, true, false,
		false, false, true, true,
	}

	for levelID, expected := range tt {
		result := world.IsConsideredCyberspaceByDefault(levelID)
		assert.Equal(t, expected, result, "Cyberspace mismatch")
	}
}
