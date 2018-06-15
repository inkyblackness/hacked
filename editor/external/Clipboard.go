package external

// Clipboard represents a temporary storage.
type Clipboard interface {
	// String returns the current value of the clipboard, if it is compatible with UTF-8.
	String() (string, error)
	// SetString sets the current value of the clipboard as UTF-8 string.
	SetString(value string)
}
