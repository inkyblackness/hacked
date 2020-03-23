package movie

import "github.com/inkyblackness/hacked/ss1/content/movie/internal/format"

// Container wraps the information and data of a MOVI container.
type Container struct {
	Audio     Audio
	Video     Video
	Subtitles Subtitles
}

// Duration returns the length of the entire movie. It is the highest timestamp of all media streams.
func (container Container) Duration() format.Timestamp {
	var max format.Timestamp
	for _, ts := range []format.Timestamp{container.Audio.Duration(), container.Video.Duration(), container.Subtitles.Duration()} {
		if max.IsBefore(ts) {
			max = ts
		}
	}
	return max
}
