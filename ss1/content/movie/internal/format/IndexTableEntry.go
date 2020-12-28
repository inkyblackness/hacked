package format

// IndexTableEntrySize specifies how many bytes a serialized IndexTableEntry uses.
const IndexTableEntrySize = 8

const indexHeaderSizeIncrement = 0x0400

// IndexTableEntry describes one data entry of a MOVI container.
type IndexTableEntry struct {
	TimestampFraction uint16
	TimestampSecond   byte
	Type              byte
	DataOffset        int32
}

// IndexTableSizeFor returns the number of bytes reserved for given amount of table entries.
func IndexTableSizeFor(entryCount int) int {
	size := indexHeaderSizeIncrement
	requiredSize := entryCount * IndexTableEntrySize

	if requiredSize > size {
		size *= requiredSize/size + 2
	}

	return size
}
