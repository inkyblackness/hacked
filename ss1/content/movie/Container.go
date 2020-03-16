package movie

import (
	"github.com/inkyblackness/hacked/ss1/content/bitmap"
)

// Container wraps the information and data of a MOVI container.
type Container struct {
	// EndTimestamp is the time of the end of the movie.
	EndTimestamp Timestamp

	// VideoWidth is the width of a video in pixel.
	VideoWidth uint16
	// VideoHeight is the height of a video in pixel.
	VideoHeight uint16
	// StartPalette is the initial palette of a video.
	StartPalette bitmap.Palette

	// AudioSampleRate is the sample frequency used for audio entries.
	AudioSampleRate uint16

	// Entries are all the parts of the movie.
	Entries []Entry

	// TODO: fill by Read()
	// TODO: merge by Write(), ordered by bucket priority
	// TODO: remove other members, they should all no longer be necessary in the end.
	Audio     Audio
	Video     Video
	Subtitles Subtitles
}

// AddEntry adds the given entry to the existing list.
func (container *Container) AddEntry(entry Entry) {
	container.Entries = append(container.Entries, entry)
}
