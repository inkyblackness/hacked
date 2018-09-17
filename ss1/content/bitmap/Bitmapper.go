package bitmap

import (
	"image"
	"image/color"
	"math"
)

// reference white point
var d65 = [3]float64{0.95047, 1.00000, 1.08883}

func labF(t float64) float64 {
	if t > 6.0/29.0*6.0/29.0*6.0/29.0 {
		return math.Cbrt(t)
	}
	return t/3.0*29.0/6.0*29.0/6.0 + 4.0/29.0
}

func square(value float64) float64 {
	return value * value
}

type labEntry struct {
	l float64
	a float64
	b float64
}

func labEntryFromColor(clr color.Color) labEntry {
	rLinear, gLinear, bLinear, _ := clr.RGBA()
	r, g, b := float64(rLinear)/float64(0xFFFF), float64(gLinear)/float64(0xFFFF), float64(bLinear)/float64(0xFFFF)
	x := 0.4124564*r + 0.3575761*g + 0.1804375*b
	y := 0.2126729*r + 0.7151522*g + 0.0721750*b
	z := 0.0193339*r + 0.1191920*g + 0.9503041*b
	whiteRef := d65
	fy := labF(y / whiteRef[1])
	entry := labEntry{
		l: 1.16*fy - 0.16,
		a: 5.0 * (labF(x/whiteRef[0]) - fy),
		b: 2.0 * (fy - labF(z/whiteRef[2]))}

	return entry
}

func (entry labEntry) distanceTo(other labEntry) float64 {
	return math.Sqrt(square(entry.l-other.l) + square(entry.a-other.a) + square(entry.b-other.b))
}

// Bitmapper creates bitmap images from generic images.
type Bitmapper struct {
	pal []labEntry
}

// NewBitmapper returns a new bitmapper instance based on the given palette.
func NewBitmapper(palette *Palette) *Bitmapper {
	bitmapper := &Bitmapper{}

	for _, clr := range palette {
		bitmapper.pal = append(bitmapper.pal, labEntryFromColor(clr.Color(0xFF)))
	}

	return bitmapper
}

// Map maps the provided image to a bitmap based on the internal palette.
func (bitmapper *Bitmapper) Map(img image.Image) Bitmap {
	var bmp Bitmap
	bounds := img.Bounds()

	bmp.Header.Width = int16(math.Max(0, math.Min(float64(bounds.Dx()), math.MaxInt16)))
	bmp.Header.Height = int16(math.Max(0, math.Min(float64(bounds.Dy()), math.MaxInt16)))
	bmp.Pixels = make([]byte, int(bmp.Header.Width)*int(bmp.Header.Height))
	for row := 0; row < int(bmp.Header.Height); row++ {
		for column := 0; column < int(bmp.Header.Width); column++ {
			bmp.Pixels[row*int(bmp.Header.Width)+column] = bitmapper.MapColor(img.At(column, row))
		}
	}

	return bmp
}

// MapColor maps the provided color to the nearest index in the palette.
func (bitmapper *Bitmapper) MapColor(clr color.Color) (palIndex byte) {
	_, _, _, a := clr.RGBA()
	indexWithin := func(index, from, to int) bool {
		return (index >= from) && (index <= to)
	}

	if a > 0 {
		clrEntry := labEntryFromColor(clr)
		palDistance := 1000.0

		for colorIndex, palEntry := range bitmapper.pal {
			isRegularColor := indexWithin(colorIndex, 0x01, 0x02) || indexWithin(colorIndex, 0x08, 0x0A) || indexWithin(colorIndex, 0x20, 0xFF)
			if isRegularColor {
				distance := palEntry.distanceTo(clrEntry)
				if distance < palDistance {
					palDistance = distance
					palIndex = byte(colorIndex)
				}
			}
		}
	}
	return
}
