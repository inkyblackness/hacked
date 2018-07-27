package compression

// BitstreamReader is a utility to read big-endian integer values of arbitrary bit size from a bitstream.
type BitstreamReader struct {
	source          []byte
	nextSourceIndex int
	sourceBitLen    uint64

	currentBitPos uint64

	buffer       uint64
	bitsBuffered uint64
}

// NewBitstreamReader returns a new instance of a bistream reader using the provided byte slice as source.
func NewBitstreamReader(data []byte) *BitstreamReader {
	return &BitstreamReader{
		source:       data,
		sourceBitLen: uint64(len(data)) * 8}
}

func (reader *BitstreamReader) bufferNextByte() {
	reader.buffer = (reader.buffer << 8) | uint64(reader.source[reader.nextSourceIndex])
	reader.bitsBuffered += 8
	reader.nextSourceIndex++
}

// Exhausted returns true if the current position is at or beyond the possible length.
// Reading an exhausted stream yields only zeroes.
func (reader *BitstreamReader) Exhausted() bool {
	return reader.currentBitPos >= reader.sourceBitLen
}

// Read returns a value with the requested bit size, right aligned, as a uint32.
// Reading does not advance the current position. A successful read of a certain size will return the same
// value when called repeatedly with the same parameter.
// Reading is possible beyond the available size. For the missing bits, a value of 0 is provided.
//
// The function panics when reading more than 32 bits.
func (reader *BitstreamReader) Read(bits int) (result uint32) {
	if bits > 32 {
		panic("Limit of bit count: 32")
	}

	available := reader.sourceBitLen - reader.currentBitPos
	toRead := uint64(bits)
	if toRead > available {
		toRead = available
	}

	for reader.bitsBuffered < uint64(toRead) {
		reader.bufferNextByte()
	}
	result = uint32(reader.buffer >> (reader.bitsBuffered - uint64(toRead)) & ^(uint64(0xFFFFFFFFFFFFFFFF) << uint64(toRead)))
	result <<= uint64(bits) - toRead

	return
}

// Advance skips the provided amount of bits and puts the current position there.
// Successive read operations will return values from the new position.
// It is possible to advance beyond the available length.
//
// This function panics when advancing with negative values.
func (reader *BitstreamReader) Advance(bits int) {
	if bits < 0 {
		panic("Can only advance forward")
	}

	newIndex := reader.currentBitPos + uint64(bits)
	if newIndex > reader.sourceBitLen {
		newIndex = reader.sourceBitLen
	}
	reader.currentBitPos = newIndex
	if uint64(bits) > reader.bitsBuffered {
		reader.bitsBuffered = 0
		reader.nextSourceIndex = int(reader.currentBitPos / 8)
		remainder := reader.currentBitPos % 8
		if remainder > 0 {
			reader.bufferNextByte()
			reader.bitsBuffered -= remainder
		}
	} else {
		reader.bitsBuffered -= uint64(bits)
	}
}
