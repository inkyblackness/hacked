package wav

import (
	"fmt"
	"io"

	"github.com/inkyblackness/hacked/ss1/content/audio"
)

var errNotASupportedWave = fmt.Errorf("not a supported WAV")

// Load reads from the provided source and returns the data.
func Load(source io.Reader) (data audio.L8, err error) {
	if source == nil {
		return data, fmt.Errorf("source is nil")
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
