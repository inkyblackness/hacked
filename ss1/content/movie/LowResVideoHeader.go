package movie

// LowResVideoHeaderSize is the size, in bytes, of the header structure.
const LowResVideoHeaderSize = 8

// LowResVideoHeader is for video entries with low resolution
type LowResVideoHeader struct {
	BoundingBox [4]uint16
}
