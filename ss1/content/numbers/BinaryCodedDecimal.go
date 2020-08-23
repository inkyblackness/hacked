package numbers

// ToBinaryCodedDecimal converts given integer value to a 4-digit BCD value.
func ToBinaryCodedDecimal(value uint16) (bcd uint16) {
	factors := []uint16{1000, 100, 10, 1}
	value %= 10000

	for _, factor := range factors {
		bcd = (bcd << 4) | (value / factor)
		value %= factor
	}

	return
}

// FromBinaryCodedDecimal converts given 4-digit BCD value into an integer.
func FromBinaryCodedDecimal(bcd uint16) (value uint16) {
	factors := []uint16{1, 10, 100, 1000}

	for _, factor := range factors {
		value += (bcd & 0xF) * factor
		bcd >>= 4
	}

	return
}
