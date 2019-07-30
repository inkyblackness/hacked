package compression

import (
	"math/bits"
	"sync"
)

type paletteLookupEntry struct {
	start int
	size  int
}

// PaletteLookup is a dictionary of tile delta data, found in a palette buffer.
type PaletteLookup struct {
	buffer  []byte
	entries map[tilePaletteKey]paletteLookupEntry
}

// Buffer returns the underlying slice.
func (lookup *PaletteLookup) Buffer() []byte {
	return lookup.buffer
}

// Lookup finds the given tile again and returns the properties where and how to reproduce it.
func (lookup *PaletteLookup) Lookup(tile tileDelta) (index int, pal []byte, mask uint64) {
	key := tilePaletteKeyFrom(tile[:])
	entry, inLookup := lookup.entries[key]
	if inLookup {
		index = entry.start
		pal = lookup.buffer[entry.start : entry.start+entry.size]
	} else {
		pal = key.buffer()
	}
	var mapped [256]byte
	for mappedIndex, b := range pal {
		mapped[b] = byte(mappedIndex)
	}
	bitSize := uint(bits.Len(uint(len(pal) - 1)))
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

type nestedEntryCache struct {
	m       sync.Mutex
	keys    map[tilePaletteKey]struct{}
	entries map[tilePaletteKey]*nestedEntry
}

func (cache *nestedEntryCache) query(key tilePaletteKey, found func(*nestedEntry)) {
	go func() {
		cache.m.Lock()
		entry, hasEntry := cache.entries[key]
		cache.m.Unlock()
		if hasEntry {
			found(entry)
			return
		}

		entry = &nestedEntry{key: key}
		entry.populate(cache)
		cache.m.Lock()
		cache.entries[key] = entry
		cache.m.Unlock()
		found(entry)
	}()
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

func (entry *nestedEntry) populate(cache *nestedEntryCache) {
	remainingKey := entry.key
	foundSomething := true
	keySearchSize := remainingKey.size - 1
	for remainingKey.size > 2 && foundSomething {
		var lastAddedKey tilePaletteKey
		lastAddedKey, foundSomething = entry.populateRemaining(cache, remainingKey, keySearchSize)
		remainingKey = remainingKey.without(&lastAddedKey)
		keySearchSize = remainingKey.size
	}
}

func (entry *nestedEntry) populateRemaining(cache *nestedEntryCache,
	remainingKey tilePaletteKey, startSize int) (tilePaletteKey, bool) {
	maxByteSize := 0
	var maxNested *nestedEntry
	keySize := startSize
	for (keySize > 2) && (maxNested == nil) {
		results := make(chan *nestedEntry)
		resultsPending := 0
		for otherKey := range cache.keys {
			if otherKey.size == keySize && remainingKey.contains(&otherKey) {
				resultsPending++
				cache.query(otherKey, func(nested *nestedEntry) {
					results <- nested
				})
			}
		}
		for resultsPending > 0 {
			nested := <-results
			nestedSize := nested.byteSize()
			if nestedSize > maxByteSize {
				maxByteSize = nestedSize
				maxNested = nested
			}
			resultsPending--
		}
		close(results)
		keySize--
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
	lookup.entries = make(map[tilePaletteKey]paletteLookupEntry)

	remainder := make(map[tilePaletteKey]struct{})
	for key := range gen.keyUses {
		remainder[key] = struct{}{}
	}

	for size := PixelPerTile; size > 2; size-- {

		{
			var earlyRemoved []tilePaletteKey
			for key := range remainder {
				wasRemoved := false
				for _, fitSize := range []int{4, 8, 16} {
					if key.size <= fitSize {
						for start := 0; start < (len(lookup.buffer)-fitSize) && !wasRemoved; start++ {
							tempKey := tilePaletteKeyFrom(lookup.buffer[start : start+fitSize])
							if tempKey.contains(&key) {
								earlyRemoved = append(earlyRemoved, key)
								wasRemoved = true

								lookup.entries[key] = paletteLookupEntry{start: start, size: fitSize}
							}
						}
					}
				}
			}
			for _, key := range earlyRemoved {
				delete(remainder, key)
			}
		}

		var keysInSize []tilePaletteKey
		// find all keys with this current size
		for key := range remainder {
			if key.size == size {
				keysInSize = append(keysInSize, key)
			}
		}

		for _, key := range keysInSize {
			var toRemove []tilePaletteKey
			cache := nestedEntryCache{
				keys:    remainder,
				entries: make(map[tilePaletteKey]*nestedEntry),
			}
			nestedRoot := nestedEntry{key: key}
			nestedRoot.populate(&cache)

			bytes := nestedRoot.extractBuffer(len(lookup.buffer), func(nestedKey tilePaletteKey, offset int) {
				toRemove = append(toRemove, nestedKey)
				lookup.entries[nestedKey] = paletteLookupEntry{start: offset, size: nestedKey.size}
			})
			lookup.buffer = append(lookup.buffer, bytes...)
			for _, key := range toRemove {
				delete(remainder, key)
			}
		}
	}

	for key := range remainder {
		bytes := key.buffer()
		lookup.entries[key] = paletteLookupEntry{start: len(lookup.buffer), size: len(bytes)}
		lookup.buffer = append(lookup.buffer, bytes...)
	}

	return lookup
}

// Add registers a further delta to the generator.
func (gen *PaletteLookupGenerator) Add(delta tileDelta) {
	key := tilePaletteKeyFrom(delta[:])
	if key.size > 2 {
		if gen.keyUses == nil {
			gen.keyUses = make(map[tilePaletteKey]int)
		}
		gen.keyUses[key]++
	}
}
