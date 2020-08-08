package world

import "path/filepath"

// FileLocation describes the path and name of a file.
type FileLocation struct {
	DirPath string
	Name    string
}

// AbsolutePathFrom returns the complete path, based on given base, of the location.
func (loc FileLocation) AbsolutePathFrom(base string) string {
	rel, err := filepath.Rel(base, loc.DirPath)
	if err != nil {
		return filepath.Join(base, loc.Name)
	}
	return filepath.Join(base, rel, loc.Name)
}
