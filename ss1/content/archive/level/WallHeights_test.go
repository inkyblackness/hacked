package level_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/inkyblackness/hacked/ss1/content/archive/level"
)

type WallHeightsSuite struct {
	suite.Suite

	tileMap    level.TileMap
	heightsMap level.WallHeightsMap

	tile *level.TileMapEntry
}

func TestWallHeightsSuite(t *testing.T) {
	suite.Run(t, new(WallHeightsSuite))
}

func (suite *WallHeightsSuite) SetupTest() {
	suite.tileMap = level.NewTileMap(4, 4)
	suite.heightsMap = level.NewWallHeightsMap(3, 3)
	suite.tile = nil
}

func (suite *WallHeightsSuite) TestCalculateFromResetsToZeroForAllSolids() {
	suite.givenARandomHeightsMap()
	suite.whenHeightsAreCalculatedFromMap()
	suite.thenHeightsShouldBeForTile(1, 1, 0.0)
}

func (suite *WallHeightsSuite) TestCalculateFromSetsToMaximumHeightUnitFromOpenToSolid() {
	suite.givenTile(1, 1).typed(level.TileTypeOpen)
	suite.whenHeightsAreCalculatedFromMap()
	suite.thenHeightsShouldBeForTile(1, 1, 32.0)
	suite.thenHeightShouldBeForTileAt(1, 2, level.DirSouth, 32, 32, 32)
	suite.thenHeightShouldBeForTileAt(1, 0, level.DirNorth, 32, 32, 32)
	suite.thenHeightShouldBeForTileAt(0, 1, level.DirEast, 32, 32, 32)
	suite.thenHeightShouldBeForTileAt(2, 1, level.DirWest, 32, 32, 32)
}

func (suite *WallHeightsSuite) TestCalculateFromSetsToMaximumHeightUnitFromOpenToSolidSide() {
	suite.givenTile(1, 1).typed(level.TileTypeOpen)
	suite.givenTile(2, 1).typed(level.TileTypeDiagonalOpenNorthEast)
	suite.whenHeightsAreCalculatedFromMap()
	suite.thenHeightsShouldBeForTile(1, 1, 32)
	suite.thenHeightShouldBeForTileAt(2, 1, level.DirWest, 32, 32, 32)
}

func (suite *WallHeightsSuite) TestCalculateFromSetsHeightDifferencesForFloors() {
	suite.givenTile(1, 1).typed(level.TileTypeOpen).floorAt(0)
	suite.givenTile(1, 2).typed(level.TileTypeOpen).floorAt(10)
	suite.whenHeightsAreCalculatedFromMap()
	suite.thenHeightShouldBeForTileAt(1, 1, level.DirNorth, 10, 10, 10)
	suite.thenHeightShouldBeForTileAt(1, 2, level.DirSouth, -10, -10, -10)
}

func (suite *WallHeightsSuite) TestCalculateFromSetsHeightDifferencesForFloorsRespectingSlope() {
	suite.givenTile(1, 1).typed(level.TileTypeSlopeEastToWest).floorAt(0).sloped(level.TileSlopeControlCeilingFlat, 2)
	suite.givenTile(1, 2).typed(level.TileTypeOpen).floorAt(10)
	suite.whenHeightsAreCalculatedFromMap()
	suite.thenHeightShouldBeForTileAt(1, 1, level.DirNorth, 8, 9, 10)
	suite.thenHeightShouldBeForTileAt(1, 2, level.DirSouth, -10, -9, -8)
}

func (suite *WallHeightsSuite) TestCalculateFromSetsHeightDifferencesForFloorsRespectingSlopeSameFloorHeight() {
	suite.givenTile(1, 1).typed(level.TileTypeOpen).floorAt(0)
	suite.givenTile(1, 2).typed(level.TileTypeSlopeEastToWest).floorAt(0).sloped(level.TileSlopeControlCeilingFlat, 4)
	suite.whenHeightsAreCalculatedFromMap()
	suite.thenHeightShouldBeForTileAt(1, 1, level.DirNorth, 4, 2, 0)
	suite.thenHeightShouldBeForTileAt(1, 2, level.DirSouth, 0, -2, -4)
}

func (suite *WallHeightsSuite) TestCalculateFromConsidersCeilingHeightsIfWallingOff() {
	suite.givenTile(1, 1).typed(level.TileTypeOpen).floorAt(0)
	suite.givenTile(1, 2).typed(level.TileTypeSlopeSouthToNorth).floorAt(0).sloped(level.TileSlopeControlFloorFlat, 8).ceilingAt(8)
	suite.whenHeightsAreCalculatedFromMap()
	suite.thenHeightShouldBeForTileAt(1, 1, level.DirNorth, 32, 32, 32)
	suite.thenHeightShouldBeForTileAt(1, 2, level.DirSouth, 32, 32, 32)
}

