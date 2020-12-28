package compression_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/inkyblackness/hacked/ss1/content/movie/internal/compression"
)

func TestTilePaletteKeyNullValueHasNoColorsSet(t *testing.T) {
	var key compression.TilePaletteKey
	for color := 0; color < 256; color++ {
		assert.False(t, key.HasColor(byte(color)), "color set although should not %v", color)
	}
}

func TestTilePaletteKeyCanMarkSingleColors(t *testing.T) {
	for color := 0; color < 256; color++ {
		var key compression.TilePaletteKey
		key.UseColor(byte(color))
		assert.True(t, key.HasColor(byte(color)), "color not set %v", color)
		for otherColor := 0; otherColor < 256; otherColor++ {
			if otherColor != color {
				assert.False(t, key.HasColor(byte(otherColor)),
					"color %v set although should not for %v", otherColor, color)
			}
		}
	}
}

func TestTilePaletteKeyCanMarkAllColors(t *testing.T) {
	var key compression.TilePaletteKey
	for color := 0; color < 256; color++ {
		key.UseColor(byte(color))
	}
	for color := 0; color < 256; color++ {
		assert.True(t, key.HasColor(byte(color)), "should have color set %v", color)
	}
}

func TestTilePaletteKeyContains(t *testing.T) {
	baseKey := compression.TilePaletteKeyFrom([]byte{1, 2, 3, 4})
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
			other := compression.TilePaletteKeyFrom(td.other)
			result := baseKey.Contains(&other)
			assert.Equal(t, result, td.expected)
		})
	}
}

func TestTilePaletteKeyWithout(t *testing.T) {
	baseKey := compression.TilePaletteKeyFrom([]byte{1, 2, 3, 4})
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
			other := compression.TilePaletteKeyFrom(td.other)
			result := baseKey.Without(&other)
			expected := compression.TilePaletteKeyFrom(td.expected)
			assert.Equal(t, result, expected, "Mismatch for %v", td.other)
		})
	}
}
