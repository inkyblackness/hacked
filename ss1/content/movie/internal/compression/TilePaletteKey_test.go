package compression

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTilePaletteKeyNullValueHasNoColorsSet(t *testing.T) {
	var key tilePaletteKey
	for color := 0; color < 256; color++ {
		assert.False(t, key.hasColor(byte(color)), "color set although should not %v", color)
	}
}

func TestTilePaletteKeyCanMarkSingleColors(t *testing.T) {
	for color := 0; color < 256; color++ {
		var key tilePaletteKey
		key.useColor(byte(color))
		assert.True(t, key.hasColor(byte(color)), "color not set %v", color)
		for otherColor := 0; otherColor < 256; otherColor++ {
			if otherColor != color {
				assert.False(t, key.hasColor(byte(otherColor)),
					"color %v set although should not for %v", otherColor, color)
			}
		}
	}
}

func TestTilePaletteKeyCanMarkAllColors(t *testing.T) {
	var key tilePaletteKey
	for color := 0; color < 256; color++ {
		key.useColor(byte(color))
	}
	for color := 0; color < 256; color++ {
		assert.True(t, key.hasColor(byte(color)), "should have color set %v", color)
	}
}

func TestTilePaletteKeyContains(t *testing.T) {
	baseKey := tilePaletteKeyFrom([]byte{1, 2, 3, 4})
	tt := []struct {
		other    []byte
		expected bool
	}{
		{other: []byte{1, 2, 3, 4}, expected: true},
		{other: []byte{2, 3, 4}, expected: true},
		{other: []byte{1, 2, 4}, expected: true},
		{other: []byte{4}, expected: true},
		{other: []byte{}, expected: true},
		{other: []byte{1, 2, 3, 4, 5}, expected: false},
		{other: []byte{3, 4, 5}, expected: false},
		{other: []byte{5}, expected: false},
	}

	for index, tc := range tt {
		td := tc
		t.Run(fmt.Sprintf("%v", index), func(t *testing.T) {
			other := tilePaletteKeyFrom(td.other)
			result := baseKey.contains(&other)
			assert.Equal(t, result, td.expected)
		})
	}
}

func TestTilePaletteKeyWithout(t *testing.T) {
	baseKey := tilePaletteKeyFrom([]byte{1, 2, 3, 4})
	tt := []struct {
		other    []byte
		expected []byte
	}{
		{other: []byte{1, 2, 3, 4}, expected: []byte{}},
		{other: []byte{2, 3, 4}, expected: []byte{1}},
		{other: []byte{1, 2, 4}, expected: []byte{3}},
		{other: []byte{4}, expected: []byte{1, 2, 3}},
		{other: []byte{}, expected: []byte{1, 2, 3, 4}},
		{other: []byte{1, 2, 3, 4, 5}, expected: []byte{}},
		{other: []byte{3, 4, 5}, expected: []byte{1, 2}},
		{other: []byte{5}, expected: []byte{1, 2, 3, 4}},
	}

	for index, tc := range tt {
		td := tc
		t.Run(fmt.Sprintf("%v", index), func(t *testing.T) {
			other := tilePaletteKeyFrom(td.other)
			result := baseKey.without(&other)
			expected := tilePaletteKeyFrom(td.expected)
			assert.Equal(t, result, expected, "Mismatch for %v", td.other)
		})
	}
}
