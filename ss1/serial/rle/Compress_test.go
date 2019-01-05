package rle_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/inkyblackness/hacked/ss1/serial/rle"
)

func TestCompressEmptyArrayResultsInTerminator(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	err := rle.Compress(writer, nil, nil)
	assert.Nil(t, err, "no error expected")
	assert.Equal(t, []byte{0x80, 0x00, 0x00}, writer.Bytes())
}

func TestCompressOfZeroBytesLong(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	input := make([]byte, 1000)
	input[len(input)-1] = 1
	err := rle.Compress(writer, input, nil)
	assert.Nil(t, err, "no error expected")
	assert.Equal(t, []byte{0x80, 0xE7, 0x03, 0x01, 0x01, 0x80, 0x0, 0x0}, writer.Bytes())
}

func TestCompressOfEqualBytesLong(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	input := make([]byte, 1000)
	for i := 0; i < len(input); i++ {
		input[i] = 0xAA
	}
	input[len(input)-1] = 1
	err := rle.Compress(writer, input, nil)
	assert.Nil(t, err, "no error expected")
	assert.Equal(t, []byte{0x80, 0xE7, 0xC3, 0xAA, 0x01, 0x01, 0x80, 0x0, 0x0}, writer.Bytes())
}

func TestCompressOfEqualBytesShort(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	input := make([]byte, 40)
	for i := 0; i < len(input); i++ {
		input[i] = 0xAA
	}
	input[len(input)-1] = 1
	err := rle.Compress(writer, input, nil)
	assert.Nil(t, err, "no error expected")
	assert.Equal(t, []byte{0x00, 0x27, 0xAA, 0x01, 0x01, 0x80, 0x0, 0x0}, writer.Bytes())
}

func TestCompressOfRandomBytesShort(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	input := make([]byte, 10)
	for i := 0; i < len(input); i++ {
		input[i] = byte(i)
	}
	err := rle.Compress(writer, input, nil)
	assert.Nil(t, err, "no error expected")
	assert.Equal(t, []byte{0x0A, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x80, 0x00, 0x00}, writer.Bytes())
}

func TestCompressWithIdenticalReference(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	input := []byte{0xAA, 0xBB, 0xCC, 0xDD}
	reference := input
	err := rle.Compress(writer, input, reference)
	assert.Nil(t, err, "no error expected")
	assert.Equal(t, []byte{0x80, 0x00, 0x00}, writer.Bytes())
}

func TestCompressWithPartlyIdenticalReference(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	input := []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}
	reference := []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0x0A}
	err := rle.Compress(writer, input, reference)
	assert.Nil(t, err, "no error expected")
	assert.Equal(t, []byte{0x85, 0x01, 0xFF, 0x80, 0x00, 0x00}, writer.Bytes())
}
