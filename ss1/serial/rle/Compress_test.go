package rle

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompressEmptyArrayResultsInTerminator(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	Compress(writer, nil, nil)
	assert.Equal(t, []byte{0x80, 0x00, 0x00}, writer.Bytes())
}

func TestCompressOfZeroBytesLong(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	input := make([]byte, 1000)
	input[len(input)-1] = 1
	Compress(writer, input, nil)
	assert.Equal(t, []byte{0x80, 0xE7, 0x03, 0x01, 0x01, 0x80, 0x0, 0x0}, writer.Bytes())
}

func TestCompressOfEqualBytesLong(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	input := make([]byte, 1000)
	for i := 0; i < len(input); i++ {
		input[i] = 0xAA
	}
	input[len(input)-1] = 1
	Compress(writer, input, nil)
	assert.Equal(t, []byte{0x80, 0xE7, 0xC3, 0xAA, 0x01, 0x01, 0x80, 0x0, 0x0}, writer.Bytes())
}

func TestCompressOfEqualBytesShort(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	input := make([]byte, 40)
	for i := 0; i < len(input); i++ {
		input[i] = 0xAA
	}
	input[len(input)-1] = 1
	Compress(writer, input, nil)
	assert.Equal(t, []byte{0x00, 0x27, 0xAA, 0x01, 0x01, 0x80, 0x0, 0x0}, writer.Bytes())
}

func TestCompressOfRandomBytesShort(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	input := make([]byte, 10)
	for i := 0; i < len(input); i++ {
		input[i] = byte(i)
	}
	Compress(writer, input, nil)
	assert.Equal(t, []byte{0x81, 0x09, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x80, 0x00, 0x00}, writer.Bytes())
}

func TestCompressWriteZeroOfLessThan80(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	writeZero(writer, 0x7F)
	assert.Equal(t, []byte{0xFF}, writer.Bytes())
}

func TestCompressWriteZeroOfLessThanFF(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	writeZero(writer, 0xFD)
	assert.Equal(t, []byte{0xFF, 0xFE}, writer.Bytes())
}

func TestCompressWriteZeroOfLessThan8000(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	writeZero(writer, 0x7FFC)
	assert.Equal(t, []byte{0x80, 0xFC, 0x7F}, writer.Bytes())
}

func TestCompressWriteZeroOf8000(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	writeZero(writer, 0x8000)
	assert.Equal(t, []byte{0x80, 0xFF, 0x7F, 0x81}, writer.Bytes())
}

func TestCompressWriteRawOfLessThan80(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	writeRaw(writer, []byte{0x0A, 0x0B, 0x0C})
	assert.Equal(t, []byte{0x03, 0x0A, 0x0B, 0x0C}, writer.Bytes())
}
