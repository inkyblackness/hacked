package format

import (
	"fmt"
)

// Tag is the identifier key of a MOVI container.
const Tag string = "MOVI"

// HeaderSize specifies the length of a serialized Header in bytes.
const HeaderSize = 256

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
	// DurationFraction is the length of the media in 1/0x10000 second units.
	DurationFraction uint16
	// DurationSeconds is the length of the media in seconds.
	DurationSeconds byte

	Unused0013 byte

	Unknown0014 [4]byte

	// VideoWidth specifies the width of the video.
	VideoWidth uint16
	// VideoHeight specifies the height of the video.
	VideoHeight uint16

	Unknown001C int16
	Unknown001E int16
	Unknown0020 int16
	Unknown0022 int32

	// SampleRate specifies the sample rate of the audio.
	SampleRate uint16

	Zero [216]byte
}

func (header *Header) String() (result string) {
	result += fmt.Sprintf("Index: %d entries (%d bytes)\n", header.IndexEntryCount, header.IndexSize)
	result += fmt.Sprintf("Content: %d bytes, %f seconds\n", header.ContentSize, float32(header.DurationSeconds)+float32(header.DurationFraction)/65536.0)
	result += fmt.Sprintf("Video: %dx%d pixel @%d sample rate\n", header.VideoWidth, header.VideoHeight, header.SampleRate)
	return
}
