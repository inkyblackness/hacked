package external

import (
	"image"
	"image/color"
	"math"
	"os"

	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/content/audio/wav"
	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ui/gui"
)

// Import starts an import dialog series, calling the given callback with a file name.
func Import(machine gui.ModalStateMachine, info string, callback func(string), lastFailed bool) {
	machine.SetState(&importStartState{
		machine:   machine,
		callback:  callback,
		info:      info,
		withError: lastFailed,
	})
}

// ImportAudio is a helper to handle audio file import. The callback is called with the loaded audio.
func ImportAudio(machine gui.ModalStateMachine, callback func(l8 audio.L8)) {
	info := "File must be a WAV file, 22050 Hz, 8-bit or 16-bit, uncompressed."
	var fileHandler func(string)

	fileHandler = func(filename string) {
		reader, err := os.Open(filename)
		if err != nil {
			Import(machine, info, fileHandler, true)
			return
		}
		defer func() { _ = reader.Close() }()
		sound, err := wav.Load(reader)
		if err != nil {
			Import(machine, info, fileHandler, true)
			return
		}
		callback(sound)
	}

	Import(machine, info, fileHandler, false)
}

// ImportImage is a helper to handle image file import. The callback is called with the loaded image.
func ImportImage(machine gui.ModalStateMachine, paletteRetriever func() (bitmap.Palette, error), callback func(bitmap.Bitmap)) {
	info := "File should be either a PNG or a GIF file.\nPaletted images matching game palette are taken 1:1,\nothers are mapped closest fitting."
	var fileHandler func(string)

	fileHandler = func(filename string) {
		reader, err := os.Open(filename)
		if err != nil {
			Import(machine, "Could not open file.\n"+info, fileHandler, true)
			return
		}
		defer func() { _ = reader.Close() }()
		img, _, err := image.Decode(reader)
		if err != nil {
			Import(machine, "File not recognized as image.\n"+info, fileHandler, true)
			return
		}

		var bmp bitmap.Bitmap
		importMapped := true
		rawPalette, err := paletteRetriever()
		if err != nil {
			Import(machine, "Can not import image without having a palette loaded.\n"+info, fileHandler, true)
			return
		}
		if palettedImg, isPaletted := img.(image.PalettedImage); isPaletted {
			imgPalette, hasPalette := palettedImg.ColorModel().(color.Palette)
			if hasPalette && paletteMatches(imgPalette, rawPalette.ColorPalette(false)) {
				bounds := img.Bounds()

				bmp.Header.Width = int16(math.Max(0, math.Min(float64(bounds.Dx()), math.MaxInt16)))
				bmp.Header.Height = int16(math.Max(0, math.Min(float64(bounds.Dy()), math.MaxInt16)))
				bmp.Pixels = make([]byte, int(bmp.Header.Width)*int(bmp.Header.Height))
				for row := 0; row < int(bmp.Header.Height); row++ {
					for column := 0; column < int(bmp.Header.Width); column++ {
						bmp.Pixels[row*int(bmp.Header.Width)+column] = palettedImg.ColorIndexAt(column, row)
					}
				}
				importMapped = false
			}
		}
		if importMapped {
			bitmapper := bitmap.NewBitmapper(&rawPalette)
			bmp = bitmapper.Map(img)
		}
		callback(bmp)
	}

	Import(machine, info, fileHandler, false)
}

func paletteMatches(imgPalette color.Palette, rawPalette color.Palette) bool {
	if len(imgPalette) > len(rawPalette) {
		return false
	}

	for index, clr := range imgPalette {
		imgR, imgG, imgB, _ := clr.RGBA()
		rawR, rawG, rawB, _ := rawPalette[index].RGBA()

		if (imgR != rawR) || (imgG != rawG) || (imgB != rawB) {
			return false
		}
	}

	return true
}