func (suite *WallHeightsSuite) TestCalculateFromConsidersMirroredSlopesWallingOff() {
	suite.givenTile(1, 1).typed(level.TileTypeOpen).floorAt(0)
	suite.givenTile(1, 2).typed(level.TileTypeSlopeNorthToSouth).floorAt(0).sloped(level.TileSlopeControlCeilingMirrored, 16)
	suite.whenHeightsAreCalculatedFromMap()
	suite.thenHeightShouldBeForTileAt(1, 1, level.DirNorth, 32, 32, 32)
	suite.thenHeightShouldBeForTileAt(1, 2, level.DirSouth, 32, 32, 32)
}

func (suite *WallHeightsSuite) TestCalculateFromWithFloorCrossingOtherCeiling() {
	suite.givenTile(1, 1).typed(level.TileTypeOpen).floorAt(20)
	suite.givenTile(1, 2).typed(level.TileTypeOpen).floorAt(0).ceilingAt(8)
	suite.whenHeightsAreCalculatedFromMap()
	suite.thenHeightShouldBeForTileAt(1, 1, level.DirNorth, 32, 32, 32)
	suite.thenHeightShouldBeForTileAt(1, 2, level.DirSouth, 32, 32, 32)
}

func (suite *WallHeightsSuite) TestCalculateFromWithSameSlopes() {
	suite.givenTile(1, 1).typed(level.TileTypeSlopeEastToWest).sloped(level.TileSlopeControlCeilingInverted, 8)
	suite.givenTile(1, 2).typed(level.TileTypeSlopeEastToWest).sloped(level.TileSlopeControlCeilingInverted, 8)
	suite.whenHeightsAreCalculatedFromMap()
	suite.thenHeightShouldBeForTileAt(1, 1, level.DirNorth, 0, 0, 0)
	suite.thenHeightShouldBeForTileAt(1, 2, level.DirSouth, 0, 0, 0)
}

func (suite *WallHeightsSuite) givenTile(x, y int) *WallHeightsSuite {
	suite.tile = suite.Tile(x, y)
	return suite
}

func (suite *WallHeightsSuite) typed(tileType level.TileType) *WallHeightsSuite {
	suite.tile.Type = tileType
	return suite
}

func (suite *WallHeightsSuite) sloped(ctrl level.TileSlopeControl, height level.TileHeightUnit) *WallHeightsSuite {
	suite.tile.Flags = suite.tile.Flags.WithSlopeControl(ctrl)
	suite.tile.SlopeHeight = height
	return suite
}

func (suite *WallHeightsSuite) floorAt(value level.TileHeightUnit) *WallHeightsSuite {
	suite.tile.Floor = suite.tile.Floor.WithAbsoluteHeight(value)
	return suite
}

func (suite *WallHeightsSuite) ceilingAt(value level.TileHeightUnit) *WallHeightsSuite {
	suite.tile.Ceiling = suite.tile.Ceiling.WithAbsoluteHeight(value)
	return suite
}

func (suite *WallHeightsSuite) givenARandomHeightsMap() {
	var random float32
	randomizeSide := func(side *[3]float32) {
		for i := 0; i < 3; i++ {
			random += 0.5
			if random > 32 {
				random = -32
			}
			side[i] = random
		}
	}
	randomizeTile := func(heights *level.WallHeights) {
		randomizeSide(&heights.North)
		randomizeSide(&heights.East)
		randomizeSide(&heights.South)
		randomizeSide(&heights.West)
	}

	for _, row := range suite.heightsMap {
		for i := 0; i < len(row); i++ {
			heights := &row[i]
			randomizeTile(heights)
		}
	}
}

func (suite *WallHeightsSuite) whenHeightsAreCalculatedFromMap() {
	suite.heightsMap.CalculateFrom(suite)
}

func (suite *WallHeightsSuite) thenHeightsShouldBeForTile(x int, y int, expected float32) {
	tile := suite.heightsMap.Tile(x, y)

	verifySide := func(name string, values [3]float32) {
		for i := 0; i < len(values); i++ {
			assert.Equal(suite.T(), expected, values[i], fmt.Sprintf("%s[%d] mismatch", name, i))
		}
	}
	verifySide("North", tile.North)
	verifySide("East", tile.East)
	verifySide("South", tile.South)
	verifySide("West", tile.West)
}

func (suite *WallHeightsSuite) thenHeightShouldBeForTileAt(x int, y int, side level.Direction,
	expectedLeft float32, expectedCenter float32, expectedRight float32) {
	tile := suite.heightsMap.Tile(x, y)
	var values [3]float32

	switch side {
	case level.DirNorth:
		values = tile.North
	case level.DirEast:
		values = tile.East
	case level.DirSouth:
		values = tile.South
	case level.DirWest:
		values = tile.West
	default:
		require.Fail(suite.T(), "Invalid side specified")
	}

	assert.Equal(suite.T(), expectedLeft, values[0], "left mismatch")
	assert.Equal(suite.T(), expectedCenter, values[1], "center mismatch")
	assert.Equal(suite.T(), expectedRight, values[2], "right mismatch")
}

func (suite *WallHeightsSuite) Tile(x, y int) *level.TileMapEntry {
	return suite.tileMap.Tile(x, y, 2)
}
