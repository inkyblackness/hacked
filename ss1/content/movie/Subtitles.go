package movie

import (
	"time"

	"github.com/inkyblackness/hacked/ss1/content/movie/internal/format"
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/resource"
)

// Subtitles contains all the subtitles in all languages.
type Subtitles struct {
	PerLanguage [resource.LanguageCount]SubtitleList
}

func (sub *Subtitles) add(lang resource.Language, timestamp time.Duration, text string) {
	sub.PerLanguage[lang].Entries = append(sub.PerLanguage[lang].Entries, Subtitle{
		Timestamp: timestamp,
		Text:      text,
	})
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

func (sub Subtitles) duration() format.Timestamp {
	var highest time.Duration
	for _, list := range sub.PerLanguage {
		for _, sub := range list.Entries {
			if highest < sub.Timestamp {
				highest = sub.Timestamp
			}
		}
	}
	return format.TimestampFromDuration(highest)
}

func (sub Subtitles) encode(cp text.Codepage) [][]format.EntryBucket {
	if !sub.ArePresent() {
		return nil
	}
	bucketsPerLanguage := make([][]format.EntryBucket, len(sub.PerLanguage)+1)

	// Ensure a subtitle area is defined.
	// The area is hardcoded. While the engine respects any area, placing the text in the
	// frame area will have the pixels become overwritten. As such, there are many "wrong" options,
	// and only a few right ones. There's no need to make them editable.
	bucketsPerLanguage[0] = []format.EntryBucket{{
		Priority:  format.EntryBucketPrioritySubtitleControl,
		Timestamp: format.Timestamp{},
		Entries: []format.Entry{{
			Timestamp: format.Timestamp{},
			Data: format.SubtitleEntryData{
				Control: format.SubtitleArea,
				Text:    cp.Encode("20 365 620 395 CLR"),
			},
		}},
	}}
	for index, lang := range sub.PerLanguage {
		bucketsPerLanguage[index+1] = lang.encode(format.SubtitleControlForLanguage(resource.Language(index)), cp)
	}
	return bucketsPerLanguage
}

// SubtitleList describes the textual representation of a movie in one language.
type SubtitleList struct {
	Entries []Subtitle
}

func (sub SubtitleList) encode(control format.SubtitleControl, cp text.Codepage) []format.EntryBucket {
	buckets := make([]format.EntryBucket, 0, len(sub.Entries))
	for _, entry := range sub.Entries {
		buckets = append(buckets, entry.encode(control, cp))
	}
	return buckets
}

// Subtitle is a timestamped text for subtitles.
type Subtitle struct {
	Timestamp time.Duration
	Text      string
}

func (sub Subtitle) encode(control format.SubtitleControl, cp text.Codepage) format.EntryBucket {
	return format.EntryBucket{
		Priority:  format.EntryBucketPrioritySubtitle,
		Timestamp: format.TimestampFromDuration(sub.Timestamp),
		Entries: []format.Entry{{
			Timestamp: format.TimestampFromDuration(sub.Timestamp),
			Data: format.SubtitleEntryData{
				Control: control,
				Text:    cp.Encode(sub.Text),
			},
		}},
	}
}
