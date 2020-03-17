package movie

import (
	"github.com/inkyblackness/hacked/ss1/content/movie/internal/format"
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/resource"
)

// Subtitles contains all the subtitles in all languages.
type Subtitles struct {
	PerLanguage [resource.LanguageCount]SubtitleList
}

// Add appends the given text at given timestamp in given language.
func (sub *Subtitles) Add(lang resource.Language, timestamp Timestamp, text string) {
	sub.PerLanguage[lang].Entries = append(sub.PerLanguage[lang].Entries, Subtitle{
		Timestamp: timestamp,
		Text:      text,
	})
}

// Duration returns the highest timestamp of all the subtitles
func (sub Subtitles) Duration() Timestamp {
	var highest Timestamp
	for _, list := range sub.PerLanguage {
		for _, sub := range list.Entries {
			if highest.IsBefore(sub.Timestamp) {
				highest = sub.Timestamp
			}
		}
	}
	return highest
}

// ArePresent returns true if at least one language makes use of subtitles.
func (sub Subtitles) ArePresent() bool {
	for _, lang := range sub.PerLanguage {
		if len(lang.Entries) > 0 {
			return true
		}
	}
	return false
}

// Encode serializes all the subtitles of all languages into buckets.
func (sub Subtitles) Encode(cp text.Codepage) [][]EntryBucket {
	if !sub.ArePresent() {
		return nil
	}
	bucketsPerLanguage := make([][]EntryBucket, len(sub.PerLanguage)+1)

	// Ensure a subtitle area is defined.
	// The area is hardcoded. While the engine respects any area, placing the text in the
	// frame area will have the pixels become overwritten. As such, there are many "wrong" options,
	// and only a few right ones. There's no need to make them editable.
	bucketsPerLanguage[0] = []EntryBucket{{
		Priority:  EntryBucketPrioritySubtitleControl,
		Timestamp: Timestamp{},
		Entries: []Entry{{
			Timestamp: Timestamp{},
			Data: SubtitleEntryData{
				Control: format.SubtitleArea,
				Text:    cp.Encode("20 365 620 395 CLR"),
			},
		}},
	}}
	for index, lang := range sub.PerLanguage {
		bucketsPerLanguage[index+1] = lang.Encode(format.SubtitleControlForLanguage(resource.Language(index)), cp)
	}
	return bucketsPerLanguage
}

// SubtitleList describes the textual representation of a movie in one language.
type SubtitleList struct {
	Entries []Subtitle
}

// Encode serializes the list of subtitles into buckets.
func (sub SubtitleList) Encode(control format.SubtitleControl, cp text.Codepage) []EntryBucket {
	buckets := make([]EntryBucket, 0, len(sub.Entries))
	for _, entry := range sub.Entries {
		buckets = append(buckets, entry.Encode(control, cp))
	}
	return buckets
}

func (sub *SubtitleList) add(ts Timestamp, text string) {
	sub.Entries = append(sub.Entries, Subtitle{
		Timestamp: ts,
		Text:      text,
	})
}

// Subtitle is a timestamped text for subtitles.
type Subtitle struct {
	Timestamp Timestamp
	Text      string
}

// Encode serializes the subtitle into a bucket.
func (sub Subtitle) Encode(control format.SubtitleControl, cp text.Codepage) EntryBucket {
	return EntryBucket{
		Priority:  EntryBucketPrioritySubtitle,
		Timestamp: sub.Timestamp,
		Entries: []Entry{{
			Timestamp: sub.Timestamp,
			Data: SubtitleEntryData{
				Control: control,
				Text:    cp.Encode(sub.Text),
			},
		}},
	}
}
