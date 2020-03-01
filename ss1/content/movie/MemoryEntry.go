package movie

import (
	"bytes"
	"encoding/binary"
)

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

// NewSubtitleEntry creates a new entry for given subtitle information.
func NewSubtitleEntry(timestamp Timestamp, control SubtitleControl, text []byte) Entry {
	buf := bytes.NewBuffer(nil)
	var subtitleHeader SubtitleHeader
	subtitleHeader.Control = control
	subtitleHeader.TextOffset = SubtitleDefaultTextOffset
	_ = binary.Write(buf, binary.LittleEndian, &subtitleHeader)
	buf.Write(make([]byte, SubtitleDefaultTextOffset-buf.Len()))
	buf.Write(text)

	return NewMemoryEntry(timestamp, Subtitle, buf.Bytes())
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
