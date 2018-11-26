package serial_test

import (
	"io"
	"testing"

	"github.com/inkyblackness/hacked/ss1/serial"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PositioningEncoderSuite struct {
	suite.Suite
	coder *serial.PositioningEncoder
	store *serial.ByteStore
}

func TestPositioningEncoderSuite(t *testing.T) {
	suite.Run(t, new(PositioningEncoderSuite))
}

func (suite *PositioningEncoderSuite) SetupTest() {
	suite.store = serial.NewByteStore()
	suite.coder = serial.NewPositioningEncoder(suite.store)
}

func (suite *PositioningEncoderSuite) TestImplementsPositioningCoderInterface() {
	instance := interface{}(suite.coder)
	_, ok := instance.(serial.PositioningCoder)

	assert.True(suite.T(), ok)
}

func (suite *PositioningEncoderSuite) TestCurPosReturnsCurrentOffset() {
	suite.coder.Code(uint32(0))
	assert.Equal(suite.T(), uint32(4), suite.coder.CurPos())
}

func (suite *PositioningEncoderSuite) TestSetCurPosRepositionsWritePointer() {
	suite.coder.Code(uint32(0))
	suite.coder.SetCurPos(0)
	suite.coder.Code(uint32(0x13243546))
	result := suite.store.Data()

	assert.Equal(suite.T(), []byte{0x46, 0x35, 0x24, 0x13}, result)
}

func (suite *PositioningEncoderSuite) TestFirstErrorBySetCurPos() {
	var target errorBuffer
	suite.coder = serial.NewPositioningEncoder(&target)

	suite.coder.Code(uint32(0))
	target.errorOnNextCall = true
	suite.coder.SetCurPos(0)

	assert.EqualError(suite.T(), suite.coder.FirstError(), "errorBuffer on call number 2")
}

func (suite *PositioningEncoderSuite) TestSetCurPosDoesNothingOnPreviousError() {
	var target errorBuffer
	suite.coder = serial.NewPositioningEncoder(&target)
	target.errorOnNextCall = true
	suite.coder.Code(uint32(0))
	target.errorOnNextCall = true
	suite.coder.SetCurPos(0)

	assert.Equal(suite.T(), target.callCounter, 1)
	assert.EqualError(suite.T(), suite.coder.FirstError(), "errorBuffer on call number 1")
}

func (suite *PositioningEncoderSuite) TestCurPosIsNotChangedOnSetError() {
	var target errorBuffer
	suite.coder = serial.NewPositioningEncoder(&target)
	target.errorOnNextCall = true
	suite.coder.SetCurPos(4)
	assert.Equal(suite.T(), uint32(0), suite.coder.CurPos())
}

func (suite *PositioningEncoderSuite) TestCurPosIsRelativeToConstructionTime() {
	_, _ = suite.store.Write([]byte{0x01, 0x02, 0x03})

	assert.Equal(suite.T(), uint32(0), suite.coder.CurPos(), "CurPos should be zero at the start")
	suite.coder.Code(uint16(0xAAAA))
	suite.coder.SetCurPos(1)
	suite.coder.Code(byte(0xCC))

	result := suite.store.Data()
	assert.Equal(suite.T(), []byte{0x01, 0x02, 0x03, 0xAA, 0xCC}, result, "SetCurPos did not work relative to original start")
}

func (suite *PositioningEncoderSuite) TestImplementsSeekerInterface() {
	instance := interface{}(suite.coder)
	_, ok := instance.(io.Seeker)

	assert.True(suite.T(), ok)
}

func (suite *PositioningEncoderSuite) TestSeekReturnsErrorForWrongWhence() {
	pos, err := suite.coder.Seek(1, 1234)

	assert.Equal(suite.T(), int64(0), pos, "zero position should be returned")
	assert.EqualError(suite.T(), err, "seek: invalid whence")
}

func (suite *PositioningEncoderSuite) TestSeekAbsolute() {
	suite.coder.Code(uint32(0))
	pos, err := suite.coder.Seek(1, io.SeekStart)
	assert.Nil(suite.T(), err, "no error expected")
	assert.Equal(suite.T(), int64(1), pos, "invalid position returned")

	suite.coder.Code(uint16(0xAABB))
	result := suite.store.Data()
	assert.Equal(suite.T(), []byte{0x00, 0xBB, 0xAA, 0x00}, result)
	assert.Equal(suite.T(), uint32(3), suite.coder.CurPos(), "should update current position")
}

func (suite *PositioningEncoderSuite) TestSeekRelative() {
	suite.coder.Code(uint32(0))
	pos, err := suite.coder.Seek(-1, io.SeekCurrent)
	assert.Equal(suite.T(), int64(3), pos, "should return new position")
	assert.Nil(suite.T(), err, "should have no error")

	suite.coder.Code(uint16(0xAABB))

	result := suite.store.Data()
	assert.Equal(suite.T(), []byte{0x00, 0x00, 0x00, 0xBB, 0xAA}, result)
	assert.Equal(suite.T(), uint32(5), suite.coder.CurPos(), "should update current position")
}

func (suite *PositioningEncoderSuite) TestSeekDoesNothingOnPreviousError() {
	var target errorBuffer
	suite.coder = serial.NewPositioningEncoder(&target)
	target.errorOnNextCall = true
	suite.coder.Code(uint32(0))
	_, _ = suite.coder.Seek(0, io.SeekStart)

	assert.Equal(suite.T(), target.callCounter, 1)
	assert.EqualError(suite.T(), suite.coder.FirstError(), "errorBuffer on call number 1")
}

func (suite *PositioningEncoderSuite) TestSeekRegistersReturnedError() {
	var target errorBuffer
	suite.coder = serial.NewPositioningEncoder(&target)
	suite.coder.Code(uint32(0))
	target.errorOnNextCall = true
	_, _ = suite.coder.Seek(1, io.SeekStart)

	pos, err := suite.coder.Seek(1, io.SeekStart)

	assert.Equal(suite.T(), int64(4), pos, "Should return old position value for seek")
	assert.EqualError(suite.T(), err, "errorBuffer on call number 2")
	assert.EqualError(suite.T(), suite.coder.FirstError(), "errorBuffer on call number 2", "should set first error")
}

func (suite *PositioningEncoderSuite) TestSeekReturnsErrorForSeekBeforeStart() {
	pos, err := suite.coder.Seek(-1, io.SeekStart)

	assert.Equal(suite.T(), int64(0), pos, "Should return old position value for seek")
	assert.EqualError(suite.T(), err, "seek: seeking before start")
}
