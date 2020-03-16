package movie

import (
	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ss1/content/movie/internal/compression"
)

type Video struct {
	// Width is the width of the video in pixel.
	Width uint16
	// Height is the height of the video in pixel.
	Height uint16

	Scenes []HighResScene
}

func (video Video) StartPalette() bitmap.Palette {
	if len(video.Scenes) == 0 {
		return bitmap.Palette{}
	}
	return video.Scenes[0].palette
}

func (video Video) Duration() Timestamp {
	var sum Timestamp
	for _, scene := range video.Scenes {
		sum = sum.Plus(scene.Duration())
	}
	return sum
}

func (video Video) Encode() []EntryBucket {
	var sceneTime Timestamp
	var buckets []EntryBucket
	for index, scene := range video.Scenes {
		buckets = append(buckets, scene.Encode(sceneTime, index != 0)...)
		sceneTime = sceneTime.Plus(scene.Duration())
	}
	return buckets
}

type HighResScene struct {
	palette       bitmap.Palette
	controlWords  []compression.ControlWord
	paletteLookup []byte
	frames        []HighResFrame
}

func (scene HighResScene) Duration() Timestamp {
	var sum Timestamp
	for _, frame := range scene.frames {
		sum = sum.Plus(frame.Duration())
	}
	return sum
}

func (scene HighResScene) Encode(start Timestamp, withPalette bool) []EntryBucket {
	buckets := make([]EntryBucket, 0, len(scene.frames)+1)
	controlEntries := []Entry{
		{Timestamp: Timestamp{}, Data: PaletteLookupEntryData{List: scene.paletteLookup}},
		{Timestamp: Timestamp{}, Data: ControlDictionaryEntryData{Words: scene.controlWords}},
	}
	if withPalette {
		controlEntries = append(controlEntries,
			Entry{Timestamp: start, Data: PaletteResetEntryData{}},
			Entry{Timestamp: start, Data: PaletteEntryData{Colors: scene.palette}})
	}
	buckets = append(buckets,
		EntryBucket{
			Priority:  EntryBucketPriorityVideoControl,
			Timestamp: start,
			Entries:   controlEntries,
		})
	frameTime := start
	for _, frame := range scene.frames {
		buckets = append(buckets, frame.Encode(frameTime))
		frameTime = frameTime.Plus(frame.displayTime)
	}

	return buckets
}

type HighResFrame struct {
	bitstream   []byte
	maskstream  []byte
	displayTime Timestamp
}

func (frame HighResFrame) Duration() Timestamp {
	return frame.displayTime
}

func (frame HighResFrame) Encode(start Timestamp) EntryBucket {
	return EntryBucket{
		Priority:  EntryBucketPriorityFrame,
		Timestamp: start,
		Entries: []Entry{
			{
				Timestamp: start,
				Data: HighResVideoEntryData{
					Bitstream:  frame.bitstream,
					Maskstream: frame.maskstream,
				},
			},
		},
	}
}
