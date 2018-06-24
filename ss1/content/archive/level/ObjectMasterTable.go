package level

import "github.com/inkyblackness/hacked/ss1/content/object"

const (
	// MasterObjectEntrySize describes the size, in bytes, of a MasterObjectEntry.
	MasterObjectEntrySize = 27

	defaultMasterObjectEntryCount = 872
)

// MasterObjectEntry describes an object in the level.
type MasterObjectEntry struct {
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

// MasterObjectTable is a list of entries.
type MasterObjectTable []MasterObjectEntry

// DefaultMasterObjectTable returns an initialized table with a default size.
func DefaultMasterObjectTable() MasterObjectTable {
	table := make(MasterObjectTable, defaultMasterObjectEntryCount)
	table.Reset()
	return table
}

// Reset wipes the entire table and initializes all links.
func (table MasterObjectTable) Reset() {
	// TODO
}
