package bitmap

// Type describes the data layout of a bitmap.
type Type byte

// Type constants are listed below.
const (
	// TypeFlat8Bit bitmaps are 8-bit paletted bitmaps that have their pixel stored in a flat layout.
	TypeFlat8Bit Type = 2
	// TypeCompressed8Bit bitmaps are 8-bit paletted bitmaps that have their pixel compressed in storage.
	// Compression is using run-length-encoding (RLE); See package rle.
	TypeCompressed8Bit Type = 4
)

// Flag adds further properties.
type Flag uint16

// Flag constants are listed below.
const (
	// FlagTransparent is set for bitmaps that shall treat palette index 0x00 as fully transparent.
	FlagTransparent Flag = 0x0001
)

// HeaderSize is the size of the Header structure, in bytes.
const HeaderSize = 28

// Area is a placeholder for either a rectangle, or an anchoring point (first two entries).
type Area [4]int16

// Header contains the meta information for a bitmap.
type Header struct {
	_             [4]byte
	Type          Type
	_             byte
	Flags         Flag
	Width         int16
	Height        int16
	Stride        uint16
	WidthFactor   byte
	HeightFactor  byte
	Area          Area
	PaletteOffset int32
}
