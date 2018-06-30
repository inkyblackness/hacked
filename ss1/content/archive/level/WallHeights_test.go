package level_test

import (
	"fmt"
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type WallHeightsSuite struct {
	suite.Suite

	tileMap    level.TileMap
	heightsMap level.WallHeightsMap
}

func TestWallHeightsSuite(t *testing.T) {
	suite.Run(t, new(WallHeightsSuite))
}

func (suite *WallHeightsSuite) SetupTest() {
	suite.tileMap = level.NewTileMap(3, 3)
	suite.heightsMap = level.NewWallHeightsMap(3, 3)
}

func (suite *WallHeightsSuite) TestCalculateFromResetsToZeroForAllSolids() {
	suite.givenARandomHeightsMap()
	suite.whenHeightsAreCalculatedFromMap()
	suite.thenHeightsShouldBeForTile(1, 1, 0.0)
}

func (suite *WallHeightsSuite) TestCalculateFromSetsToMaximumHeightUnitFromOpenToSolid() {
	suite.givenTileHasType(1, 1, level.TileTypeOpen)
	suite.whenHeightsAreCalculatedFromMap()
	suite.thenHeightsShouldBeForTile(1, 1, 32.0)
}

func (suite *WallHeightsSuite) givenTileHasType(x int, y int, tileType level.TileType) {
	suite.tileMap.Tile(x, y).Type = tileType
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
	suite.heightsMap.CalculateFrom(suite.tileMap)
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
