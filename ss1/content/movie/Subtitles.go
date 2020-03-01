package movie

// Subtitles describe the textual representation of a movie.
type Subtitles struct {
	Entries []SubtitleEntry
}

// SubtitleEntry is a timestamped text for subtitles.
type SubtitleEntry struct {
	Timestamp Timestamp
	Text      string
}

func (sub *Subtitles) add(ts Timestamp, text string) {
	sub.Entries = append(sub.Entries, SubtitleEntry{
		Timestamp: ts,
		Text:      text,
	})
}
