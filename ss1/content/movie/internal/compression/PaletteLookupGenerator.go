package compression

import "math/bits"

type tilePaletteKey struct {
	usedColors [4]uint64
	size       byte
}

func (entry *tilePaletteKey) useColor(index byte) {
	if !entry.hasColor(index) {
		entry.usedColors[index/64] |= 1 << uint(index%64)
		entry.size++
	}
}

func (entry *tilePaletteKey) hasColor(index byte) bool {
	return (entry.usedColors[index/64] & (1 << uint(index%64))) != 0
}

func (entry *tilePaletteKey) contains(other *tilePaletteKey) bool {
	return ((^entry.usedColors[0] & other.usedColors[0]) == 0) &&
		((^entry.usedColors[1] & other.usedColors[1]) == 0) &&
		((^entry.usedColors[2] & other.usedColors[2]) == 0) &&
		((^entry.usedColors[3] & other.usedColors[3]) == 0)
}

// PaletteLookup is a dictionary of tile delta data, found in a palette buffer.
type PaletteLookup struct {
	buffer []byte
	starts map[tilePaletteKey]int
}

// Buffer returns the underlying slice.
func (lookup *PaletteLookup) Buffer() []byte {
	return lookup.buffer
}

// Lookup finds the given tile again and returns the properties where and how to reproduce it.
func (lookup *PaletteLookup) Lookup(tile tileDelta) (index int, pal []byte, mask uint64) {
	var key tilePaletteKey
	for _, pal := range tile {
		key.useColor(pal)
	}
	index = lookup.starts[key]
	pal = lookup.buffer[index : index+int(key.size)]
	var mapped [256]byte
	for mappedIndex, b := range pal {
		mapped[b] = byte(mappedIndex)
	}
	bitSize := uint(bits.Len(uint(key.size - 1)))
	for tileIndex := PixelPerTile - 1; tileIndex >= 0; tileIndex-- {
		mask <<= bitSize
		mask |= uint64(mapped[tile[tileIndex]])
	}
	return
}

// PaletteLookupGenerator creates palette lookups based on a set of registered tiles.
type PaletteLookupGenerator struct {
	// deltaToKey map[tileDelta]tilePaletteKey
	keyUses map[tilePaletteKey]int
}

// Generate creates a lookup based on all currently registered tile deltas.
func (gen *PaletteLookupGenerator) Generate() PaletteLookup {
	var lookup PaletteLookup
	lookup.starts = make(map[tilePaletteKey]int)
	for key := range gen.keyUses {
		lookup.starts[key] = len(lookup.buffer)
		for i := 0; i < 256; i++ {
			if key.hasColor(byte(i)) {
				lookup.buffer = append(lookup.buffer, byte(i))
			}
		}
	}
	return lookup
}

// Add registers a further delta to the generator.
func (gen *PaletteLookupGenerator) Add(delta tileDelta) {
	var key tilePaletteKey
	for _, pal := range delta {
		key.useColor(pal)
	}
	if gen.keyUses == nil {
		gen.keyUses = make(map[tilePaletteKey]int)
	}
	gen.keyUses[key]++
}
