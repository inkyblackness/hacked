package voc

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/inkyblackness/hacked/ss1"
	"github.com/inkyblackness/hacked/ss1/content/audio"
)

const (
	errSourceIsNil  ss1.StringError = "source is nil"
	errNoAudioFound ss1.StringError = "no audio found"
)

// Load reads from the provided source a Creative Voice Sound and returns the data.
func Load(source io.Reader) (data audio.L8, err error) {
	if source == nil {
		return data, errSourceIsNil
	}

	err = readAndVerifyHeader(source)
	if err != nil {
		return
	}
	return readSoundData(source)
}

func readAndVerifyHeader(source io.Reader) error {
	start := make([]byte, len(FileHeader))
	headerSize := uint16(0)
	version := uint16(0)
	versionValidity := uint16(0)

	_, err := source.Read(start)
	if err != nil {
		return err
	}
	err = binary.Read(source, binary.LittleEndian, &headerSize)
	if err != nil {
		return err
	}
	err = binary.Read(source, binary.LittleEndian, &version)
	if err != nil {
		return err
	}
	err = binary.Read(source, binary.LittleEndian, &versionValidity)
	if err != nil {
		return err
	}

	calculated := ^version + versionCheckValue
	if calculated != versionValidity {
		return fmt.Errorf("version validity failed: 0x%04X != 0x%04X", calculated, versionValidity)
	}

	skip := make([]byte, headerSize-standardHeaderSize)
	_, err = source.Read(skip)
	return err
}

func readSoundData(source io.Reader) (data audio.L8, err error) {
	sampleRate := float32(0.0)
	var samples []byte
	done := false

	for !done {
		blockStart := make([]byte, 4)

		_, err = source.Read(blockStart)
		if err != nil {
			return
		}
		switch blockType(blockStart[0]) {
		case soundData:
			{
				meta := make([]byte, 2)
				_, err = source.Read(meta)
				if err != nil {
					return
				}
				sampleRate = divisorToSampleRate(meta[0])

				newCount := lengthFromBlockStart(blockStart) - len(meta)
				buf := make([]byte, newCount)
				_, err = source.Read(buf)
				if err != nil {
					return
				}

				oldCount := len(samples)
				newSamples := make([]byte, oldCount+newCount)
				copy(newSamples, samples)
				copy(newSamples[oldCount:], buf)
				samples = newSamples
			}
		case terminator:
			{
				done = true
			}
		}
	}

	if len(samples) == 0 {
		return data, errNoAudioFound
	}

	return audio.L8{
		SampleRate: sampleRate,
		Samples:    samples,
	}, nil
}
