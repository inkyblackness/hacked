package rle

import (
	"errors"
	"io"
)

// Decompress decompresses from the given reader and writes into the provided output buffer.
// The output buffer must be pre-allocated.
// If it contains non-zero data, this data may be preserved if the compressed data specifies to.
func Decompress(reader io.Reader, output []byte) (err error) {
	outIndex := 0
	done := false
	nextByte := func() byte {
		zz := []byte{0x00}
		_, err = reader.Read(zz)
		return zz[0]
	}

	for !done && (err == nil) {
		first := nextByte()

		switch {
		case first == 0x00:
			nn := nextByte()
			zz := nextByte()

			outIndex += writeBytesOfValue(output[outIndex:outIndex+int(nn)], func() byte { return zz })
		case first < 0x80:
			outIndex += writeBytesOfValue(output[outIndex:outIndex+int(first)], nextByte)
		case first == 0x80:
			control := uint16(nextByte())
			control += uint16(nextByte()) << 8
			switch {
			case control == 0x0000:
				done = true
			case control < 0x8000:
				outIndex += int(control)
			case control < 0xC000:
				outIndex += writeBytesOfValue(output[outIndex:outIndex+int(control&0x3FFF)], nextByte)
			case (control & 0xFF00) == 0xC000:
				err = errors.New("undefined case 80 nn C0")
			default:
				zz := nextByte()

				outIndex += writeBytesOfValue(output[outIndex:outIndex+int(control&0x3FFF)], func() byte { return zz })
			}
		default:
			outIndex += int(first & 0x7F)
		}
	}

	return
}

func writeBytesOfValue(buffer []byte, producer func() byte) int {
	count := len(buffer)
	for i := 0; i < count; i++ {
		buffer[i] = producer()
	}
	return count
}
