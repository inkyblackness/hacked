package movie

// Container wraps the information and data of a MOVI container.
type Container struct {
	Audio     Audio
	Video     Video
	Subtitles Subtitles
}

// Duration returns the length of the entire movie. It is the highest timestamp of all media streams.
func (container Container) Duration() Timestamp {
	var max Timestamp
	for _, ts := range []Timestamp{container.Audio.Duration(), container.Video.Duration(), container.Subtitles.Duration()} {
		if max.IsBefore(ts) {
			max = ts
		}
	}
	return max
}
