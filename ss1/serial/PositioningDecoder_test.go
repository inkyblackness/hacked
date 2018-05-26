package serial_test

import (
	"bytes"
	"testing"

	"github.com/inkyblackness/hacked/ss1/serial"

	"github.com/stretchr/testify/assert"
)

func TestPositioningDecoderSetCurPosRepositionsReadOffset(t *testing.T) {
	var source = bytes.NewReader([]byte{0x78, 0x12, 0x34})
	coder := serial.NewPositioningDecoder(source)
	arrayValue := make([]byte, 3)
	var intValue uint16

	coder.Code(arrayValue)
	coder.SetCurPos(1)
	coder.Code(&intValue)

	assert.Equal(t, uint16(0x3412), intValue)
}

func TestPositioningDecoderCurPosReturnsCurrentOffset(t *testing.T) {
	var source = bytes.NewReader([]byte{0x78, 0x12, 0x34})
	coder := serial.NewPositioningDecoder(source)

	coder.Code(uint16(0))

	assert.Equal(t, uint32(2), coder.CurPos())
}

func TestPositioningDecoderFirstErrorBySetCurPos(t *testing.T) {
	var target errorBuffer
	coder := serial.NewPositioningDecoder(&target)

	coder.Code(uint32(0))
	target.errorOnNextCall = true
	coder.SetCurPos(0)

	assert.EqualError(t, coder.FirstError(), "errorBuffer on call number 2")
}

func TestPositioningDecoderSetCurPosDoesNothingOnPreviousError(t *testing.T) {
	var target errorBuffer
	coder := serial.NewPositioningDecoder(&target)
	target.errorOnNextCall = true
	coder.Code(uint32(0))
	target.errorOnNextCall = true
	coder.SetCurPos(0)

	assert.Equal(t, target.callCounter, 1)
	assert.EqualError(t, coder.FirstError(), "errorBuffer on call number 1")
}

func TestCurPosIsNotChangedOnSetError(t *testing.T) {
	var target errorBuffer
	coder := serial.NewPositioningDecoder(&target)
	target.errorOnNextCall = true
	coder.SetCurPos(4)
	assert.Equal(t, uint32(0), coder.CurPos())
}
