package compression

import (
	"math/bits"
	"sort"
)

type tilePaletteKey struct {
	usedColors [4]uint64
	size       int
}

func (entry *tilePaletteKey) buffer() []byte {
	result := make([]byte, 0, entry.size)
	for i := 0; i < 256; i++ {
		if entry.hasColor(byte(i)) {
			result = append(result, byte(i))
		}
	}
	return result
}

func (entry *tilePaletteKey) joinedBuffer(source []byte) []byte {
	result := make([]byte, 0, entry.size)
	var addedColors tilePaletteKey
	for _, color := range source {
		addedColors.useColor(color)
		result = append(result, color)
	}
	for color := 0; color < 256; color++ {
		if entry.hasColor(byte(color)) && !addedColors.hasColor(byte(color)) {
			result = append(result, byte(color))
		}
	}
	return result
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
	index, inLookup := lookup.starts[key]
	if inLookup {
		pal = lookup.buffer[index : index+int(key.size)]
	} else {
		pal = key.buffer()
	}
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

	remainder := make(map[tilePaletteKey]struct{})
	for key := range gen.keyUses {
		remainder[key] = struct{}{}
	}

	for size := PixelPerTile; size > 2; size-- {
		var keysInSize []tilePaletteKey

		{ // TODO: consider removing this block again should it not bring too much of a benefit.
			var earlyRemoved []tilePaletteKey
			for key := range remainder {
				if key.size == size {
					wasRemoved := false
					for start := 0; start < (len(lookup.buffer)-key.size) && !wasRemoved; start++ {
						var tempKey tilePaletteKey
						for _, color := range lookup.buffer[start : start+key.size] {
							tempKey.useColor(color)
						}
						if tempKey.contains(&key) {
							earlyRemoved = append(earlyRemoved, key)
							wasRemoved = true

							lookup.starts[key] = start
						}
					}
				}
			}
			for _, key := range earlyRemoved {
				delete(remainder, key)
			}
		}

		// find all keys with this current size
		for key := range remainder {
			if key.size == size {
				keysInSize = append(keysInSize, key)
			}
		}

		toRemove := keysInSize[:]
		for _, key := range keysInSize {
			var containedKeys []tilePaletteKey

			// find all contained keys, sort them by usage
			for nestedKey := range remainder {
				if (nestedKey.size < key.size) && key.contains(&nestedKey) {
					containedKeys = append(containedKeys, nestedKey)
				}
			}

			sort.Slice(containedKeys, func(a, b int) bool {
				// return gen.keyUses[containedKeys[a]] > gen.keyUses[containedKeys[b]] // sort by use has no point
				return containedKeys[a].size > containedKeys[b].size
			})

			lookup.starts[key] = len(lookup.buffer)
			if len(containedKeys) > 0 {
				containedKey := containedKeys[0]
				toRemove = append(toRemove, containedKey)
				lookup.starts[containedKey] = len(lookup.buffer)
				lookup.buffer = append(lookup.buffer, key.joinedBuffer(containedKey.buffer())...)
			} else {
				lookup.buffer = append(lookup.buffer, key.buffer()...)
			}
		}
		for _, key := range toRemove {
			delete(remainder, key)
		}
	}

	for key := range remainder {
		lookup.starts[key] = len(lookup.buffer)
		lookup.buffer = append(lookup.buffer, key.buffer()...)
	}

	return lookup
}

// Add registers a further delta to the generator.
func (gen *PaletteLookupGenerator) Add(delta tileDelta) {
	var key tilePaletteKey
	for _, pal := range delta {
		key.useColor(pal)
	}
	if key.size > 2 {
		if gen.keyUses == nil {
			gen.keyUses = make(map[tilePaletteKey]int)
		}
		gen.keyUses[key]++
	}
}
