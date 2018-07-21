package movie

// HighResVideoHeaderSize is the size, in bytes, of the header structure.
const HighResVideoHeaderSize = 2

// HighResVideoHeader is for video entries with high resolution.
type HighResVideoHeader struct {
	PixelDataOffset uint16
}
