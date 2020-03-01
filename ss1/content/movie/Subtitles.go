package movie

// Subtitles describe the textual representation of a movie.
type Subtitles struct {
	entries []SubtitleEntry
}

// SubtitleEntry is a timestamped text for subtitles.
type SubtitleEntry struct {
	Timestamp Timestamp
	Text      string
}

func (sub *Subtitles) add(ts Timestamp, text string) {
	sub.entries = append(sub.entries, SubtitleEntry{
		Timestamp: ts,
		Text:      text,
	})
}
