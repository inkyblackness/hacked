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

	MasterObjectTableIndex int16
	Next                   int16
	NextTile               int16
}

// ObjectCrossReferenceTable is a list of entries.
type ObjectCrossReferenceTable []ObjectCrossReferenceEntry

// DefaultObjectCrossReferenceTable returns an initialized table with a default size.
func DefaultObjectCrossReferenceTable() ObjectCrossReferenceTable {
	table := make(ObjectCrossReferenceTable, defaultObjectCrossReferenceEntryCount)
	table.Reset()
	return table
}

// Reset wipes the entire table and initializes all links.
func (table ObjectCrossReferenceTable) Reset() {
	// TODO
}
