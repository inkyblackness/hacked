package external

import (
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/content/audio/wav"
	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ui/gui"
)

// Export starts an export dialog series, calling the given callback with a folder name.
func Export(machine gui.ModalStateMachine, info string, callback func(string), lastFailed bool) {
	machine.SetState(&exportStartState{
		machine:   machine,
		callback:  callback,
		info:      info,
		withError: lastFailed,
	})
}

// ExportAudio is a helper wrapper for exporting audio.
func ExportAudio(machine gui.ModalStateMachine, filename string, sound audio.L8) {
	info := "File to be written: " + filename
	var dirHandler func(string)

	dirHandler = func(dirname string) {
		writer, err := os.Create(filepath.Join(dirname, filename))
		if err != nil {
			Export(machine, info, dirHandler, true)
			return
		}
		defer func() { _ = writer.Close() }()
		err = wav.Save(writer, sound.SampleRate, sound.Samples)
		if err != nil {
			Export(machine, info, dirHandler, true)
		}
	}

	Export(machine, info, dirHandler, false)
}

// ExportImage is a helper wrapper for exporting single images.
func ExportImage(machine gui.ModalStateMachine, filename string, bmp bitmap.Bitmap) {
	info := "File to be written: " + filename
	var exportTo func(string)

	exportTo = func(dirname string) {
		writer, err := os.Create(filepath.Join(dirname, filename))
		if err != nil {
			Export(machine, "Could not create file.\n"+info, exportTo, true)
			return
		}
		defer func() { _ = writer.Close() }()

		imageRect := image.Rect(0, 0, int(bmp.Header.Width), int(bmp.Header.Height))
		imagePal := bmp.Palette.ColorPalette(false)
		paletted := image.NewPaletted(imageRect, imagePal)
		paletted.Pix = bmp.Pixels
		err = png.Encode(writer, paletted)
		if err != nil {
			Export(machine, info, exportTo, true)
			return
		}
	}

	Export(machine, info, exportTo, false)
}
