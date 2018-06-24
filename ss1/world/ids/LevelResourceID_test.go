package ids_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"

	"github.com/stretchr/testify/assert"
)

func TestLevelResourceIDForLevel(t *testing.T) {
	tt := []struct {
		in       ids.LevelResourceID
		level    int
		expected resource.ID
	}{
		{0, 0, 4000},
		{1, 0, 4001},

		{5, 2, 4205},
		{42, 15, 5542},
	}

	for _, tc := range tt {
		out := tc.in.ForLevel(tc.level)
		assert.Equal(t, tc.expected, out, "Mismatch for %v -> %v", tc, out)
	}
}
