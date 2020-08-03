package level

import "github.com/inkyblackness/hacked/ss1/content/object"

const (
	// ObjectMainEntrySize describes the size, in bytes, of a ObjectMasterEntry.
	ObjectMainEntrySize = 27

	defaultObjectMasterEntryCount = 872
)

// ObjectMasterEntry describes an object in the level.
type ObjectMasterEntry struct {
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
func (entry ObjectMasterEntry) Triple() object.Triple {
	return object.TripleFrom(int(entry.Class), int(entry.Subclass), int(entry.Type))
}

// Reset clears the entry and resets all members.
func (entry *ObjectMasterEntry) Reset() {
	*entry = ObjectMasterEntry{}
}

// ObjectMasterTable is a list of entries.
// The first entry is reserved for internal use. For the reserved entry,
// the Next field refers to the head of the single-linked free chain,
// and the CrossReferenceTableIndex refers to the head of the double-linked used chain.
type ObjectMasterTable []ObjectMasterEntry

// DefaultObjectMasterTable returns an initialized table with a default size.
func DefaultObjectMasterTable() ObjectMasterTable {
	table := make(ObjectMasterTable, defaultObjectMasterEntryCount)
	table.Reset()
	return table
}

// Reset wipes the entire table and initializes all links.
func (table ObjectMasterTable) Reset() {
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
func (table ObjectMasterTable) Allocate() ObjectID {
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
func (table ObjectMasterTable) Release(id ObjectID) {
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
