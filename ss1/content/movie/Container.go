package movie

import "github.com/inkyblackness/hacked/ss1/content/movie/internal/format"

// Container wraps the information and data of a MOVI container.
type Container struct {
	Audio     Audio
	Video     Video
	Subtitles Subtitles
}

func (container Container) duration() format.Timestamp {
	var max format.Timestamp
	for _, ts := range []format.Timestamp{container.Audio.duration(), container.Video.duration(), container.Subtitles.duration()} {
		if max.IsBefore(ts) {
			max = ts
		}
	}
	return max
}
