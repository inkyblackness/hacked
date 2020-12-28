package compression

import (
	"context"
	"math/bits"
	"sort"
)

type paletteLookupEntry struct {
	start int
	size  int
}

// PaletteLookup is a dictionary of tile delta data, found in a palette buffer.
type PaletteLookup struct {
	buffer  []byte
	entries map[TilePaletteKey]paletteLookupEntry
}

// Buffer returns the underlying slice.
func (lookup *PaletteLookup) Buffer() []byte {
	return lookup.buffer
}

// Lookup finds the given tile again and returns the properties where and how to reproduce it.
func (lookup *PaletteLookup) Lookup(tile tileDelta) (index int, pal []byte, mask uint64) {
	key := TilePaletteKeyFrom(tile[:])
	entry, inLookup := lookup.entries[key]
	if inLookup {
		index = entry.start
		pal = lookup.buffer[entry.start : entry.start+entry.size]
	} else {
		pal = key.Buffer()
	}
	var mapped [256]byte
	for mappedIndex := len(pal) - 1; mappedIndex >= 0; mappedIndex-- {
		mapped[pal[mappedIndex]] = byte(mappedIndex)
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
	keyUses map[TilePaletteKey]int
}

// Generate creates a lookup based on all currently registered tile deltas.
func (gen *PaletteLookupGenerator) Generate(ctx context.Context) (PaletteLookup, error) {
	var lookup PaletteLookup
	lookup.entries = make(map[TilePaletteKey]paletteLookupEntry)

	remainder := make(map[TilePaletteKey]struct{})
	for key := range gen.keyUses {
		remainder[key] = struct{}{}
	}

	type sizedEntry struct {
		entries    map[TilePaletteKey]paletteLookupEntry
		lastOffset int
	}
	sizedEntries := make(map[int]*sizedEntry)
	knownSizes := []int{4, 8, 16}
	for _, size := range knownSizes {
		sizedEntries[size] = &sizedEntry{
			entries: make(map[TilePaletteKey]paletteLookupEntry),
		}
	}

	toDeleteCache := make([]TilePaletteKey, 0, len(remainder))
	addToBuffer := func(data []byte) {
		lookup.buffer = append(lookup.buffer, data...)

		newSize := len(lookup.buffer)

		for _, fitSize := range knownSizes {
			fitLimit := newSize - fitSize
			entry := sizedEntries[fitSize]

			// remove all entries beyond a certain limit. as these bytes don't change, retrying won't help.
			toDelete := toDeleteCache[0:0]
			limit := newSize - 16 - len(data)
			for key, entry := range entry.entries {
				if entry.start < limit {
					toDelete = append(toDelete, key)
				}
			}
			for _, key := range toDelete {
				delete(entry.entries, key)
			}

			// find any new keys
			for start := entry.lastOffset; start < fitLimit; start++ {
				tempKey := TilePaletteKeyFrom(lookup.buffer[start : start+fitSize])
				if _, existing := entry.entries[tempKey]; !existing {
					entry.entries[tempKey] = paletteLookupEntry{
						start: start,
						size:  fitSize,
					}
				}
			}
			if fitLimit > 0 {
				entry.lastOffset = fitLimit
			}
		}
	}

	addEarlyEntry := func(key TilePaletteKey, limitSize int) bool {
		for _, fitSize := range knownSizes {
			if key.size <= fitSize && fitSize <= limitSize {
				entry := sizedEntries[fitSize]
				for tempKey, paletteEntry := range entry.entries {
					if tempKey.Contains(&key) && (!key.HasColor(0x00) || (lookup.buffer[paletteEntry.start] == 0x00)) {
						lookup.entries[key] = paletteEntry
						return true
					}
				}
			}
		}
		return false
	}

	sizeLimitForSize := map[int]int{3: 4, 4: 8, 5: 8, 6: 8, 7: 8, 8: 8, 9: 16, 10: 16, 11: 16, 12: 16, 13: 16, 14: 16, 15: 16, 16: 16}
	for size := PixelPerTile; (size > 2) && (ctx.Err() == nil); size-- {
		keysInSize := make([]TilePaletteKey, 0, len(remainder))
		for key := range remainder {
			if key.size == size {
				keysInSize = append(keysInSize, key)
			}
		}
		sort.Slice(keysInSize, func(a, b int) bool { return keysInSize[a].LessThan(&keysInSize[b]) })

		// fmt.Printf("Working on key size %v, have %v sized, %v total remaining\n", size, len(keysInSize), len(remainder))
		for _, sizedKey := range keysInSize {
			if ctx.Err() != nil {
				return PaletteLookup{}, ctx.Err()
			}

			{
				var earlyRemoved []TilePaletteKey
				for key := range remainder {
					if addEarlyEntry(key, sizeLimitForSize[size]) {
						earlyRemoved = append(earlyRemoved, key)
					}
				}
				for _, key := range earlyRemoved {
					delete(remainder, key)
				}
			}
			if _, stillRemaining := remainder[sizedKey]; stillRemaining {
				bytes := sizedKey.Buffer()
				lookup.entries[sizedKey] = paletteLookupEntry{start: len(lookup.buffer), size: len(bytes)}
				addToBuffer(bytes)

				delete(remainder, sizedKey)
			}
		}
	}

	return lookup, ctx.Err()
}

// Add registers a further delta to the generator.
func (gen *PaletteLookupGenerator) Add(delta tileDelta) {
	key := TilePaletteKeyFrom(delta[:])
	if key.size > 2 {
		if gen.keyUses == nil {
			gen.keyUses = make(map[TilePaletteKey]int)
		}
		gen.keyUses[key]++
	}
}
