package compression

type dictEntry struct {
	prev  *dictEntry
	depth int

	next [256]*dictEntry

	key   word
	value byte
	first byte
	used  bool
}

func rootDictEntry() *dictEntry {
	return &dictEntry{prev: nil, depth: 0, value: 0x00, key: reset}
}

var emptyList [256]*dictEntry

func (entry *dictEntry) Add(value byte, key word, newEntry *dictEntry) *dictEntry {
	newEntry.prev = entry
	newEntry.depth = entry.depth + 1
	newEntry.value = value
	newEntry.key = key
	newEntry.first = entry.first
	if newEntry.used {
		copy(newEntry.next[:], emptyList[:])
	}
	newEntry.used = true

	entry.next[value] = newEntry
	if entry.depth == 0 {
		newEntry.first = value
	}

	return newEntry
}

func (entry *dictEntry) Data(bytes []byte) {
	cur := entry
	for i := entry.depth - 1; i >= 0; i-- {
		bytes[i] = cur.value
		cur = cur.prev
	}
}

func (entry *dictEntry) FirstByte() byte {
	return entry.first
}
