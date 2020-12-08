package world

import (
	"path/filepath"
	"strings"
)

// FileLocation describes the path and name of a file.
type FileLocation struct {
	DirPath string
	Name    string
}

// FileLocationFrom returns an instance for given filename.
func FileLocationFrom(filename string) FileLocation {
	return FileLocation{
		DirPath: filepath.Dir(filename),
		Name:    filepath.Base(filename),
	}
}

// AbsolutePathFrom returns the complete path, based on given base, of the location.
func (loc FileLocation) AbsolutePathFrom(base string) string {
	rel, err := filepath.Rel(base, loc.DirPath)
	if err != nil {
		return filepath.Join(base, loc.DirPath, loc.Name)
	}
	return filepath.Join(base, rel, loc.Name)
}

// NestedRelativeTo returns a string that describes the relative nested location to given base - if possible.
func (loc FileLocation) NestedRelativeTo(base string) string {
	absoluteDir := loc.DirPath
	if !filepath.IsAbs(absoluteDir) {
		absoluteDir = filepath.Clean(filepath.Join(base, loc.DirPath))
	}
	if !strings.HasPrefix(absoluteDir+string(filepath.Separator), base+string(filepath.Separator)) {
		return filepath.Join(absoluteDir, loc.Name)
	}
	rel, err := filepath.Rel(base, absoluteDir)
	if err != nil {
		return filepath.Join(base, loc.Name)
	}
	return filepath.Join(".", rel, loc.Name)
}
