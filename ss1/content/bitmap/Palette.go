package bitmap

import (
	"bytes"
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
)

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

// PaletteSize indicates how many colors a palette can hold with a byte index.
const PaletteSize = 256

// Palette describes a list of colors to be used by bitmaps.
// It is an array of 256 RGB values.
type Palette [PaletteSize]RGB

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

// IndexClosestTo returns the index into this palette that matches the given color the closest.
// This search excludes the provided indices.
func (pal Palette) IndexClosestTo(rgb RGB, excluding []byte) byte {
	targetColor, _ := colorful.MakeColor(rgb.Color(0xFF))
	closestDist := float64(-1)
	closestIndex := byte(0)
	for index, other := range pal {
		if bytes.IndexByte(excluding, byte(index)) >= 0 {
			continue
		}
		otherColor, _ := colorful.MakeColor(other.Color(0xFF))
		dist := targetColor.DistanceLab(otherColor)
		if (closestDist < 0) || (dist < closestDist) {
			closestDist = dist
			closestIndex = byte(index)
		}
	}
	return closestIndex
}
