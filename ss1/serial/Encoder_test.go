package serial_test

import (
	"io"
	"testing"

	"github.com/inkyblackness/hacked/ss1/serial"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type EncoderSuite struct {
	suite.Suite
	store *serial.ByteStore
	coder *serial.Encoder
}

func TestEncoderSuite(t *testing.T) {
	suite.Run(t, new(EncoderSuite))
}

func (suite *EncoderSuite) SetupTest() {
	suite.store = serial.NewByteStore()
	suite.coder = serial.NewEncoder(suite.store)
}

func (suite *EncoderSuite) TestImplementsCoderInterface() {
	instance := interface{}(suite.coder)
	_, ok := instance.(serial.Coder)
	assert.True(suite.T(), ok)
}

func (suite *EncoderSuite) TestDataOnEmptyEncoderReturnsEmptyArray() {
	suite.verifyData(make([]byte, 0))
}

func (suite *EncoderSuite) TestCodeUint32Pointer() {
	var value uint32 = 0x12345678
	suite.verifyCodedData(&value, []byte{0x78, 0x56, 0x34, 0x12})
}

func (suite *EncoderSuite) TestCodeUint32() {
	var value uint32 = 0x12345678
	suite.verifyCodedData(value, []byte{0x78, 0x56, 0x34, 0x12})
}

func (suite *EncoderSuite) TestCodeUint16Pointer() {
	var value uint16 = 0x3456
	suite.verifyCodedData(&value, []byte{0x56, 0x34})
}

func (suite *EncoderSuite) TestCodeUint16() {
	var value uint16 = 0x3456
	suite.verifyCodedData(value, []byte{0x56, 0x34})
}

func (suite *EncoderSuite) TestCodeBytePointer() {
	var value byte = 0x42
	suite.verifyCodedData(&value, []byte{0x42})
}

func (suite *EncoderSuite) TestCodeByte() {
	var value byte = 0x42
	suite.verifyCodedData(value, []byte{0x42})
}

func (suite *EncoderSuite) TestCodeByteSlice() {
	value := []byte{0x01, 0x02, 0x03}
	suite.verifyCodedData(value, []byte{0x01, 0x02, 0x03})
}

func (suite *EncoderSuite) TestFirstErrorByCode() {
	var target errorBuffer
	suite.coder = serial.NewEncoder(&target)
	suite.coder.Code(uint16(0))
	target.errorOnNextCall = true
	suite.coder.Code(uint32(0))

	assert.EqualError(suite.T(), suite.coder.FirstError(), "errorBuffer on call number 2")
}

func (suite *EncoderSuite) TestCodeDoesNothingOnPreviousError() {
	var target errorBuffer
	suite.coder = serial.NewEncoder(&target)
	target.errorOnNextCall = true
	suite.coder.Code(uint32(0))
	target.errorOnNextCall = true
	suite.coder.Code(uint32(0))

	assert.Equal(suite.T(), target.callCounter, 1)
	assert.EqualError(suite.T(), suite.coder.FirstError(), "errorBuffer on call number 1")
}

func (suite *EncoderSuite) TestImplementsWriterInterface() {
	instance := interface{}(suite.coder)
	_, ok := instance.(io.Writer)
	assert.True(suite.T(), ok)
}

func (suite *EncoderSuite) TestWriteToWriter() {
	written, err := suite.coder.Write([]byte{0x0A, 0x0B, 0x0C})
	assert.Nil(suite.T(), err, "no error expected")
	assert.Equal(suite.T(), 3, written, "unexpected length written")
	suite.verifyData([]byte{0x0A, 0x0B, 0x0C})
}

func (suite *EncoderSuite) TestWriteDoesNothingOnPreviousError() {
	var target errorBuffer
	suite.coder = serial.NewEncoder(&target)
	target.errorOnNextCall = true
	suite.coder.Code(uint32(0))
	_, _ = suite.coder.Write([]byte{0x0A, 0x0B, 0x0C})

	assert.Equal(suite.T(), target.callCounter, 1)
}

func (suite *EncoderSuite) TestWriteReturnsFirstErrorOnPreviousError() {
	var target errorBuffer
	suite.coder = serial.NewEncoder(&target)
	target.errorOnNextCall = true
	suite.coder.Code(uint32(0))
	n, err := suite.coder.Write([]byte{0x0A, 0x0B, 0x0C})

	assert.Equal(suite.T(), 0, n, "Should return zero written bytes")
	assert.EqualError(suite.T(), err, "errorBuffer on call number 1")
}

func (suite *EncoderSuite) TestWriteOnErrorCase() {
	var target errorBuffer
	suite.coder = serial.NewEncoder(&target)
	target.errorOnNextCall = true
	target.errorByteCount = 2
	n, err := suite.coder.Write([]byte{0x0A, 0x0B, 0x0C})

	assert.Equal(suite.T(), 2, n, "Should return count of successfully written bytes")
	assert.EqualError(suite.T(), err, "errorBuffer on call number 1")
	assert.EqualError(suite.T(), suite.coder.FirstError(), "errorBuffer on call number 1", "should set first error")
}

func (suite *EncoderSuite) TestCodeWithCodable() {
	codable := new(MockedCodable)
	suite.coder.Code(codable)

	assert.Equal(suite.T(), suite.coder, codable.calledCoder, "Should be called with coder")
}

func (suite *EncoderSuite) verifyData(expected []byte) {
	assert.Nil(suite.T(), suite.coder.FirstError())
	result := suite.store.Data()
	assert.Equal(suite.T(), expected, result)
}

func (suite *EncoderSuite) verifyCodedData(value interface{}, expected []byte) {
	suite.coder.Code(value)
	suite.verifyData(expected)
}
