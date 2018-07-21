package compression_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/inkyblackness/hacked/ss1/content/movie/compression"
)

type StandardTileColorerSuite struct {
	suite.Suite

	buffer []byte
	stride int

	colorer compression.TileColorFunction
}

func TestStandardTileColorerSuite(t *testing.T) {
	suite.Run(t, new(StandardTileColorerSuite))
}

func (suite *StandardTileColorerSuite) SetupTest() {
	suite.buffer = make([]byte, compression.PixelPerTile*9)
	suite.stride = compression.TileSideLength * 3

	for i := 0; i < len(suite.buffer); i++ {
		suite.buffer[i] = 0xDD
	}

	suite.colorer = compression.StandardTileColorer(suite.buffer, suite.stride)
}

func (suite *StandardTileColorerSuite) TestFunctionColorsAllPixelOfATile() {
	suite.whenTileIsColored(1, 1, []byte{0xFF, 0x00}, 0, 1)

	suite.thenTileShouldBe(1, 1, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})
}

func (suite *StandardTileColorerSuite) TestColoringATileLeavesPixelOfOthersAlone() {
	suite.whenTileIsColored(1, 1, []byte{0xFF}, 0, 1)

	suite.thenTileShouldBeUntouched(0, 0)
	suite.thenTileShouldBeUntouched(1, 0)
	suite.thenTileShouldBeUntouched(2, 0)
	suite.thenTileShouldBeUntouched(0, 1)
	//suite.thenTileShouldBeUntouched(1, 1)
	suite.thenTileShouldBeUntouched(2, 1)
	suite.thenTileShouldBeUntouched(0, 2)
	suite.thenTileShouldBeUntouched(1, 2)
	suite.thenTileShouldBeUntouched(2, 2)
}

func (suite *StandardTileColorerSuite) TestFunctionWorksWithMaximumIndexSize() {
	suite.whenTileIsColored(1, 1, []byte{0x01, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF},
		0xFEDCBA9876543210, 4)

	suite.thenTileShouldBe(1, 1, []byte{0x01, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF})
}

func (suite *StandardTileColorerSuite) TestFunctionSkipsZeroPaletteIndices() {
	suite.givenTileIsFilled(1, 1, 0xFF)

	suite.whenTileIsColored(1, 1, []byte{0x33, 0x00}, 0xAAAA, 1)

	suite.thenTileShouldBe(1, 1, []byte{0x33, 0xFF, 0x33, 0xFF, 0x33, 0xFF, 0x33, 0xFF, 0x33, 0xFF, 0x33, 0xFF, 0x33, 0xFF, 0x33, 0xFF})
}

func (suite *StandardTileColorerSuite) givenTileIsFilled(hTile int, vTile int, initValue byte) {
	suite.colorer(hTile, vTile, []byte{initValue, 0xEE}, 0x0000, 1)
}

func (suite *StandardTileColorerSuite) whenTileIsColored(hTile int, vTile int, lookupArray []byte, mask uint64, indexBitSize uint64) {
	suite.colorer(hTile, vTile, lookupArray, mask, indexBitSize)
}

func (suite *StandardTileColorerSuite) thenTileShouldBe(hTile, vTile int, expected []byte) {
	tile := suite.getTile(hTile, vTile)

	assert.Equal(suite.T(), expected, tile)
}

func (suite *StandardTileColorerSuite) thenTileShouldBeUntouched(hTile, vTile int) {
	expected := make([]byte, compression.PixelPerTile)
	for i := 0; i < len(expected); i++ {
		expected[i] = 0xDD
	}

	suite.thenTileShouldBe(hTile, vTile, expected)
}

func (suite *StandardTileColorerSuite) getTile(hTile, vTile int) []byte {
	tileBuffer := make([]byte, compression.PixelPerTile)
	start := vTile*compression.TileSideLength*suite.stride + hTile*compression.TileSideLength

	for i := 0; i < compression.TileSideLength; i++ {
		outOffset := compression.TileSideLength * i
		inOffset := start + suite.stride*i
		copy(tileBuffer[outOffset:outOffset+compression.TileSideLength], suite.buffer[inOffset:inOffset+compression.TileSideLength])
	}

	return tileBuffer
}
