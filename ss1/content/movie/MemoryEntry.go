package movie

type memoryEntry struct {
	timestamp float32
	dataType  DataType
	data      []byte
}

// NewMemoryEntry returns an Entry instance that has the properties in memory.
func NewMemoryEntry(timestamp float32, dataType DataType, data []byte) Entry {
	entry := &memoryEntry{
		timestamp: timestamp,
		dataType:  dataType,
		data:      data}

	return entry
}

func (entry *memoryEntry) Timestamp() float32 {
	return entry.timestamp
}

func (entry *memoryEntry) Type() DataType {
	return entry.dataType
}

func (entry *memoryEntry) Data() []byte {
	return entry.data
}
