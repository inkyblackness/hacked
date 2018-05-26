package compression

import (
	"bytes"
	"io"
	"math/rand"
	"testing"
	"time"

	"github.com/inkyblackness/hacked/ss1/serial"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
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
	suite.writeWords(0x0001, 0x0002, 0x0100, reset, 0x0003, 0x0004, 0x0100, endOfStream)

	suite.verifyOutput([]byte{0x01, 0x02, 0x01, 0x02, 0x03, 0x04, 0x03, 0x04})
}

func (suite *DecompressorSuite) TestDecompressHandlesSelfReferencingWords() {
	suite.writeWords(0x0001, 0x0002, 0x0101, endOfStream)

	suite.verifyOutput([]byte{0x01, 0x02, 0x02, 0x02})
}

func (suite *DecompressorSuite) writeWords(values ...word) {
	suite.store = serial.NewByteStore()
	coder := serial.NewEncoder(suite.store)
	writer := newWordWriter(coder)

	for _, value := range values {
		writer.write(value)
	}
	writer.close()
}

func (suite *DecompressorSuite) verify(input []byte) {
	suite.store = serial.NewByteStore()
	suite.compressor = NewCompressor(serial.NewEncoder(suite.store))

	suite.compressor.Write(input)
	suite.compressor.Close()

	suite.verifyOutput(input)
}

func (suite *DecompressorSuite) verifyOutput(expected []byte) {
	output := suite.buffer(len(expected))
	source := bytes.NewReader(suite.store.Data())
	decompressor := NewDecompressor(source)
	decompressor.Read(output)

	assert.Equal(suite.T(), expected, output)
}

func (suite *DecompressorSuite) buffer(byteCount int) []byte {
	result := make([]byte, byteCount)
	for i := range result {
		result[i] = 0xFF
	}
	return result
}
