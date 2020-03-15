package movie

import (
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/resource"
)

type AllSubtitles struct {
	PerLanguage [resource.LanguageCount]Subtitles
}

func (all AllSubtitles) AreInUse() bool {
	for _, lang := range all.PerLanguage {
		if len(lang.Entries) > 0 {
			return true
		}
	}
	return false
}

func (all AllSubtitles) Encode(cp text.Codepage) [][]EntryBucket {
	if !all.AreInUse() {
		return nil
	}
	bucketsPerLanguage := make([][]EntryBucket, len(all.PerLanguage)+1)

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
				Control: SubtitleArea,
				Text:    cp.Encode("20 365 620 395 CLR"),
			},
		}},
	}}
	for index, lang := range all.PerLanguage {
		bucketsPerLanguage[index+1] = lang.Encode(SubtitleControlForLanguage(resource.Language(index)), cp)
	}
	return bucketsPerLanguage
}

// Subtitles describe the textual representation of a movie.
type Subtitles struct {
	Entries []Subtitle
}

func (sub Subtitles) Encode(control SubtitleControl, cp text.Codepage) []EntryBucket {
	buckets := make([]EntryBucket, 0, len(sub.Entries))
	for _, entry := range sub.Entries {
		buckets = append(buckets, entry.Encode(control, cp))
	}
	return buckets
}

func (sub *Subtitles) add(ts Timestamp, text string) {
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

func (sub Subtitle) Encode(control SubtitleControl, cp text.Codepage) EntryBucket {
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
