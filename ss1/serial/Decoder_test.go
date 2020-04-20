package serial_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/inkyblackness/hacked/ss1/serial"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DecoderSuite struct {
	suite.Suite
	errorBuf *errorBuffer
	coder    *serial.Decoder
}

func TestDecoderSuite(t *testing.T) {
	suite.Run(t, new(DecoderSuite))
}

func (suite *DecoderSuite) SetupTest() {
	suite.errorBuf = nil
	suite.coder = nil
}

func (suite *DecoderSuite) TestImplementsCoderInterface() {
	instance := interface{}(suite.coder)
	_, ok := instance.(serial.Coder)
	assert.True(suite.T(), ok)
}

func (suite *DecoderSuite) TestCodeUint32() {
	suite.whenDecodingFrom([]byte{0x78, 0x56, 0x34, 0x12})

	var value uint32
	suite.coder.Code(&value)

	assert.Equal(suite.T(), uint32(0x12345678), value)
}

func (suite *DecoderSuite) TestCodeUint16() {
	suite.whenDecodingFrom([]byte{0x45, 0x23})

	var value uint16
	suite.coder.Code(&value)

	assert.Equal(suite.T(), uint16(0x2345), value)
}

func (suite *DecoderSuite) TestCodeByte() {
	suite.whenDecodingFrom([]byte{0xAB})

	var value byte
	suite.coder.Code(&value)

	assert.Equal(suite.T(), byte(0xAB), value)
}

func (suite *DecoderSuite) TestCodeByteSlice() {
	suite.whenDecodingFrom([]byte{0x78, 0x12, 0x34})

	value := make([]byte, 3)
	suite.coder.Code(value)

	assert.Equal(suite.T(), []byte{0x78, 0x12, 0x34}, value)
}

func (suite *DecoderSuite) TestFirstErrorReturnsFirstError() {
	suite.whenDecodingWithErrors()

	data := uint32(0)
	suite.coder.Code(&data)
	suite.errorBuf.errorOnNextCall = true
	suite.coder.Code(uint16(0))

	assert.EqualError(suite.T(), suite.coder.FirstError(), "errorBuffer on call number 2")
}

func (suite *DecoderSuite) TestFirstErrorIgnoresFurtherErrors() {
	suite.whenDecodingWithErrors()

	suite.errorBuf.errorOnNextCall = true
	suite.coder.Code(uint32(0))
	suite.errorBuf.errorOnNextCall = true
	suite.coder.Code(uint32(0))

	assert.Equal(suite.T(), suite.errorBuf.callCounter, 1)
	assert.EqualError(suite.T(), suite.coder.FirstError(), "errorBuffer on call number 1")
}

func (suite *DecoderSuite) TestImplementsReaderInterface() {
	instance := interface{}(suite.coder)
	_, ok := instance.(io.Reader)
	assert.True(suite.T(), ok)
}

func (suite *DecoderSuite) TestReadFromReader() {
	suite.whenDecodingFrom([]byte{0x78, 0x12, 0x34})

	value := make([]byte, 3)
	read, err := suite.coder.Read(value)

	assert.Nil(suite.T(), err, "no error expected")
	assert.Equal(suite.T(), 3, read, "unexpected length read")
	assert.Equal(suite.T(), []byte{0x78, 0x12, 0x34}, value)
}

func (suite *DecoderSuite) TestReadHandlesErrors() {
	suite.whenDecodingWithErrors()

	suite.errorBuf.errorOnNextCall = true
	_, _ = suite.coder.Read(make([]byte, 3))
	_, _ = suite.coder.Read(make([]byte, 5))
	assert.Equal(suite.T(), 1, suite.errorBuf.callCounter)
	assert.NotNil(suite.T(), suite.coder.FirstError())
}

func (suite *DecoderSuite) TestCodeWithCodable() {
	codable := new(MockedCodable)

	suite.whenDecodingFrom([]byte{0x78, 0x12, 0x34})
	suite.coder.Code(codable)

	assert.Equal(suite.T(), suite.coder, codable.calledCoder, "Should be called with coder")
}

func (suite *DecoderSuite) whenDecodingFrom(data []byte) {
	source := bytes.NewReader(data)
	suite.coder = serial.NewDecoder(source)
}

func (suite *DecoderSuite) whenDecodingWithErrors() {
	suite.errorBuf = new(errorBuffer)
	suite.coder = serial.NewDecoder(suite.errorBuf)
}
