package format

// IndexTableEntrySize specifies how many bytes a serialized IndexTableEntry uses.
const IndexTableEntrySize = 8

// IndexTableEntry describes one data entry of a MOVI container.
type IndexTableEntry struct {
	TimestampFraction uint16
	TimestampSecond   byte
	Type              byte
	DataOffset        int32
}
