package rle_test

import (
	"bytes"
	"testing"

	"github.com/inkyblackness/hacked/ss1/serial/rle"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecompressEmptyArrayReturnsError(t *testing.T) {
	err := rle.Decompress(bytes.NewReader(nil), make([]byte, 100))

	assert.NotNil(t, err)
}

func TestDecompress800000IsEndOfStream(t *testing.T) {
	var result []byte
	err := rle.Decompress(bytes.NewReader([]byte{0x80, 0x00, 0x00}), result)

	assert.Nil(t, err)
}

func TestDecompress800000AppendsRemainingZeroes(t *testing.T) {
	result := make([]byte, 10)
	err := rle.Decompress(bytes.NewReader([]byte{0x80, 0x00, 0x00}), result)

	require.Nil(t, err)
	assert.Equal(t, make([]byte, 10), result)
}

func TestDecompress800000IsConsumedAtEnd(t *testing.T) {
	var result []byte
	reader := bytes.NewReader([]byte{0x80, 0x00, 0x00})
	err := rle.Decompress(reader, result)

	require.Nil(t, err)
	pos, _ := reader.Seek(0, 1)
	assert.Equal(t, int64(3), pos)
}

func TestDecompress00WritesNNBytesOfColorZZ(t *testing.T) {
	result := make([]byte, 5)
	err := rle.Decompress(bytes.NewReader([]byte{0x00, 0x05, 0xCC, 0x80, 0x00, 0x00}), result)

	require.Nil(t, err)
	assert.Equal(t, []byte{0xCC, 0xCC, 0xCC, 0xCC, 0xCC}, result)
}

func TestDecompress00ReturnsErrorIfZZIsMissing(t *testing.T) {
	err := rle.Decompress(bytes.NewReader([]byte{0x00, 0x05}), make([]byte, 5))

	assert.NotNil(t, err)
}

func TestDecompress00ReturnsErrorIfNNIsMissing(t *testing.T) {
	err := rle.Decompress(bytes.NewReader([]byte{0x00}), make([]byte, 5))

	assert.NotNil(t, err)
}

func TestDecompressNNLess80WritesNNFollowingBytes(t *testing.T) {
	result := make([]byte, 2)
	err := rle.Decompress(bytes.NewReader([]byte{0x02, 0xAA, 0xBB, 0x80, 0x00, 0x00}), result)

	require.Nil(t, err)
	assert.Equal(t, []byte{0xAA, 0xBB}, result)
}

func TestDecompressNNLess80ReturnsErrorIfEndOfFile(t *testing.T) {
	err := rle.Decompress(bytes.NewReader([]byte{0x02, 0xAA}), make([]byte, 2))

	assert.NotNil(t, err)
}

func TestDecompress80NNNNSkipsBytes(t *testing.T) {
	expected := make([]byte, 0x123)
	result := make([]byte, 0x123)
	for i := 0; i < len(result); i++ {
		result[i] = byte(i)
	}
	copy(expected, result)
	err := rle.Decompress(bytes.NewReader([]byte{0x80, 0x23, 0x01, 0x80, 0x00, 0x00}), result)

	require.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestDecompress80CopiesNNBytes(t *testing.T) {
	input := bytes.NewBuffer(nil)
	input.Write([]byte{0x80, 0x04, 0x80})
	input.Write([]byte{0x01, 0x02, 0x03, 0x04})
	input.Write([]byte{0x80, 0x00, 0x00})
	result := make([]byte, 4)
	err := rle.Decompress(bytes.NewReader(input.Bytes()), result)

	require.Nil(t, err)
	assert.Equal(t, []byte{0x01, 0x02, 0x03, 0x04}, result)
}

func TestDecompress80CopiesNNBytesExtended(t *testing.T) {
	input := bytes.NewBuffer(nil)
	expected := make([]byte, 0x3FFF)
	input.Write([]byte{0x80, 0xFF, 0xBF})
	input.Write(expected)
	input.Write([]byte{0x80, 0x00, 0x00})
	result := make([]byte, len(expected))
	err := rle.Decompress(bytes.NewReader(input.Bytes()), result)

	require.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestDecompress80ReturnsErrorForUndefinedCase(t *testing.T) {
	input := bytes.NewBuffer(nil)
	input.Write([]byte{0x80, 0x00, 0xC0})
	input.Write([]byte{0x80, 0x00, 0x00})
	err := rle.Decompress(bytes.NewReader(input.Bytes()), make([]byte, 1))

	assert.NotNil(t, err)
}

func TestDecompress80WritesNNBytesOfValue(t *testing.T) {
	input := bytes.NewBuffer(nil)
	expected := make([]byte, 0x3FFF)
	for index := range expected {
		expected[index] = 0xCD
	}
	input.Write([]byte{0x80, 0xFF, 0xFF, 0xCD})
	input.Write([]byte{0x80, 0x00, 0x00})
	result := make([]byte, len(expected))
	err := rle.Decompress(bytes.NewReader(input.Bytes()), result)

	require.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestDecompressNNMore80WritesZeroes(t *testing.T) {
	result := make([]byte, 3)
	err := rle.Decompress(bytes.NewReader([]byte{0x83, 0x80, 0x00, 0x00}), result)

	require.Nil(t, err)
	assert.Equal(t, []byte{0, 0, 0}, result)
}
