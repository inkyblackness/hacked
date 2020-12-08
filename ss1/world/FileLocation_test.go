package world_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/inkyblackness/hacked/ss1/world"
)

func TestFileLocationAbsolutePathFrom(t *testing.T) {
	temp := os.TempDir()
	base := filepath.Join(temp, "base")

	tt := []struct {
		name     string
		loc      world.FileLocation
		base     string
		expected string
	}{
		{
			name:     "file loaded in same base",
			loc:      world.FileLocation{DirPath: base, Name: "test.res"},
			base:     base,
			expected: filepath.Join(base, "test.res"),
		},
		{
			name:     "new file",
			loc:      world.FileLocation{DirPath: ".", Name: "test.res"},
			base:     base,
			expected: filepath.Join(base, "test.res"),
		},
		{
			name:     "new file without dir path",
			loc:      world.FileLocation{DirPath: "", Name: "test.res"},
			base:     base,
			expected: filepath.Join(base, "test.res"),
		},
		{
			name:     "file found in different sibling path",
			loc:      world.FileLocation{DirPath: filepath.Join(temp, "other"), Name: "test.res"},
			base:     base,
			expected: filepath.Join(temp, "other", "test.res"),
		},
		{
			name:     "file found in sub path",
			loc:      world.FileLocation{DirPath: filepath.Join(base, "nested"), Name: "test.res"},
			base:     base,
			expected: filepath.Join(base, "nested", "test.res"),
		},
		{
			name:     "nested relative file dir path",
			loc:      world.FileLocation{DirPath: "nested", Name: "test.res"},
			base:     base,
			expected: filepath.Join(base, "nested", "test.res"),
		},
	}
	for _, tc := range tt {
		func(name string, loc world.FileLocation, base, expected string) {
			t.Run(name, func(t *testing.T) {
				abs := loc.AbsolutePathFrom(base)
				assert.Equal(t, expected, abs)
			})
		}(tc.name, tc.loc, tc.base, tc.expected)
	}
}

func TestFileLocationNestedRelativeTo(t *testing.T) {
	temp := os.TempDir()
	base := filepath.Join(temp, "base")

	tt := []struct {
		name     string
		loc      world.FileLocation
		base     string
		expected string
	}{
		{
			name:     "file loaded in same base",
			loc:      world.FileLocation{DirPath: base, Name: "test.res"},
			base:     base,
			expected: filepath.Join(".", "test.res"),
		},
		{
			name:     "new file",
			loc:      world.FileLocation{DirPath: ".", Name: "test.res"},
			base:     base,
			expected: filepath.Join(".", "test.res"),
		},
		{
			name:     "new file without dir path",
			loc:      world.FileLocation{DirPath: "", Name: "test.res"},
			base:     base,
			expected: filepath.Join(".", "test.res"),
		},
		{
			name:     "file found in different path",
			loc:      world.FileLocation{DirPath: filepath.Join(temp, "other"), Name: "test.res"},
			base:     base,
			expected: filepath.Join(temp, "other", "test.res"),
		},
		{
			name:     "file found in sub path",
			loc:      world.FileLocation{DirPath: filepath.Join(base, "nested"), Name: "test.res"},
			base:     base,
			expected: filepath.Join(".", "nested", "test.res"),
		},
		{
			name:     "file found in similar named different path",
			loc:      world.FileLocation{DirPath: filepath.Join(temp, "baseOther"), Name: "test.res"},
			base:     base,
			expected: filepath.Join(temp, "baseOther", "test.res"),
		},
	}
	for _, tc := range tt {
		func(name string, loc world.FileLocation, base, expected string) {
			t.Run(name, func(t *testing.T) {
				rel := loc.NestedRelativeTo(base)
				assert.Equal(t, expected, rel)
			})
		}(tc.name, tc.loc, tc.base, tc.expected)
	}
}
