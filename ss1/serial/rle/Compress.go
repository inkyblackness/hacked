package rle

import "io"

// Compress compresses the given byte array into the given writer.
// The optional reference array is used as a delta basis. If provided, bytes will be skipped
// where the data equals the reference.
func Compress(writer io.Writer, data []byte, reference []byte) {
	end := len(data)
	refLen := len(reference)
	ref := func(index int) byte {
		if index < refLen {
			return reference[index]
		}
		return 0x00
	}

	countIdenticalBytes := func(from int) int {
		index := from
		for (index < end) && (data[index] == ref(index)) {
			index++
		}
		return index - from
	}
	countConstantBytes := func(from int, value byte) int {
		index := from
		for (index < end) && (data[index] == value) {
			index++
		}
		return index - from
	}

	{
		trailingSkip := 0

		for temp := end - 1; (temp >= 0) && (data[temp] == ref(temp)); temp-- {
			trailingSkip++
		}
		end -= trailingSkip % 0x7FFF
	}
	start := 0
	for start < end {
		identicalByteCount := countIdenticalBytes(start)
		constByteCount := countConstantBytes(start, data[start])

		if identicalByteCount > 3 {
			skipBytes(writer, identicalByteCount)
			start += identicalByteCount
		} else if constByteCount > 3 {
			writeConstant(writer, constByteCount, data[start])
			start += constByteCount
		} else {
			diffByteCount := 0
			abort := false

			for (start+diffByteCount) < end && !abort {
				nextIdenticalByteCount := countIdenticalBytes(start + diffByteCount)
				nextConstByteCount := countConstantBytes(start+diffByteCount, data[start+diffByteCount])

				if nextIdenticalByteCount < 4 && nextConstByteCount < 4 {
					if nextIdenticalByteCount > nextConstByteCount {
						diffByteCount += nextIdenticalByteCount
					} else {
						diffByteCount += nextConstByteCount
					}
				} else {
					abort = true
				}
			}
			writeRaw(writer, data[start:start+diffByteCount])
			start += diffByteCount
		}
	}
	writeExtended(writer, 0x0000)
}

func writeExtended(writer io.Writer, control uint16, extra ...byte) {
	writer.Write([]byte{0x80, byte(control & 0xFF), byte((control >> 8) & 0xFF)})
	writer.Write(extra)
}

func skipBytes(writer io.Writer, size int) {
	remain := size

	for remain > 0 {
		if remain < 0x80 {
			writer.Write([]byte{byte(0x80 + remain)})
			remain = 0
		} else {
			lenControl := 0x7FFF
			if remain < lenControl {
				lenControl = remain
			}
			writeExtended(writer, uint16(lenControl))
			remain -= lenControl
		}
	}
}

func writeConstant(writer io.Writer, size int, value byte) {
	start := 0

	for start < size {
		remain := size - start
		if remain < 0x100 {
			writer.Write([]byte{0x00, byte(remain), value})
			start = size
		} else {
			lenControl := 0x3FFF
			if remain < lenControl {
				lenControl = remain
			}
			writeExtended(writer, 0xC000+uint16(lenControl), value)
			start += lenControl
		}
	}
}

func writeRaw(writer io.Writer, data []byte) {
	end := len(data)
	start := 0

	for start < end {
		remain := end - start
		if remain < 0x80 {
			writer.Write([]byte{byte(remain)})
			writer.Write(data[start:])
			start = end
		} else {
			lenControl := 0x3FFF

			if remain < lenControl {
				lenControl = remain
			}
			writeExtended(writer, 0x8000+uint16(lenControl), data[start:start+lenControl]...)
			start += lenControl
		}
	}
}
