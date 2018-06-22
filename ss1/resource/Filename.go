package resource

// Filename defines a wrapper for a file that is possibly language-specific.
type Filename interface {
	// For returns the name of the file for the given language.
	For(lang Language) string
}

// I18nFile is for internationalized resource files - i.e., those that store resources per file.
type I18nFile [LanguageCount]string

// For returns the string per language index.
func (spec I18nFile) For(lang Language) string {
	return spec[int(lang)]
}

// AnyLanguage is for generic resource files that are language agnostic.
type AnyLanguage string

// For returns the string itself.
func (any AnyLanguage) For(lang Language) string {
	return string(any)
}
