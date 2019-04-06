package compression

import (
	"math/bits"
)

type tilePaletteKey struct {
	usedColors [4]uint64
	size       int
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

func (key *tilePaletteKey) joinedBuffer(source []byte) []byte {
	result := make([]byte, 0, key.size)
	var addedColors tilePaletteKey
	for _, color := range source {
		addedColors.useColor(color)
		result = append(result, color)
	}
	for color := 0; color < 256; color++ {
		if key.hasColor(byte(color)) && !addedColors.hasColor(byte(color)) {
			result = append(result, byte(color))
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
	for i := 0; i < 256; i++ {
		if key.hasColor(byte(i)) && !(*other).hasColor(byte(i)) {
			result.useColor(byte(i))
		}
	}
	return result
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

type nestedEntry struct {
	key    tilePaletteKey
	nested []nestedEntry
}

func (entry nestedEntry) buffer() []byte {
	return entry.extractBuffer(0, func(tilePaletteKey, int) {})
}

func (entry nestedEntry) byteSize() int {
	nestedSize := 0
	for _, nested := range entry.nested {
		nestedSize += nested.byteSize()
	}
	return entry.key.size + nestedSize
}

func (entry *nestedEntry) populate(keys map[tilePaletteKey]struct{}) {
	remainingKey := entry.key
	foundSomething := true
	for remainingKey.size > 2 && foundSomething {
		var lastAddedKey tilePaletteKey
		lastAddedKey, foundSomething = entry.populateRemaining(keys, remainingKey)
		remainingKey = remainingKey.without(&lastAddedKey)
	}
}

func (entry *nestedEntry) populateRemaining(keys map[tilePaletteKey]struct{}, remainingKey tilePaletteKey) (tilePaletteKey, bool) {
	maxByteSize := 0
	var maxNested *nestedEntry
	keySize := remainingKey.size
	for (keySize > 2) && (maxNested == nil) {
		keySize--
		for otherKey := range keys {
			if otherKey.size == keySize && remainingKey.contains(&otherKey) {
				nested := nestedEntry{key: otherKey}
				nested.populate(keys)
				nestedSize := nested.byteSize()
				if nestedSize > maxByteSize {
					maxByteSize = nestedSize
					maxNested = &nested
				}
			}
		}
	}
	if maxNested == nil {
		return tilePaletteKey{}, false
	}
	entry.nested = append(entry.nested, *maxNested)
	return maxNested.key, true
}

func (entry *nestedEntry) extractBuffer(startOffset int, marker func(tilePaletteKey, int)) []byte {
	var nestedBuffer []byte
	marker(entry.key, startOffset)
	relativeOffset := 0
	for _, nested := range entry.nested {
		bufferPart := nested.extractBuffer(startOffset+relativeOffset, marker)
		nestedBuffer = append(nestedBuffer, bufferPart...)
		relativeOffset += nested.key.size
	}
	return entry.key.joinedBuffer(nestedBuffer)
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

		{
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
			nestedRoot := nestedEntry{key: key}
			nestedRoot.populate(remainder)

			bytes := nestedRoot.extractBuffer(len(lookup.buffer), func(nestedKey tilePaletteKey, offset int) {
				toRemove = append(toRemove, nestedKey)
				lookup.starts[nestedKey] = offset
			})
			lookup.buffer = append(lookup.buffer, bytes...)
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
