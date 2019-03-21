package compression_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/inkyblackness/hacked/ss1/content/movie/internal/compression"

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

func (suite *CompressionSuite) TestTwoTiles() {
	frame0 := []byte{
		0x01, 0x02, 0x03, 0x04, 0x20, 0x21, 0x22, 0x23,
		0x05, 0x06, 0x07, 0x08, 0x24, 0x25, 0x26, 0x27,
		0x09, 0x0A, 0x0B, 0x0C, 0x28, 0x29, 0x2A, 0x2B,
		0x0D, 0x0E, 0x0F, 0x10, 0x2C, 0x2D, 0x2E, 0x2F,
	}
	suite.verifyCompression(8, 4, frame0)
}

func (suite *CompressionSuite) verifyCompression(width, height int, inFrames ...[]byte) {
	suite.T().Helper()
	verifyCompression(suite.T(), width, height, inFrames...)
}

func BenchmarkRandomFrames(b *testing.B) {
	// With a frame size of 600x300 we have 150*75=11250 tiles per frame.
	// The maximum palette offset is 0x1FFFF, allowing for 8191 entries of full 16 byte (fully random tile) palettes.
	// This means not even one frame of completely randomized tiles can be encoded.
	// As a result, this test will first create a set of 4000 random tiles, and then create a sequence of frames
	// using this set of tiles, instead of creating complete random frames.
	// Still, this test is rather slow, so it is moved to be a benchmark for now.

	seed := time.Now().UnixNano()
	b.Logf("Running with seed %v", seed)
	hTiles := 150
	vTiles := 75
	width := 4 * hTiles
	height := 4 * vTiles
	r := rand.New(rand.NewSource(seed))
	tiles := make([][4][4]byte, 4000)
	for tileIndex := 0; tileIndex < len(tiles); tileIndex++ {
		tile := &tiles[tileIndex]
		for y := 0; y < 4; y++ {
			for x := 0; x < 4; x++ {
				tile[y][x] = byte(1 + r.Intn(255)) // value of 0x00 may not be encoded.
			}
		}
	}

	frames := make([][]byte, 5)
	for frameIndex := 0; frameIndex < len(frames); frameIndex++ {
		frame := make([]byte, width*height)
		for y := 0; y < height; y += 4 {
			for x := 0; x < width; x += 4 {
				tile := &tiles[r.Intn(len(tiles))]
				for line := 0; line < 4; line++ {
					copy(frame[(y+line)*width+x:], tile[line][:])
				}
			}
		}
		frames[frameIndex] = frame
	}

	b.ResetTimer()
	verifyCompression(b, width, height, frames...)
}

func verifyCompression(t testing.TB, width, height int, inFrames ...[]byte) {
	t.Helper()
	encoder := compression.NewSceneEncoder(width, height)
	for frameIndex, frame := range inFrames {
		require.Equal(t, width*height, len(frame), fmt.Sprintf("Length of frame %d is wrong for dimension", frameIndex))
		err := encoder.AddFrame(frame)
		require.Nil(t, err, fmt.Sprintf("no error expected adding frame %d: %v", frameIndex, err))
	}
	controlWords, paletteLookup, encodedFrames, err := encoder.Encode()
	require.Equal(t, len(inFrames), len(encodedFrames), "expected equal amount of encoded frames for input frames")
	require.Nil(t, err, fmt.Sprintf("no error expected encoding: %v", err))

	t.Logf("Statistics: ControlWords: %v, PaletteLookup: %v", len(controlWords), len(paletteLookup))

	decoderBuilder := compression.NewFrameDecoderBuilder(width, height)
	decoderBuilder.WithControlWords(controlWords)
	decoderBuilder.WithPaletteLookupList(paletteLookup)
	frameBuffer := make([]byte, width*height)
	decoderBuilder.ForStandardFrame(frameBuffer, width)
	for frameIndex, encodedFrame := range encodedFrames {
		decoder := decoderBuilder.Build()
		t.Logf("Statistics F%2d: Bitstream: %v, Maskstream: %v",
			frameIndex, len(encodedFrame.Bitstream), len(encodedFrame.Maskstream))
		decoder.Decode(encodedFrame.Bitstream, encodedFrame.Maskstream)
		assert.Equal(t, inFrames[frameIndex], frameBuffer, fmt.Sprintf("Frame content mismatch for frame %d", frameIndex))
	}
}
