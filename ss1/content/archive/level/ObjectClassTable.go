package level

import "github.com/inkyblackness/hacked/ss1/content/object"

// ObjectClassEntry describes an entry in a object-class specific list.
type ObjectClassEntry struct {
	ObjectID ObjectID
	Next     int16
	Prev     int16
	Data     []byte
}

// ObjectClassTable is a list of entries.
type ObjectClassTable []ObjectClassEntry

// DefaultObjectClassTable returns an initialized table for given object class.
func DefaultObjectClassTable(class object.Class) ObjectClassTable {
	info := ObjectClassInfoFor(class)
	table := make(ObjectClassTable, info.EntryCount)
	table.AllocateData(info.DataSize)
	table.Reset()
	return table
}

// AllocateData prepares each entry to be able to store the given amount of bytes.
// This function drops any previously assigned data.
func (table ObjectClassTable) AllocateData(size int) {
	// TODO
}

// Reset wipes the entire table and initializes all links.
// The bytes of the data members of each entry are reset to zero.
func (table ObjectClassTable) Reset() {
	// TODO
}
