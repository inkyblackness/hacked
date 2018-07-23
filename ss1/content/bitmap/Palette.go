package bitmap

import "image/color"

// RGB is describing Red, Green, and Blue intensities with 8-bit resolution.
type RGB struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

// Color returns the RGB data as a regular Color entry.
func (col RGB) Color(alpha byte) color.Color {
	return color.RGBA{
		R: col.Red,
		G: col.Green,
		B: col.Blue,
		A: alpha,
	}
}

// Palette describes a list of colors to be used by bitmaps.
// It is an array of 256 RGB values.
type Palette [256]RGB

// ColorPalette returns a palette usable for the image packages.
func (pal Palette) ColorPalette(firstIndexTransparent bool) color.Palette {
	result := make(color.Palette, len(pal))
	for index, col := range pal {
		alpha := byte(0xFF)
		if (index == 0) && firstIndexTransparent {
			alpha = 0x00
		}
		result[index] = col.Color(alpha)
	}
	return result
}
