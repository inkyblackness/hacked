package bitmap

import (
	"encoding/binary"
	"errors"
	"io"

	"bytes"
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
	if reader == nil {
		return nil, errors.New("reader is nil")
	}

	var bmp Bitmap

	err := binary.Read(reader, binary.LittleEndian, &bmp.Header)
	if err != nil {
		return nil, err
	}
	bmp.Pixels = make([]byte, int(bmp.Header.Height)*int(bmp.Header.Stride))
	if bmp.Header.Type == TypeCompressed8Bit {
		err = rle.Decompress(reader, bmp.Pixels)
	} else {
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
		rle.Compress(buf, rawData, nil)
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
		_ = binary.Write(buf, binary.LittleEndian, uint32(1))
		_ = binary.Write(buf, binary.LittleEndian, bmp.Palette)
	}

	return buf.Bytes()
}
