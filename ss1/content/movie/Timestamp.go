package movie

const fractionDivisor = float32(0x10000)

// Timestamp represents a point in time using fixed resolution.
type Timestamp struct {
	Second   uint8
	Fraction uint16
}

// TimestampFromSeconds creates a timestamp instance from given floating point value.
func TimestampFromSeconds(value float32) Timestamp {
	second := uint8(value)
	return Timestamp{
		Second:   second,
		Fraction: uint16((value - float32(second)) * fractionDivisor),
	}
}
