package movie

type memoryEntry struct {
	timestamp Timestamp
	dataType  DataType
	data      []byte
}

// NewMemoryEntry returns an Entry instance that has the properties in memory.
func NewMemoryEntry(timestamp Timestamp, dataType DataType, data []byte) Entry {
	entry := &memoryEntry{
		timestamp: timestamp,
		dataType:  dataType,
		data:      data}

	return entry
}

func (entry *memoryEntry) Timestamp() Timestamp {
	return entry.timestamp
}

func (entry *memoryEntry) Type() DataType {
	return entry.dataType
}

func (entry *memoryEntry) Data() []byte {
	return entry.data
}
