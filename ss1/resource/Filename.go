package resource

import "strings"

// Filename defines a wrapper for a file that is possibly language-specific.
type Filename interface {
	// For returns the name of the file for the given language.
	For(lang Language) string

	// Matches returns true if the given filename matches the described one.
	Matches(filename string) bool
}

// I18nFile is for internationalized resource files - i.e., those that store resources per file.
type I18nFile [LanguageCount]string

// For returns the string per language index.
func (spec I18nFile) For(lang Language) string {
	return spec[int(lang)]
}

// Matches returns true if the given filename matches one of the localized filenames.
func (spec I18nFile) Matches(filename string) bool {
	lowercase := strings.ToLower(filename)
	for _, entry := range spec {
		if entry == lowercase {
			return true
		}
	}
	return false
}

// AnyLanguage is for generic resource files that are language agnostic.
type AnyLanguage string

// For returns the string itself.
func (any AnyLanguage) For(lang Language) string {
	return string(any)
}

// Matches returns true if the given filename matches this one.
func (any AnyLanguage) Matches(filename string) bool {
	return strings.ToLower(filename) == string(any)
}

// FilenameList is a list of filenames
type FilenameList []Filename

// Matches returns true if the given filename matches any of the contained entries.
func (list FilenameList) Matches(filename string) bool {
	for _, entry := range list {
		if entry.Matches(filename) {
			return true
		}
	}
	return false
}
