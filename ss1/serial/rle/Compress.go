package rle

import "io"

// Compress compresses the given byte array into the given writer.
func Compress(writer io.Writer, data []byte) {
	end := len(data)
	start := 0

	countSameBytes := func(from int, value byte) int {
		index := from
		for index < end && data[index] == value {
			index++
		}
		return index - from
	}

	{
		trailingZeroes := 0

		for temp := end - 1; temp >= 0 && data[temp] == 0; temp-- {
			trailingZeroes++
		}
		end -= trailingZeroes % 0x7FFF
	}
	for start < end {
		startValue := data[start]
		sameByteCount := countSameBytes(start, startValue)

		if startValue == 0 {
			writeZero(writer, sameByteCount)
			start += sameByteCount
		} else if sameByteCount > 3 {
			writeConstant(writer, sameByteCount, startValue)
			start += sameByteCount
		} else {
			diffByteCount := 0
			abort := false

			for (start+diffByteCount) < end && !abort {
				startValue = data[start+diffByteCount]
				temp := countSameBytes(start+diffByteCount, startValue)

				if startValue != 0 && temp < 4 {
					diffByteCount += temp
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

func writeZero(writer io.Writer, size int) {
	remain := size

	for remain > 0 {
		if remain < 0x80 {
			writer.Write([]byte{byte(0x80 + remain)})
			remain = 0
		} else if remain < 0xFF {
			writer.Write([]byte{0xFF})
			remain -= 0x7F
		} else {
			lenControl := 0x7FFF

			if lenControl > remain {
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
