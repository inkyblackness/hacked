package bitmap

// RGB is describing Red, Green, and Blue intensities with 8-bit resolution.
type RGB struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

// Palette describes a list of colors to be used by bitmaps.
// It is an array of 256 RGB values.
type Palette [256]RGB
