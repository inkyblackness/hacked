package compression

import "math/bits"

type tilePaletteKey struct {
	usedColors [4]uint64
	size       int
}

func tilePaletteKeyFrom(colors []byte) tilePaletteKey {
	var key tilePaletteKey
	for _, c := range colors {
		key.useColor(c)
	}
	return key
}

func (key *tilePaletteKey) buffer() []byte {
	result := make([]byte, 0, key.size)
	for i := 0; i < 256; i++ {
		if key.hasColor(byte(i)) {
			result = append(result, byte(i))
		}
	}
	return result
}

func (key *tilePaletteKey) useColor(index byte) {
	if !key.hasColor(index) {
		key.usedColors[index/64] |= 1 << uint(index%64)
		key.size++
	}
}

func (key *tilePaletteKey) hasColor(index byte) bool {
	return (key.usedColors[index/64] & (1 << uint(index%64))) != 0
}

func (key *tilePaletteKey) contains(other *tilePaletteKey) bool {
	return ((^key.usedColors[0] & other.usedColors[0]) == 0) &&
		((^key.usedColors[1] & other.usedColors[1]) == 0) &&
		((^key.usedColors[2] & other.usedColors[2]) == 0) &&
		((^key.usedColors[3] & other.usedColors[3]) == 0)
}

func (key *tilePaletteKey) without(other *tilePaletteKey) tilePaletteKey {
	var result tilePaletteKey
	for i := 0; i < 4; i++ {
		result.usedColors[i] = key.usedColors[i] & ^other.usedColors[i]
		result.size += bits.OnesCount64(result.usedColors[i])
	}
	return result
}

func (key *tilePaletteKey) lessThan(other *tilePaletteKey) bool {
	for i := 0; i < 4; i++ {
		if key.usedColors[i] < other.usedColors[i] {
			return true
		}
	}
	return false
}
