package level_test

import (
	"fmt"
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/archive/level"

	"github.com/stretchr/testify/assert"
)

func TestTileTypeInfoSlopeFactorsSealTileWithMirror(t *testing.T) {
	tt := []level.TileType{
		level.TileTypeSlopeSouthToNorth,
		level.TileTypeSlopeWestToEast,
		level.TileTypeSlopeNorthToSouth,
		level.TileTypeSlopeEastToWest,

		level.TileTypeValleySouthEastToNorthWest,
		level.TileTypeValleySouthWestToNorthEast,
		level.TileTypeValleyNorthWestToSouthEast,
		level.TileTypeValleyNorthEastToSouthWest,

		level.TileTypeRidgeNorthWestToSouthEast,
		level.TileTypeRidgeNorthEastToSouthWest,
		level.TileTypeRidgeSouthEastToNorthWest,
		level.TileTypeRidgeSouthWestToNorthEast,
	}
	for _, tc := range tt {
		info := tc.Info()
		mirroredInfo := info.SlopeMirrorType.Info()
		var factors level.SlopeFactors
		for i := 0; i < 8; i++ {
			factors[i] = info.SlopeFloorFactors[i] + mirroredInfo.SlopeFloorFactors[i]
		}
		assert.Equal(t, level.SlopeFactors{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0}, factors, fmt.Sprintf("Factors don't seal tile for type 0x%02X", tc))
	}
}

func TestTileTypeInfoDefaultsToSolidEquivalent(t *testing.T) {
	tileType := level.TileType(0x77)
	expected := level.TileTypeInfo{
		SolidSides:        level.DirNorth.Plus(level.DirEast).Plus(level.DirSouth).Plus(level.DirWest),
		SlopeFloorFactors: level.SlopeFactors{0, 0, 0, 0, 0, 0, 0, 0},
		SlopeMirrorType:   tileType,
	}
	assert.Equal(t, expected, tileType.Info())
}
