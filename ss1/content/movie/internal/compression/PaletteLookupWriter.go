package compression

import "bytes"

// PaletteLookupWriter writes palette lookup tables.
type PaletteLookupWriter struct {
	Buffer []byte
}

// Write ensures the given sequence of bytes is found in the buffer.
// The returned value is the offset into the buffer.
// It may add the given bytes at the end, or find the existing sequence within the buffer.
func (w *PaletteLookupWriter) Write(data []byte) uint32 {
	byteCount := len(data)
	for offset := 0; offset <= len(w.Buffer)-byteCount; offset++ {
		if bytes.Equal(data, w.Buffer[offset:offset+byteCount]) {
			return uint32(offset)
		}
	}
	for remaining := 1; remaining < byteCount; remaining++ {
		existing := byteCount - remaining
		if existing <= len(w.Buffer) {
			if bytes.Equal(data[:existing], w.Buffer[len(w.Buffer)-existing:]) {
				w.Buffer = append(w.Buffer, data[existing:]...)
				return uint32(len(w.Buffer) - byteCount)
			}
		}
	}
	w.Buffer = append(w.Buffer, data...)
	return uint32(len(w.Buffer) - byteCount)
}
