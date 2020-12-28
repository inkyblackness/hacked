package compression

import "math/bits"

// TilePaletteKey marks use of colors of a 256 palette.
type TilePaletteKey struct {
	usedColors [4]uint64
	size       int
}

// TilePaletteKeyFrom returns a key for the given slice of colors.
func TilePaletteKeyFrom(colors []byte) TilePaletteKey {
	var key TilePaletteKey
	for _, c := range colors {
		key.UseColor(c)
	}
	return key
}

// Buffer returns the color indices marked in use by this key.
func (key *TilePaletteKey) Buffer() []byte {
	result := make([]byte, 0, key.size)
	for i := 0; i < 256; i++ {
		if key.HasColor(byte(i)) {
			result = append(result, byte(i))
		}
	}
	return result
}

// UseColor marks the given color index in use.
func (key *TilePaletteKey) UseColor(index byte) {
	if !key.HasColor(index) {
		key.usedColors[index/64] |= 1 << uint(index%64)
		key.size++
	}
}

// HasColor returns true if the given color index is in use by this key.
func (key *TilePaletteKey) HasColor(index byte) bool {
	return (key.usedColors[index/64] & (1 << uint(index%64))) != 0
}

// Contains returns true if this key is equal to, or a superset of, the given key.
func (key *TilePaletteKey) Contains(other *TilePaletteKey) bool {
	return ((^key.usedColors[0] & other.usedColors[0]) == 0) &&
		((^key.usedColors[1] & other.usedColors[1]) == 0) &&
		((^key.usedColors[2] & other.usedColors[2]) == 0) &&
		((^key.usedColors[3] & other.usedColors[3]) == 0)
}

// Without returns a new key instance that has all remaining colors marked in use.
func (key *TilePaletteKey) Without(other *TilePaletteKey) TilePaletteKey {
	var result TilePaletteKey
	for i := 0; i < 4; i++ {
		result.usedColors[i] = key.usedColors[i] & ^other.usedColors[i]
		result.size += bits.OnesCount64(result.usedColors[i])
	}
	return result
}

// LessThan is a sorting function to order keys.
func (key *TilePaletteKey) LessThan(other *TilePaletteKey) bool {
	for i := 0; i < 4; i++ {
		if key.usedColors[i] < other.usedColors[i] {
			return true
		}
	}
	return false
}
