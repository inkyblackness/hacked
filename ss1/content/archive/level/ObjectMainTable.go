package level

import "github.com/inkyblackness/hacked/ss1/content/object"

const (
	// ObjectMainEntrySize describes the size, in bytes, of a ObjectMainEntry.
	ObjectMainEntrySize = 27

	defaultObjectMainEntryCount = 872
)

// ObjectMainEntryModifier is a function that changes properties.
type ObjectMainEntryModifier func(*ObjectMainEntry)

// ObjectMainEntry describes an object in the level.
type ObjectMainEntry struct {
	InUse byte

	Class    object.Class
	Subclass object.Subclass

	ClassTableIndex          int16
	CrossReferenceTableIndex int16
	Next                     ObjectID
	Prev                     ObjectID

	X         Coordinate
	Y         Coordinate
	Z         HeightUnit
	XRotation RotationUnit
	ZRotation RotationUnit
	YRotation RotationUnit

	_ byte

	Type object.Type

	Hitpoints int16

	Extra [4]byte
}

// Triple returns the unique identifier of the entry.
func (entry ObjectMainEntry) Triple() object.Triple {
	return object.TripleFrom(int(entry.Class), int(entry.Subclass), int(entry.Type))
}

// TilePosition returns the tile coordinates of the object.
func (entry ObjectMainEntry) TilePosition() TilePosition {
	return TilePosition{X: entry.X.Tile(), Y: entry.Y.Tile()}
}

// FinePosition returns the coordinates of the object within the tile.
func (entry ObjectMainEntry) FinePosition() FinePosition {
	return FinePosition{X: entry.X.Fine(), Y: entry.Y.Fine()}
}

// Reset clears the entry and resets all members.
func (entry *ObjectMainEntry) Reset() {
	*entry = ObjectMainEntry{}
}

// ObjectMainTable is a list of entries.
// The first entry is reserved for internal use. For the reserved entry,
// the Next field refers to the head of the single-linked free chain,
// and the CrossReferenceTableIndex refers to the head of the double-linked used chain.
type ObjectMainTable []ObjectMainEntry

// DefaultObjectMainTable returns an initialized table with a default size.
func DefaultObjectMainTable() ObjectMainTable {
	table := make(ObjectMainTable, defaultObjectMainEntryCount)
	table.Reset()
	return table
}

// Capacity returns the number of how many objects the table can hold.
func (table ObjectMainTable) Capacity() int {
	size := len(table)
	if size < 2 {
		return 0
	}
	return size - 1
}

// AllocatedCount returns the number of how many objects are currently allocated in the table.
func (table ObjectMainTable) AllocatedCount() int {
	capacity := table.Capacity()
	if capacity == 0 {
		return 0
	}
	index := int(table[0].CrossReferenceTableIndex)
	count := 0
	for (index != 0) && (count < capacity) {
		count++
		index = int(table[index].Next)
	}
	return count
}

// Reset wipes the entire table and initializes all links.
func (table ObjectMainTable) Reset() {
	tableLen := len(table)
	for i := 0; i < tableLen; i++ {
		entry := &table[i]
		entry.Reset()
		entry.Next = ObjectID(i + 1)
	}
	if tableLen > 0 {
		table[tableLen-1].Next = 0
	}
}

// Allocate attempts to find an available entry in the table and activates it.
// Returns 0 if not possible (exhausted).
func (table ObjectMainTable) Allocate() ObjectID {
	if len(table) < 2 {
		return 0
	}
	start := &table[0]
	if start.Next == 0 {
		return 0
	}
	id := start.Next
	entry := &table[id]
	start.Next = entry.Next

	entry.Reset()
	entry.Next = ObjectID(start.CrossReferenceTableIndex)
	table[start.CrossReferenceTableIndex].Prev = id
	entry.Prev = 0
	start.CrossReferenceTableIndex = int16(id)

	entry.InUse = 1

	return id
}

// Release deactivates the entry with given ID.
func (table ObjectMainTable) Release(id ObjectID) {
	if (id < 1) || (int(id) >= len(table)) {
		return
	}
	start := &table[0]
	entry := &table[id]

	if entry.InUse == 0 {
		return
	}

	if entry.Prev == 0 {
		start.CrossReferenceTableIndex = int16(entry.Next)
	} else {
		table[entry.Prev].Next = entry.Next
	}
	table[entry.Next].Prev = entry.Prev

	entry.Reset()
	entry.Next = start.Next
	start.Next = id
}
