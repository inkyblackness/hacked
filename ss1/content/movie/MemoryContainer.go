package movie

import (
	"github.com/inkyblackness/hacked/ss1/content/bitmap"
)

type memoryContainer struct {
	endTimestamp Timestamp

	videoWidth   uint16
	videoHeight  uint16
	startPalette bitmap.Palette

	audioSampleRate uint16

	entries []Entry
}

func (container *memoryContainer) EndTimestamp() Timestamp {
	return container.endTimestamp
}

func (container *memoryContainer) VideoWidth() uint16 {
	return container.videoWidth
}

func (container *memoryContainer) VideoHeight() uint16 {
	return container.videoHeight
}

func (container *memoryContainer) StartPalette() bitmap.Palette {
	return container.startPalette
}

func (container *memoryContainer) AudioSampleRate() uint16 {
	return container.audioSampleRate
}

func (container *memoryContainer) EntryCount() int {
	return len(container.entries)
}

func (container *memoryContainer) Entry(index int) Entry {
	return container.entries[index]
}
