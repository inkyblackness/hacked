package level

import "github.com/inkyblackness/hacked/ss1/content/object"

const (
	// ObjectMasterEntrySize describes the size, in bytes, of a ObjectMasterEntry.
	ObjectMasterEntrySize = 27

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

// Reset clears the entry and resets all members.
func (entry *ObjectMasterEntry) Reset() {
	*entry = ObjectMasterEntry{}
}

// ObjectMasterTable is a list of entries.
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
