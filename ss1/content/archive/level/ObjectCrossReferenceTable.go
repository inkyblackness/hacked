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
