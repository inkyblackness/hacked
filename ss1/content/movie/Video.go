package movie

import (
	"context"

	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ss1/content/movie/internal/compression"
)

// Constants for video.
const (
	HighResDefaultWidth  = 600
	HighResDefaultHeight = 300
)

// Video describes the visual part of a movie.
type Video struct {
	// Width is the width of the video in pixel.
	Width uint16
	// Height is the height of the video in pixel.
	Height uint16
	// Scenes contain the frames of the video.
	Scenes []HighResScene
}

// StartPalette returns the palette of the first scene. If no scene is present, a black palette is returned.
func (video Video) StartPalette() bitmap.Palette {
	if len(video.Scenes) == 0 {
		return bitmap.Palette{}
	}
	return video.Scenes[0].palette
}

// Duration returns the sum of all scene durations.
func (video Video) Duration() Timestamp {
	var sum Timestamp
	for _, scene := range video.Scenes {
		sum = sum.Plus(scene.Duration())
	}
	return sum
}

// Encode serializes the scenes into entry buckets.
func (video Video) Encode() []EntryBucket {
	var sceneTime Timestamp
	var buckets []EntryBucket
	for index, scene := range video.Scenes {
		buckets = append(buckets, scene.Encode(sceneTime, index != 0)...)
		sceneTime = sceneTime.Plus(scene.Duration())
	}
	return buckets
}

// Decompress unpacks all the frames of all scenes.
func (video Video) Decompress() ([]Scene, error) {
	var scenes []Scene
	width := int(video.Width)
	height := int(video.Height)
	frameBuffer := make([]byte, width*height)
	decoderBuilder := compression.NewFrameDecoderBuilder(width, height)
	decoderBuilder.ForStandardFrame(frameBuffer, width)

	cloneFramebuffer := func() []byte {
		bufferCopy := make([]byte, len(frameBuffer))
		copy(bufferCopy, frameBuffer)
		return bufferCopy
	}

	for _, compressedScene := range video.Scenes {
		scenePalette := compressedScene.palette
		decoderBuilder.WithControlWords(compressedScene.controlWords)
		decoderBuilder.WithPaletteLookupList(compressedScene.paletteLookup)
		decoder := decoderBuilder.Build()
		var scene Scene
		for _, compressedFrame := range compressedScene.frames {
			err := decoder.Decode(compressedFrame.bitstream, compressedFrame.maskstream)
			if err != nil {
				return nil, err
			}

			bmp := bitmap.Bitmap{
				Header: bitmap.Header{
					Type:   bitmap.TypeFlat8Bit,
					Width:  int16(video.Width),
					Height: int16(video.Height),
					Stride: video.Width,
				},
				Palette: &scenePalette,
				Pixels:  cloneFramebuffer(),
			}
			scene.Frames = append(scene.Frames, Frame{
				Bitmap:      bmp,
				DisplayTime: compressedFrame.displayTime.ToDuration(),
			})
		}
		scenes = append(scenes, scene)
	}
	return scenes, nil
}

// HighResScene is a set of frames with high-resolution compression.
type HighResScene struct {
	palette       bitmap.Palette
	paletteLookup []byte
	controlWords  []compression.ControlWord
	frames        []HighResFrame
}

// HighResSceneFrom compresses given scene and returns the compression result.
func HighResSceneFrom(ctx context.Context, scene Scene, width, height int) (HighResScene, error) {
	encoder := compression.NewSceneEncoder(width, height)
	var palette bitmap.Palette
	for _, frame := range scene.Frames {
		if frame.Bitmap.Palette != nil {
			palette = *frame.Bitmap.Palette
		}
		err := encoder.AddFrame(frame.Bitmap.Pixels)
		if err != nil {
			return HighResScene{}, err
		}
		if ctx.Err() != nil {
			return HighResScene{}, ctx.Err()
		}
	}
	words, paletteLookup, frames, err := encoder.Encode(ctx)
	if err != nil {
		return HighResScene{}, err
	}
	compressedScene := HighResScene{
		palette:       palette,
		paletteLookup: paletteLookup,
		controlWords:  words,
		frames:        make([]HighResFrame, len(frames)),
	}
	for index, frame := range frames {
		compressedScene.frames[index] = HighResFrame{
			bitstream:   frame.Bitstream,
			maskstream:  frame.Maskstream,
			displayTime: TimestampFromDuration(scene.Frames[index].DisplayTime),
		}
	}

	return compressedScene, nil
}

// Duration returns the length of the scene, it is the sum of display times of all frames.
func (scene HighResScene) Duration() Timestamp {
	var sum Timestamp
	for _, frame := range scene.frames {
		sum = sum.Plus(frame.Duration())
	}
	return sum
}

// Encode serializes the scene into entry buckets.
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

// HighResFrame contains the compressed information of a high-resolution picture in a scene.
type HighResFrame struct {
	bitstream   []byte
	maskstream  []byte
	displayTime Timestamp
}

// Duration returns the display time of the frame.
func (frame HighResFrame) Duration() Timestamp {
	return frame.displayTime
}

// Encode serializes the frame into an entry bucket.
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
