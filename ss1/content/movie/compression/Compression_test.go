package compression_test

import (
	"fmt"
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/movie/compression"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CompressionSuite struct {
	suite.Suite
}

func TestTileSideLength(t *testing.T) {
	assert.Equal(t, 4, compression.TileSideLength, "The following tests assume a tile side length of 4")
}

func TestCompressionTest(t *testing.T) {
	suite.Run(t, new(CompressionSuite))
}

func (suite *CompressionSuite) SetupTest() {

}

func (suite *CompressionSuite) TestOneFrameOneTile() {
	frame0 := []byte{
		0x01, 0x02, 0x03, 0x04,
		0x05, 0x06, 0x07, 0x08,
		0x09, 0x0A, 0x0B, 0x0C,
		0x0D, 0x0E, 0x0F, 0x10,
	}
	suite.verifyCompression(4, 4, frame0)
}

func (suite *CompressionSuite) TestTwoFramesOneTile() {
	frame0 := []byte{
		0x01, 0x02, 0x03, 0x04,
		0x05, 0x06, 0x07, 0x08,
		0x09, 0x0A, 0x0B, 0x0C,
		0x0D, 0x0E, 0x0F, 0x10,
	}
	frame1 := []byte{
		0x20, 0x21, 0x22, 0x23,
		0x24, 0x25, 0x26, 0x27,
		0x28, 0x29, 0x2A, 0x2B,
		0x2C, 0x2D, 0x2E, 0x2F,
	}
	suite.verifyCompression(4, 4, frame0, frame1)
}

/*
func (suite *CompressionSuite) TestTwoTiles() {
	frame0 := []byte{
		0x01, 0x02, 0x03, 0x04, 0x20, 0x21, 0x22, 0x23,
		0x05, 0x06, 0x07, 0x08, 0x24, 0x25, 0x26, 0x27,
		0x09, 0x0A, 0x0B, 0x0C, 0x28, 0x29, 0x2A, 0x2B,
		0x0D, 0x0E, 0x0F, 0x10, 0x2C, 0x2D, 0x2E, 0x2F,
	}
	suite.verifyCompression(8, 4, frame0)
}
*/
func (suite *CompressionSuite) verifyCompression(width, height int, inFrames ...[]byte) {
	suite.T().Helper()
	encoder := compression.NewSceneEncoder(width, height)
	for frameIndex, frame := range inFrames {
		require.Equal(suite.T(), width*height, len(frame), fmt.Sprintf("Length of frame %d is wrong for dimension", frameIndex))
		err := encoder.AddFrame(frame)
		require.Nil(suite.T(), err, fmt.Sprintf("no error expected adding frame %d: %v", frameIndex, err))
	}
	controlWords, paletteLookup, encodedFrames, err := encoder.Encode()
	require.Equal(suite.T(), len(inFrames), len(encodedFrames), "expected equal amount of encoded frames for input frames")
	require.Nil(suite.T(), err, fmt.Sprintf("no error expected encoding: %v", err))

	decoderBuilder := compression.NewFrameDecoderBuilder(width, height)
	decoderBuilder.WithControlWords(controlWords)
	decoderBuilder.WithPaletteLookupList(paletteLookup)
	frameBuffer := make([]byte, width*height)
	decoderBuilder.ForStandardFrame(frameBuffer, width)
	for frameIndex, encodedFrame := range encodedFrames {
		decoder := decoderBuilder.Build()
		decoder.Decode(encodedFrame.Bitstream, encodedFrame.Maskstream)
		assert.Equal(suite.T(), inFrames[frameIndex], frameBuffer, fmt.Sprintf("Frame decode error for frame %d", frameIndex))
	}
}
