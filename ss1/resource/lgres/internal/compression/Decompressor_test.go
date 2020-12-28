package compression_test

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/inkyblackness/hacked/ss1/resource/lgres/internal/compression"
	"github.com/inkyblackness/hacked/ss1/serial"
)

type DecompressorSuite struct {
	suite.Suite
	store      *serial.ByteStore
	compressor io.WriteCloser
}

func TestDecompressorSuite(t *testing.T) {
	suite.Run(t, new(DecompressorSuite))
}

func (suite *DecompressorSuite) SetupTest() {
}

func (suite *DecompressorSuite) TestDecompressTest1() {
	input := []byte{0x00, 0x01, 0x00, 0x01}

	suite.verify(input)
}

func (suite *DecompressorSuite) TestDecompressTest2() {
	input := []byte{0x00, 0x01, 0x00, 0x01, 0x00, 0x01}

	suite.verify(input)
}

func (suite *DecompressorSuite) TestDecompressTest3() {
	suite.verify([]byte{})
}

func (suite *DecompressorSuite) TestDecompressTest4() {
	input := []byte{0x00, 0x01, 0x00, 0x02, 0x01, 0x00, 0x01}

	suite.verify(input)
}

func (suite *DecompressorSuite) TestDecompressTest5() {
	input := []byte{0x00, 0x01, 0x00, 0x02, 0x01, 0x00, 0x01, 0x02, 0x01, 0x02, 0x01, 0x00, 0x01, 0x02}

	suite.verify(input)
}

func (suite *DecompressorSuite) TestDecompressTestRandom() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for testCase := 0; testCase < 100; testCase++ {
		input := make([]byte, r.Intn(1024))
		for i := 0; i < len(input); i++ {
			input[i] = byte(r.Intn(256))
		}
		suite.verify(input)
	}
}

func (suite *DecompressorSuite) TestDecompressHandlesDictionaryResets() {
	suite.writeWords(0x0001, 0x0002, 0x0100, compression.Reset, 0x0003, 0x0004, 0x0100, compression.EndOfStream)

	suite.verifyOutput([]byte{0x01, 0x02, 0x01, 0x02, 0x03, 0x04, 0x03, 0x04})
}

func (suite *DecompressorSuite) TestDecompressHandlesSelfReferencingWords() {
	suite.writeWords(0x0001, 0x0002, 0x0101, compression.EndOfStream)

	suite.verifyOutput([]byte{0x01, 0x02, 0x02, 0x02})
}

func (suite *DecompressorSuite) writeWords(values ...compression.Word) {
	suite.store = serial.NewByteStore()
	coder := serial.NewEncoder(suite.store)
	writer := compression.NewWordWriter(coder)

	for _, value := range values {
		writer.Write(value)
	}
	writer.Close()
}

func (suite *DecompressorSuite) verify(input []byte) {
	suite.store = serial.NewByteStore()
	suite.compressor = compression.NewCompressor(serial.NewEncoder(suite.store))

	n, err := suite.compressor.Write(input)
	assert.Nil(suite.T(), err, "no error expected writing")
	assert.Equal(suite.T(), len(input), n, "invalid number written")
	err = suite.compressor.Close()
	assert.Nil(suite.T(), err, "no error expected closing")

	suite.verifyOutput(input)
}

func (suite *DecompressorSuite) verifyOutput(expected []byte) {
	output := suite.buffer(len(expected))
	source := bytes.NewReader(suite.store.Data())
	decompressor := compression.NewDecompressor(source)
	read, err := decompressor.Read(output)

	assert.True(suite.T(), (err == nil) || (err == io.EOF), fmt.Sprintf("unexpected error: %v", err))
	assert.Equal(suite.T(), len(output), read, "unexpected number of bytes read")
	assert.Equal(suite.T(), expected, output)
}

func (suite *DecompressorSuite) buffer(byteCount int) []byte {
	result := make([]byte, byteCount)
	for i := range result {
		result[i] = 0xFF
	}
	return result
}
