package compression

// MaskstreamReader reads mask integers from a byte array.
type MaskstreamReader struct {
	source []byte
	curPos int
}

// NewMaskstreamReader returns a new reader instance for given source.
func NewMaskstreamReader(source []byte) *MaskstreamReader {
	return &MaskstreamReader{source: source}
}

func (reader *MaskstreamReader) nextByte() (value uint64) {
	if reader.curPos < len(reader.source) {
		value = uint64(reader.source[reader.curPos])
		reader.curPos++
	}

	return
}

// Read returns a mask integer of given byte length from the current position.
// If the current position is at or beyond the available size, 0x00 is assumed for the missing bytes.
//
// Reading more than 8, or less than 0, bytes panics.
func (reader *MaskstreamReader) Read(bytes int) (value uint64) {
	if bytes > 8 {
		panic("Limit of byte count: 8")
	}
	if bytes < 0 {
		panic("Minimum byte count: 0")
	}

	for i := 0; i < bytes; i++ {
		value |= reader.nextByte() << uint64(8*i)
	}

	return
}
