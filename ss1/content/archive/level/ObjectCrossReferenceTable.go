package level

const (
	// ObjectCrossReferenceEntrySize describes the size, in bytes, of a ObjectCrossReferenceEntry.
	ObjectCrossReferenceEntrySize = 10

	defaultObjectCrossReferenceEntryCount = 1600
)

// ObjectCrossReferenceEntry links objects and tiles.
type ObjectCrossReferenceEntry struct {
	TileX int16
	TileY int16

	ObjectID       ObjectID
	NextInTile     int16
	NextTileForObj int16
}

// Reset clears the members of the entry.
func (entry *ObjectCrossReferenceEntry) Reset() {
	*entry = ObjectCrossReferenceEntry{}
}

// ObjectCrossReferenceTable is a list of entries.
// The first entry is reserved for internal use. For the reserved entry,
// The NextInTile member refers to the head of the single-linked free chain.
type ObjectCrossReferenceTable []ObjectCrossReferenceEntry

// DefaultObjectCrossReferenceTable returns an initialized table with a default size.
func DefaultObjectCrossReferenceTable() ObjectCrossReferenceTable {
	table := make(ObjectCrossReferenceTable, defaultObjectCrossReferenceEntryCount)
	table.Reset()
	return table
}

// Reset wipes the entire table and initializes all links.
func (table ObjectCrossReferenceTable) Reset() {
	tableLen := len(table)
	for i := 0; i < tableLen; i++ {
		entry := &table[i]
		entry.Reset()
		entry.NextInTile = int16(i + 1)
	}
	if tableLen > 0 {
		table[tableLen-1].NextInTile = 0
	}
}

// Allocate attempts to reserve a free entry in the table and return its index.
// Returns 0 if exhausted.
func (table ObjectCrossReferenceTable) Allocate() int {
	if len(table) < 2 {
		return 0
	}
	start := &table[0]
	if start.NextInTile == 0 {
		return 0
	}
	index := start.NextInTile
	entry := &table[index]
	start.NextInTile = entry.NextInTile

	entry.Reset()

	return int(index)
}

// Release frees the entry with given index.
func (table ObjectCrossReferenceTable) Release(index int) {
	if (index < 1) || (index >= len(table)) {
		return
	}
	start := &table[0]
	entry := &table[index]

	entry.Reset()
	entry.NextInTile = start.NextInTile
	start.NextInTile = int16(index)
}
