package compression

// BitstreamWriter is a utility to write big-endian integer values of arbitrary bit size to a bitstream.
type BitstreamWriter struct {
	buf []byte

	offset uint
}

// Buffer returns the currently stored buffer.
func (w BitstreamWriter) Buffer() []byte {
	return w.buf
}

// Write stores the given amount of bits from the provided value.
// Trying to write more than 32 bits will cause panic.
func (w *BitstreamWriter) Write(bits uint, value uint32) {
	if bits > 32 {
		panic("maximum of 32 bits are possible to be written")
	}
	for bits > 0 {
		temp := (uint64(value) & ^(^uint64(0) << bits)) << (40 - bits - w.offset)
		if w.offset == 0 {
			w.buf = append(w.buf, 0x00)
		}
		w.buf[len(w.buf)-1] |= byte(temp >> 32)
		written := 8 - w.offset
		if written > bits {
			written = bits
		}
		bits -= written
		w.offset = (w.offset + written) % 8
	}
}
