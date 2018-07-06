package bitmap_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/bitmap"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeReturnsErrorOnNilSource(t *testing.T) {
	_, err := bitmap.Decode(nil)

	assert.Error(t, err, "error expected")
}

func TestDecodeOfUncompressedDataReturnsBitmap(t *testing.T) {
	data := getTestData(bitmap.TypeFlat8Bit, []byte{0xAA}, false)
	bmp, err := bitmap.Decode(bytes.NewReader(data))

	require.Nil(t, err, "no error expected")
	require.NotNil(t, bmp, "bitmap expected")
	assert.Equal(t, int16(1), bmp.Header.Width)
	assert.Equal(t, int16(1), bmp.Header.Height)
	assert.Equal(t, []byte{0xAA}, bmp.Pixels)
}

func TestDecodeOfCompressedDataReturnsBitmap(t *testing.T) {
	data := getTestData(bitmap.TypeCompressed8Bit, []byte{0x00, 0x01, 0xBB, 0x80, 0x00, 0x00}, false)
	bmp, err := bitmap.Decode(bytes.NewReader(data))

	require.Nil(t, err, "no error expected")
	require.NotNil(t, bmp, "bitmap expected")
	assert.Equal(t, int16(1), bmp.Header.Width)
	assert.Equal(t, int16(1), bmp.Header.Height)
	assert.Equal(t, []byte{0xBB}, bmp.Pixels)
}

func TestDecodeWithPrivatePaletteReturnsPalette(t *testing.T) {
	data := getTestData(bitmap.TypeFlat8Bit, []byte{0xAA}, true)
	bmp, err := bitmap.Decode(bytes.NewReader(data))

	require.Nil(t, err, "no error expected")
	require.NotNil(t, bmp, "bitmap expected")
	assert.NotNil(t, bmp.Palette, "palette expected")
}

func TestWriteUncompressedWithoutPalette(t *testing.T) {
	sourceData := getTestData(bitmap.TypeFlat8Bit, []byte{0xAA}, false)
	bmp, _ := bitmap.Decode(bytes.NewReader(sourceData))

	result := bitmap.Encode(bmp, 0)
	assert.Equal(t, sourceData, result)
}

func TestWriteCompressedWithoutPalette(t *testing.T) {
	sourceData := getTestData(bitmap.TypeCompressed8Bit, []byte{0x01, 0xBB, 0x80, 0x00, 0x00}, false)
	bmp, _ := bitmap.Decode(bytes.NewReader(sourceData))

	result := bitmap.Encode(bmp, 0)
	assert.Equal(t, sourceData, result)
}

func TestWriteUncompressedWithPalette(t *testing.T) {
	sourceData := getTestData(bitmap.TypeFlat8Bit, []byte{0xCC}, true)
	bmp, _ := bitmap.Decode(bytes.NewReader(sourceData))

	result := bitmap.Encode(bmp, 0)
	assert.Equal(t, sourceData, result)
}

func TestWriteCompressedWithPalette(t *testing.T) {
	sourceData := getTestData(bitmap.TypeCompressed8Bit, []byte{0x01, 0xEE, 0x80, 0x00, 0x00}, true)
	bmp, _ := bitmap.Decode(bytes.NewReader(sourceData))

	result := bitmap.Encode(bmp, 0)
	assert.Equal(t, sourceData, result)
}

func TestWriteWithForceTransparency(t *testing.T) {
	sourceData := getTestData(bitmap.TypeFlat8Bit, []byte{0xAA}, false)
	sourceData[6] = 1
	bmp, _ := bitmap.Decode(bytes.NewReader(sourceData))

	result := bitmap.Encode(bmp, 0)
	assert.Equal(t, sourceData, result)
}

func getTestData(bmpType bitmap.Type, data []byte, withPalette bool) []byte {
	var header bitmap.Header
	buf := bytes.NewBuffer(nil)

	header.Type = bmpType
	header.Width = 1
	header.Stride = 1
	header.Height = 1
	header.WidthFactor = 1
	header.HeightFactor = 1
	if withPalette {
		header.PaletteOffset = int32(binary.Size(header) + len(data))
	}
	_ = binary.Write(buf, binary.LittleEndian, &header)
	buf.Write(data)
	if withPalette {
		var pal bitmap.Palette
		_ = binary.Write(buf, binary.LittleEndian, uint32(0x01000000))
		_ = binary.Write(buf, binary.LittleEndian, pal)
	}

	return buf.Bytes()
}
