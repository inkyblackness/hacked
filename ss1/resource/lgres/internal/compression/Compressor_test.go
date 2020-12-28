package compression_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/inkyblackness/hacked/ss1/resource/lgres/internal/compression"
	"github.com/inkyblackness/hacked/ss1/serial"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CompressorSuite struct {
	suite.Suite
	store      *serial.ByteStore
	compressor io.WriteCloser
}

func TestCompressorSuite(t *testing.T) {
	suite.Run(t, new(CompressorSuite))
}

func (suite *CompressorSuite) SetupTest() {
	suite.store = serial.NewByteStore()
	suite.compressor = compression.NewCompressor(suite.store)
}

func (suite *CompressorSuite) TestWriteCompressesFirstReocurrence() {
	suite.givenWrittenData([]byte{0x00, 0x01})
	suite.givenWrittenData([]byte{0x00, 0x01})

	err := suite.compressor.Close()
	assert.Nil(suite.T(), err, "no error expected")

	suite.thenWordsShouldBe(compression.Word(0x0000), compression.Word(0x0001), compression.Word(0x0100))
}

func (suite *CompressorSuite) TestWriteCompressesTest1() {
	suite.givenWrittenData([]byte{0x00, 0x01, 0x00, 0x02, 0x01, 0x00, 0x01})

	err := suite.compressor.Close()
	assert.Nil(suite.T(), err, "no error expected")

	suite.thenWordsShouldBe(
		compression.Word(0x0000), compression.Word(0x0001), compression.Word(0x0000),
		compression.Word(0x0002), compression.Word(0x0101), compression.Word(0x0001))
}

func (suite *CompressorSuite) givenWrittenData(data []byte) {
	written, err := suite.compressor.Write(data)
	require.Nil(suite.T(), err, "no error expected writing")
	require.Equal(suite.T(), len(data), written, "expected all bytes to be written")
}

func (suite *CompressorSuite) thenWordsShouldBe(expected ...compression.Word) {
	source := bytes.NewReader(suite.store.Data())
	reader := compression.NewWordReader(serial.NewDecoder(source))
	var words []compression.Word

	for read := reader.Read(); read != compression.EndOfStream; read = reader.Read() {
		words = append(words, read)
	}

	assert.Equal(suite.T(), expected, words)
}
