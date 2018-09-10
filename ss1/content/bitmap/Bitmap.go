package bitmap

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/inkyblackness/hacked/ss1/serial/rle"
)

// Bitmap describes a palette based image.
type Bitmap struct {
	Header  Header
	Pixels  []byte
	Palette *Palette
}

// Decode tries to read a bitmap from given reader.
// Should the bitmap be a compressed bitmap, then a reference image with pixel data all 0x00 is assumed.
func Decode(reader io.Reader) (*Bitmap, error) {
	return DecodeReferenced(reader, func(width, height int16) ([]byte, error) {
		return make([]byte, int(width)*int(height)), nil
	})
}

// DecodeReferenced tries to read a bitmap from given reader.
// If the serialized bitmap describes a compressed bitmap, then the pixels from the reference are used as a basis for the result.
// The returned byte array from the provider will be used as pixel buffer for the new bitmap.
func DecodeReferenced(reader io.Reader, provider func(width, height int16) ([]byte, error)) (*Bitmap, error) {
	if reader == nil {
		return nil, errors.New("reader is nil")
	}

	var bmp Bitmap

	err := binary.Read(reader, binary.LittleEndian, &bmp.Header)
	if err != nil {
		return nil, err
	}
	if bmp.Header.Type == TypeCompressed8Bit {
		if bmp.Header.Stride != uint16(bmp.Header.Width) {
			return nil, errors.New("stride not equal to width for compressed bitmap")
		}
		bmp.Pixels, err = provider(bmp.Header.Width, bmp.Header.Height)
		if err != nil {
			return nil, err
		}
		err = rle.Decompress(reader, bmp.Pixels)
	} else {
		bmp.Pixels = make([]byte, int(bmp.Header.Height)*int(bmp.Header.Stride))
		_, err = reader.Read(bmp.Pixels)
	}
	if err != nil {
		return nil, errors.New("data could not be read")
	}

	if bmp.Header.PaletteOffset != 0 {
		paletteFlag := uint32(0)
		_ = binary.Read(reader, binary.LittleEndian, &paletteFlag)
		bmp.Palette = new(Palette)
		err = binary.Read(reader, binary.LittleEndian, bmp.Palette)
		if err != nil {
			return nil, errors.New("palette could not be read")
		}
	}

	return &bmp, nil
}

// Encode writes the bitmap to a byte array and returns it.
// Compressed bitmaps will be compressed with a reference image with pixel data all 0x00.
func Encode(bmp *Bitmap, offsetBase int) []byte {
	rawData := bmp.Pixels
	if bmp.Header.Type == TypeCompressed8Bit {
		buf := bytes.NewBuffer(nil)
		_ = rle.Compress(buf, rawData, nil)
		rawData = buf.Bytes()
	}
	header := bmp.Header
	if bmp.Palette != nil {
		header.PaletteOffset = int32(offsetBase + HeaderSize + len(rawData))
	}

	buf := bytes.NewBuffer(nil)
	_ = binary.Write(buf, binary.LittleEndian, &header)
	_ = binary.Write(buf, binary.LittleEndian, rawData)
	if bmp.Palette != nil {
		_ = binary.Write(buf, binary.LittleEndian, paletteMarker)
		_ = binary.Write(buf, binary.LittleEndian, bmp.Palette)
	}

	return buf.Bytes()
}
