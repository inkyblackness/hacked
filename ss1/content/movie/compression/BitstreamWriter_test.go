package compression_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/movie/compression"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BitstreamWriterSuite struct {
	suite.Suite

	writer compression.BitstreamWriter
}

func TestBitstreamWriter(t *testing.T) {
	suite.Run(t, new(BitstreamWriterSuite))
}

func (suite *BitstreamWriterSuite) SetupTest() {
	suite.writer = compression.BitstreamWriter{}
}

func (suite *BitstreamWriterSuite) TestWriteMoreThan32BitsPanics() {
	assert.Panics(suite.T(), func() { suite.whenWriting(33, 0) })
}

func (suite *BitstreamWriterSuite) TestWriteZeroBitsIsIgnored() {
	suite.whenWriting(0, 1)
	suite.thenBufferShouldBe([]byte{})
}

func (suite *BitstreamWriterSuite) TestWriteOneBit() {
	suite.whenWriting(1, 1)
	suite.thenBufferShouldBe([]byte{0x80})
}

func (suite *BitstreamWriterSuite) TestWriteEightBits() {
	suite.whenWriting(8, 0xEF)
	suite.thenBufferShouldBe([]byte{0xEF})
}

func (suite *BitstreamWriterSuite) TestWriteTwelveBits() {
	suite.whenWriting(12, 0xFEF)
	suite.thenBufferShouldBe([]byte{0xFE, 0xF0})
}

func (suite *BitstreamWriterSuite) TestWriteOneBitTwice() {
	suite.givenWritten(1, 1)
	suite.whenWriting(1, 1)
	suite.thenBufferShouldBe([]byte{0xC0})
}

func (suite *BitstreamWriterSuite) TestWriteArbitrary1() {
	suite.givenWritten(1, 1)
	suite.givenWritten(10, 0x2AA)
	suite.whenWriting(16, 0xFFFF)
	suite.thenBufferShouldBe([]byte{0xD5, 0x5F, 0xFF, 0xE0}) // [1 10 10101] [010 11111] [11111111] [111 00000]
}

func (suite *BitstreamWriterSuite) TestWriteArbitrary2() {
	suite.givenWritten(1, 1)
	suite.givenWritten(10, 0x2AA)
	suite.whenWriting(24, 0xFFFFFF)
	suite.thenBufferShouldBe([]byte{0xD5, 0x5F, 0xFF, 0xFF, 0xE0}) // extension of Arbitrary1
}

func (suite *BitstreamWriterSuite) TestWriteLessBitsThanGiven() {
	suite.givenWritten(4, 0x0)
	suite.whenWriting(8, 0xFFF)
	suite.thenBufferShouldBe([]byte{0x0F, 0xF0})
}

func (suite *BitstreamWriterSuite) TestWriteArbitrary3() {
	suite.givenWritten(1, 1)
	suite.givenWritten(10, 0x2AA)
	suite.whenWriting(32, 0xFFFFFFFF)
	suite.thenBufferShouldBe([]byte{0xD5, 0x5F, 0xFF, 0xFF, 0xFF, 0xE0}) // extension of Arbitrary1+2
}

func (suite *BitstreamWriterSuite) givenWritten(bits uint, value uint32) {
	suite.writer.Write(bits, value)
}

func (suite *BitstreamWriterSuite) whenWriting(bits uint, value uint32) {
	suite.writer.Write(bits, value)
}

func (suite *BitstreamWriterSuite) thenBufferShouldBe(expected []byte) {
	suite.T().Helper()
	buf := suite.writer.Buffer()
	if len(expected) == 0 {
		assert.Equal(suite.T(), 0, len(buf))
	} else {
		assert.Equal(suite.T(), expected, buf)
	}
}
