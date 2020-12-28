package compression_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/resource/lgres/internal/compression"
	"github.com/inkyblackness/hacked/ss1/serial"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WordWriterSuite struct {
	suite.Suite
	writer *compression.WordWriter
	store  *serial.ByteStore
}

func TestWordWriterSuite(t *testing.T) {
	suite.Run(t, new(WordWriterSuite))
}

func (suite *WordWriterSuite) SetupTest() {
	suite.store = serial.NewByteStore()
	coder := serial.NewEncoder(suite.store)
	suite.writer = compression.NewWordWriter(coder)
}

func (suite *WordWriterSuite) TestCloseWritesEndOfStreamMarkerAndTrailingZeroByte() {
	suite.writer.Close()

	assert.Equal(suite.T(), []byte{0xFF, 0xFC, 0x00}, suite.store.Data())
}

func (suite *WordWriterSuite) TestCloseWritesRemainderOnlyIfNotEmpty() {
	suite.writer.Write(compression.Word(0x0000))
	suite.writer.Write(compression.Word(0x0000))
	suite.writer.Write(compression.Word(0x0000))
	suite.writer.Close()

	assert.Equal(suite.T(), []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x3F, 0xFF, 0x00}, suite.store.Data())
}

func (suite *WordWriterSuite) TestWriteAndCloseLinesUpBits() {
	suite.writer.Write(compression.Word(0x1FFE)) // 0111111 1111110
	suite.writer.Close()                         // 1111111 1111111

	assert.Equal(suite.T(), []byte{0x7F, 0xFB, 0xFF, 0xF0, 0x00}, suite.store.Data())
}
