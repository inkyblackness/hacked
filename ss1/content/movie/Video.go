package movie

import (
	"context"

	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ss1/content/movie/internal/compression"
	"github.com/inkyblackness/hacked/ss1/content/movie/internal/format"
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

func (video Video) duration() format.Timestamp {
	var sum format.Timestamp
	for _, scene := range video.Scenes {
		sum = sum.Plus(scene.duration())
	}
	return sum
}

func (video Video) encode() []format.EntryBucket {
	var sceneTime format.Timestamp
	var buckets []format.EntryBucket
	for index, scene := range video.Scenes {
		buckets = append(buckets, scene.encode(sceneTime, index != 0)...)
		sceneTime = sceneTime.Plus(scene.duration())
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
		decoderBuilder.WithControlWords(compressedScene.controlWords)
		decoderBuilder.WithPaletteLookupList(compressedScene.paletteLookup)
		decoder := decoderBuilder.Build()
		var scene Scene
		scene.Palette = compressedScene.palette
		for _, compressedFrame := range compressedScene.frames {
			err := decoder.Decode(compressedFrame.bitstream, compressedFrame.maskstream)
			if err != nil {
				return nil, err
			}
			scene.Frames = append(scene.Frames, Frame{
				Pixels:      cloneFramebuffer(),
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
	for _, frame := range scene.Frames {
		err := encoder.AddFrame(frame.Pixels)
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
		palette:       scene.Palette,
		paletteLookup: paletteLookup,
		controlWords:  words,
		frames:        make([]HighResFrame, len(frames)),
	}
	for index, frame := range frames {
		compressedScene.frames[index] = HighResFrame{
			bitstream:   frame.Bitstream,
			maskstream:  frame.Maskstream,
			displayTime: format.TimestampFromDuration(scene.Frames[index].DisplayTime),
		}
	}

	return compressedScene, nil
}

func (scene HighResScene) duration() format.Timestamp {
	var sum format.Timestamp
	for _, frame := range scene.frames {
		sum = sum.Plus(frame.duration())
	}
	return sum
}

func (scene HighResScene) encode(start format.Timestamp, withPalette bool) []format.EntryBucket {
	buckets := make([]format.EntryBucket, 0, len(scene.frames)+1)
	controlEntries := []format.Entry{
		{Timestamp: format.Timestamp{}, Data: format.PaletteLookupEntryData{List: scene.paletteLookup}},
		{Timestamp: format.Timestamp{}, Data: format.ControlDictionaryEntryData{Words: scene.controlWords}},
	}
	if withPalette {
		controlEntries = append(controlEntries,
			format.Entry{Timestamp: start, Data: format.PaletteResetEntryData{}},
			format.Entry{Timestamp: start, Data: format.PaletteEntryData{Colors: scene.palette}})
	}
	buckets = append(buckets,
		format.EntryBucket{
			Priority:  format.EntryBucketPriorityVideoControl,
			Timestamp: start,
			Entries:   controlEntries,
		})
	frameTime := start
	for _, frame := range scene.frames {
		buckets = append(buckets, frame.encode(frameTime))
		frameTime = frameTime.Plus(frame.displayTime)
	}

	return buckets
}

// HighResFrame contains the compressed information of a high-resolution picture in a scene.
type HighResFrame struct {
	bitstream   []byte
	maskstream  []byte
	displayTime format.Timestamp
}

func (frame HighResFrame) duration() format.Timestamp {
	return frame.displayTime
}

func (frame HighResFrame) encode(start format.Timestamp) format.EntryBucket {
	return format.EntryBucket{
		Priority:  format.EntryBucketPriorityFrame,
		Timestamp: start,
		Entries: []format.Entry{
			{
				Timestamp: start,
				Data: format.HighResVideoEntryData{
					Bitstream:  frame.bitstream,
					Maskstream: frame.maskstream,
				},
			},
		},
	}
}
