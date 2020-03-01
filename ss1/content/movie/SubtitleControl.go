package movie

import "github.com/inkyblackness/hacked/ss1/resource"

// SubtitleControl specifies how to interpret a subtitle entry.
type SubtitleControl uint32

// Subtitle constants
const (
	SubtitleArea    = SubtitleControl(0x41455241)
	SubtitleTextStd = SubtitleControl(0x20445453)
	SubtitleTextFrn = SubtitleControl(0x204E5246)
	SubtitleTextGer = SubtitleControl(0x20524547)
)

// String returns the string presentation of the control value.
func (ctrl SubtitleControl) String() string {
	return string([]rune{rune((ctrl >> 0) & 0xFF), rune((ctrl >> 8) & 0xFF), rune((ctrl >> 16) & 0xFF), rune((ctrl >> 24) & 0xFF)})
}

// SubtitleControlForLanguage returns the corresponding subtitle control for given language.
func SubtitleControlForLanguage(lang resource.Language) SubtitleControl {
	switch lang {
	case resource.LangDefault:
		return SubtitleTextStd
	case resource.LangFrench:
		return SubtitleTextFrn
	case resource.LangGerman:
		return SubtitleTextGer
	default:
		panic("unknown language")
	}
}
