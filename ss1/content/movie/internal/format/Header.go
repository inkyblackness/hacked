package format

import (
	"fmt"
	"math"
)

// Tag is the identifier key of a MOVI container.
const Tag string = "MOVI"

// HeaderSize specifies the length of a serialized Header in bytes.
const HeaderSize = 256

// Fix represents a fixed-point number, with 16 bits of fraction and 15 bits of integer number.
type Fix struct {
	Fraction uint16
	Number   int16
}

// FixFromFloat returns a fixed-point number from given floating-point number.
func FixFromFloat(value float32) Fix {
	number := int16(value)
	return Fix{
		Fraction: uint16(math.Abs(float64(value)-float64(number)) * 0x10000),
		Number:   number,
	}
}

// ToFloat returns the equivalent floating-point number.
func (f Fix) ToFloat() float32 {
	return float32(f.Number) + float32(f.Fraction)/0x10000
}

// Header is the first entry of a MOVI container.
type Header struct {
	// Tag is the identifier key of a MOVI container
	Tag [4]byte
	// IndexEntryCount specifies how many index entries are in the index.
	IndexEntryCount int32
	// IndexSize specifies how long the index is.
	IndexSize int32
	// ContentSize specifies the length of the content.
	ContentSize int32
	// Duration is the length of the media in seconds.
	Duration Fix

	// VideoFrameRate gives a hint on number of frames per second.
	VideoFrameRate Fix
	// VideoWidth specifies the width of the video.
	VideoWidth uint16
	// VideoHeight specifies the height of the video.
	VideoHeight uint16
	// VideoBitsPerPixel specifies bit count for color information: 8, 15, 24.
	VideoBitsPerPixel int16
	// VideoPalettePresent specifies whether a palette is set.
	VideoPalettePresent int16

	// AudioChannelCount specifies whether there is audio, and how many channels to use. 0 = no audio, 1 = mono, 2 = stereo.
	AudioChannelCount int16
	// AudioBytesPerSample specifies how many bytes one audio sample has. 1 = 8-bit, 2 = 16-bit.
	AudioBytesPerSample int16
	// AudioSampleRate specifies the sample rate of the audio.
	AudioSampleRate Fix

	_ [216]byte
}

func (header *Header) String() (result string) {
	result += fmt.Sprintf("Index: %d entries (%d bytes)\n", header.IndexEntryCount, header.IndexSize)
	result += fmt.Sprintf("Content: %d bytes, %f seconds\n", header.ContentSize, header.Duration.ToFloat())
	result += fmt.Sprintf("Video: %dx%d pixel @%f sample rate\n", header.VideoWidth, header.VideoHeight, header.AudioSampleRate.ToFloat())
	return
}
