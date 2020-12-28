package compression

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/serial"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WordWriterSuite struct {
	suite.Suite
	writer *wordWriter
	store  *serial.ByteStore
}

func TestWordWriterSuite(t *testing.T) {
	suite.Run(t, new(WordWriterSuite))
}

func (suite *WordWriterSuite) SetupTest() {
	suite.store = serial.NewByteStore()
	coder := serial.NewEncoder(suite.store)
	suite.writer = newWordWriter(coder)
}

func (suite *WordWriterSuite) TestCloseWritesEndOfStreamMarkerAndTrailingZeroByte() {
	suite.writer.close()

	assert.Equal(suite.T(), []byte{0xFF, 0xFC, 0x00}, suite.store.Data())
}

func (suite *WordWriterSuite) TestCloseWritesRemainderOnlyIfNotEmpty() {
	suite.writer.write(word(0x0000))
	suite.writer.write(word(0x0000))
	suite.writer.write(word(0x0000))
	suite.writer.close()

	assert.Equal(suite.T(), []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x3F, 0xFF, 0x00}, suite.store.Data())
}

func (suite *WordWriterSuite) TestWriteAndCloseLinesUpBits() {
	suite.writer.write(word(0x1FFE)) // 0111111 1111110
	suite.writer.close()             // 1111111 1111111

	assert.Equal(suite.T(), []byte{0x7F, 0xFB, 0xFF, 0xF0, 0x00}, suite.store.Data())
}
