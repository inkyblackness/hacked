package level

import (
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/serial"
)

const (
	// ObjectClassEntryHeaderSize is the size, in bytes, of the header prefix for each class entry.
	ObjectClassEntryHeaderSize = 6
)

// ObjectClassEntry describes an entry in a object-class specific list.
type ObjectClassEntry struct {
	ObjectID ObjectID
	Next     int16
	Prev     int16
	Data     []byte
}

// NewObjectClassEntry returns an instance with the given amount of bytes initialized for data.
func NewObjectClassEntry(dataSize int) *ObjectClassEntry {
	return &ObjectClassEntry{
		Data: make([]byte, dataSize),
	}
}

// Code serializes the entry with given coder, including the data array.
func (entry *ObjectClassEntry) Code(coder serial.Coder) {
	coder.Code(&entry.ObjectID)
	coder.Code(&entry.Next)
	coder.Code(&entry.Prev)
	coder.Code(entry.Data)
}

// Reset sets all members of the entry to zero, including all bytes of the data array.
func (entry *ObjectClassEntry) Reset() {
	entry.ObjectID = 0
	entry.Next = 0
	entry.Prev = 0
	for i := 0; i < len(entry.Data); i++ {
		entry.Data[i] = 0
	}
}

// ObjectClassTable is a list of entries.
//
// The first entry is reserved for internal use. For the reserved entry,
// the Next field identifies the head of the single-linked free chain,
// and the ObjectID field identifies the head of the double-linked used chain.
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
	tableLen := len(table)
	for i := 0; i < tableLen; i++ {
		entry := &table[i]
		entry.Reset()
		entry.Next = int16(i + 1)
	}
	if tableLen > 0 {
		table[tableLen-1].Next = 0
	}
}

// Code serializes the table with the provided coder.
func (table ObjectClassTable) Code(coder serial.Coder) {
	for i := 0; i < len(table); i++ {
		table[i].Code(coder)
	}
}

// Allocate attempts to reserve an entry in the table and returns the corresponding index.
// returns 0 if exhausted.
func (table ObjectClassTable) Allocate() int {
	if len(table) < 2 {
		return 0
	}
	start := &table[0]
	if start.Next == 0 {
		return 0
	}
	index := start.Next
	entry := &table[index]
	start.Next = entry.Next

	entry.Reset()
	entry.Next = int16(start.ObjectID)
	table[start.ObjectID].Prev = index
	entry.Prev = 0
	start.ObjectID = ObjectID(index)

	return int(index)
}

// Release frees the identified entry.
func (table ObjectClassTable) Release(index int) {
	if (index < 1) || (index >= len(table)) {
		return
	}
	start := &table[0]
	entry := &table[index]

	if entry.Prev == 0 {
		start.ObjectID = ObjectID(entry.Next)
	} else {
		table[entry.Prev].Next = entry.Next
	}
	table[entry.Next].Prev = entry.Prev

	entry.Reset()
	entry.Next = start.Next
	start.Next = int16(index)
}
