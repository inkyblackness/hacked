package ss1

// StringError is a string-based type implementing the error interface.
// It is used to cover all errors based on static strings.
type StringError string

// Error returns the string interpretation of the error.
func (err StringError) Error() string {
	return string(err)
}
