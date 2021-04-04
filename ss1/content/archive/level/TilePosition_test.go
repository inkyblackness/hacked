package level_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/inkyblackness/hacked/ss1/content/archive/level"
)

func TestFinePositionSlopeFactorFor(t *testing.T) {
	tt := []struct {
		pos      level.FinePosition
		tileType level.TileType
		expected float32
	}{
		{level.FinePosition{X: 127, Y: 127}, level.TileTypeOpen, 0.0},
		{level.FinePosition{X: 127, Y: 127}, level.TileTypeSlopeNorthToSouth, 0.5},
		{level.FinePosition{X: 127, Y: 127}, level.TileTypeRidgeNorthWestToSouthEast, 0.5},
		{level.FinePosition{X: 0, Y: 127}, level.TileTypeRidgeNorthWestToSouthEast, 0.0},
		{level.FinePosition{X: 127, Y: 0}, level.TileTypeRidgeNorthWestToSouthEast, 0.5},
	}

	for _, tc := range tt {
		result := tc.pos.SlopeFactorFor(tc.tileType)
		assert.InDeltaf(t, tc.expected, result, 0.01, "Type: %v, Pos: %v", tc.tileType, tc.pos)
	}
}
