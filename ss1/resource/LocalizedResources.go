package resource

import (
	"strings"
)

// LocalizedResources associates a language with a resource provider under a specific identifier.
type LocalizedResources struct {
	// ID is the identifier of the provider. This could be a filename for instance.
	ID string
	// Language specifies for which language the provider has resources.
	Language Language
	// Provider is the actual container of the resources.
	Provider Provider
}

type languageSpecificFilenames struct {
	cybstrng string
	mfdart   string
	citalog  string
	citbark  string
	lowintr  string
	svgaintr string
}

func (spec languageSpecificFilenames) hasFilename(filename string) bool {
	lowercase := strings.ToLower(filename)
	return (spec.cybstrng == lowercase) ||
		(spec.mfdart == lowercase) ||
		(spec.citalog == lowercase) ||
		(spec.citbark == lowercase) ||
		(spec.lowintr == lowercase) ||
		(spec.svgaintr == lowercase)
}

var localizedFilenames = map[Language]languageSpecificFilenames{
	LangDefault: {
		cybstrng: "cybstrng.res",
		mfdart:   "mfdart.res",
		citalog:  "citalog.res",
		citbark:  "citbark.res",
		lowintr:  "lowintr.res",
		svgaintr: "svgaintr.res",
	},
	LangFrench: {
		cybstrng: "frnstrng.res",
		mfdart:   "mfdfrn.res",
		citalog:  "frnalog.res",
		citbark:  "frnbark.res",
		lowintr:  "lofrintr.res",
		svgaintr: "svfrintr.res",
	},
	LangGerman: {
		cybstrng: "gerstrng.res",
		mfdart:   "mfdger.res",
		citalog:  "geralog.res",
		citbark:  "gerbark.res",
		lowintr:  "logeintr.res",
		svgaintr: "svgeintr.res",
	},
}

// LocalizeFilename returns the language that the resource file would typically contain.
func LocalizeFilename(filename string) Language {
	result := LangAny
	for lang, loc := range localizedFilenames {
		if loc.hasFilename(filename) {
			result = lang
		}
	}
	return result
}

// LocalizeResourcesByFilename creates an instance of LocalizedResources based on a filename and a provider.
// The given filename is taken as an ID, as well as a hint to identify the language.
func LocalizeResourcesByFilename(provider Provider, filename string) (res LocalizedResources) {
	res.ID = filename
	res.Provider = provider
	res.Language = LocalizeFilename(filename)

	return
}
