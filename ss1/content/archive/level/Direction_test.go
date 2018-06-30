package level_test

import (
	"fmt"
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/archive/level"

	"github.com/stretchr/testify/assert"
)

func TestDirectionOffset(t *testing.T) {
	tt := []struct {
		start    level.Direction
		offset   int
		expected level.Direction
	}{
		{level.DirNorth, 0, level.DirNorth},
		{level.DirSouth, 0, level.DirSouth},

		{level.DirNorth, 1, level.DirNorthEast},
		{level.DirNorth, -1, level.DirNorthWest},

		{level.DirEast, 2, level.DirSouth},
		{level.DirWest, 3, level.DirNorthEast},

		{level.DirNorth, 8, level.DirNorth},
		{level.DirNorth, -8, level.DirNorth},

		{level.DirEast, 16, level.DirEast},
		{level.DirSouthEast, 4, level.DirNorthWest},

		{level.DirSouth, -15, level.DirSouthWest},
	}

	for _, tc := range tt {
		result := tc.start.Offset(tc.offset)
		assert.Equal(t, tc.expected, result, fmt.Sprintf("Wrong for %d+%d", int(tc.start), tc.offset))
	}
}
