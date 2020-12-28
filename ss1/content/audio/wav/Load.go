package wav

import (
	"io"

	"github.com/inkyblackness/hacked/ss1"
	"github.com/inkyblackness/hacked/ss1/content/audio"
)

const (
	errSourceIsNil       ss1.StringError = "source is nil"
	errNotASupportedWave ss1.StringError = "not a supported WAV"
)

// Load reads from the provided source and returns the data.
func Load(source io.Reader) (data audio.L8, err error) {
	if source == nil {
		return data, errSourceIsNil
	}

	var loader waveLoader

	loader.load(source)
	if loader.err != nil {
		return data, loader.err
	}

	data.SampleRate = loader.sampleRate
	data.Samples = loader.samples

	return
}
