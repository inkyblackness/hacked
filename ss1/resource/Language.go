package resource

import "fmt"

// Language defines the human language of a resource.
type Language byte

const (
	// LangAny identifies language agnostic resources.
	LangAny Language = 0xFF
	// LangDefault identifies the default language, typically English - unless modded.
	LangDefault Language = 0
	// LangFrench identifies the French language.
	LangFrench Language = 1
	// LangGerman identifies the German language.
	LangGerman Language = 2
)

func (lang Language) String() string {
	switch lang {
	case LangAny:
		return "Any"
	case LangDefault:
		return "Default"
	case LangFrench:
		return "French"
	case LangGerman:
		return "German"
	default:
		return fmt.Sprintf("Unknown%02X", int(lang))
	}
}

// Languages returns a slice of all human languages. Does not include "Any" selector.
func Languages() []Language {
	return []Language{LangDefault, LangFrench, LangGerman}
}

// Includes returns true if the language includes the provided one.
// This is not symmetrical. While "Any" includes "German", "German" does not include "Any".
func (lang Language) Includes(other Language) bool {
	return (lang == LangAny) || (lang == other)
}
