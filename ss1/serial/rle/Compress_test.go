package rle

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompressEmptyArrayResultsInTerminator(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	Compress(writer, nil)
	assert.Equal(t, []byte{0x80, 0x00, 0x00}, writer.Bytes())
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
