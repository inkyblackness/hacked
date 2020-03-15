package movie

// Subtitles describe the textual representation of a movie.
type Subtitles struct {
	Entries []Subtitle
}

// Subtitle is a timestamped text for subtitles.
type Subtitle struct {
	Timestamp Timestamp
	Text      string
}

func (sub *Subtitles) add(ts Timestamp, text string) {
	sub.Entries = append(sub.Entries, Subtitle{
		Timestamp: ts,
		Text:      text,
	})
}
