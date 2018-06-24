package level

import (
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/serial"
)

// ObjectClassEntry describes an entry in a object-class specific list.
type ObjectClassEntry struct {
	ObjectID ObjectID
	Next     int16
	Prev     int16
	Data     []byte
}

// NewObjectClassEntry returns an instance with the given amount of bytes initialized for data.
func NewObjectClassEntry(dataSize int) ObjectClassEntry {
	return ObjectClassEntry{
		Data: make([]byte, dataSize),
	}
}

// Code serializes the entry with given coder, including the data array.
func (entry ObjectClassEntry) Code(coder serial.Coder) {
	coder.Code(&entry.ObjectID)
	coder.Code(&entry.Next)
	coder.Code(&entry.Prev)
	coder.Code(entry.Data)
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
	for i := 0; i < len(table); i++ {
		table[i].Data = make([]byte, size)
	}
}

// Reset wipes the entire table and initializes all links.
// The bytes of the data members of each entry are reset to zero.
func (table ObjectClassTable) Reset() {
	// TODO
}

// Code serializes the table with the provided coder.
func (table ObjectClassTable) Code(coder serial.Coder) {
	for i := 0; i < len(table); i++ {
		coder.Code(&table[i])
	}
}
