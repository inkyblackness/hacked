package external

import (
	"os"

	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/content/audio/wav"
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
