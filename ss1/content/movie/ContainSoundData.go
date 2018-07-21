package movie

import (
	"bytes"

	"github.com/inkyblackness/hacked/ss1/content/audio"
)

const audioEntrySize = 0x2000

// ContainSoundData packs a sound data into a container and encodes it.
func ContainSoundData(soundData audio.L8) []byte {
	builder := NewContainerBuilder()
	startOffset := 0
	entryStartTime := float32(0)
	timePerEntry := timeFromRaw(timeToRaw(float32(audioEntrySize) / soundData.SampleRate))

	for (startOffset + audioEntrySize) <= len(soundData.Samples) {
		endOffset := startOffset + audioEntrySize
		builder.AddEntry(NewMemoryEntry(entryStartTime, Audio, soundData.Samples[startOffset:endOffset]))
		entryStartTime += timePerEntry
		startOffset = endOffset
	}
	if startOffset < len(soundData.Samples) {
		builder.AddEntry(NewMemoryEntry(entryStartTime, Audio, soundData.Samples[startOffset:]))
		entryStartTime += timeFromRaw(timeToRaw(float32(len(soundData.Samples)-startOffset) / soundData.SampleRate))
	}

	builder.MediaDuration(entryStartTime)
	builder.AudioSampleRate(uint16(soundData.SampleRate))

	container := builder.Build()
	buffer := bytes.NewBuffer(nil)
	_ = Write(buffer, container)
	return buffer.Bytes()
}
