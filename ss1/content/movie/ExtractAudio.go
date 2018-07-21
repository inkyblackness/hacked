package movie

import (
	"bytes"

	"github.com/inkyblackness/hacked/ss1/content/audio"
)

// ExtractAudio decodes the given data array as a MOVI container and
// extracts the audio track.
func ExtractAudio(data []byte) (sound audio.L8, err error) {
	container, err := Read(bytes.NewReader(data))

	if container != nil {
		var samples []byte

		for i := 0; i < container.EntryCount(); i++ {
			entry := container.Entry(i)

			if entry.Type() == Audio {
				samples = append(samples, entry.Data()...)
			}
		}
		sound.SampleRate = float32(container.AudioSampleRate())
		sound.Samples = samples
	}
	return
}
