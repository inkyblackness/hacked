package movie

import "time"

const fractionDivisor = 0x10000

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

// ToDuration returns the equivalent duration for this timestamp.
func (ts Timestamp) ToDuration() time.Duration {
	return time.Duration((float64(ts.Second) + float64(ts.Fraction)/fractionDivisor) * float64(time.Second))
}

// IsAfter returns true if this timestamp is later than the given one.
func (ts Timestamp) IsAfter(other Timestamp) bool {
	return (ts.Second > other.Second) || ((ts.Second == other.Second) && (ts.Fraction > other.Fraction))
}

// DeltaTo returns the absolute difference between this and the given timestamp.
func (ts Timestamp) DeltaTo(other Timestamp) Timestamp {
	tsLinear := ts.toLinear()
	otherLinear := other.toLinear()
	if tsLinear < otherLinear {
		return timestampFromLinear(otherLinear - tsLinear)
	}
	return timestampFromLinear(tsLinear - otherLinear)
}

// Plus returns a timestamp with the given one added to the current one.
// The result is saturated if the addition is larger than the timestamp can hold.
func (ts Timestamp) Plus(other Timestamp) Timestamp {
	delta := ts.toLinear() + other.toLinear()
	if delta > 0xFFFFFF {
		return timestampFromLinear(0xFFFFFF)
	}
	return timestampFromLinear(delta)
}

// Minus returns a timestamp with given one removed from the current one.
// The result is clipped to 0 if subtraction would result in a negative number.
func (ts Timestamp) Minus(other Timestamp) Timestamp {
	otherLinear := other.toLinear()
	tsLinear := ts.toLinear()
	if otherLinear >= tsLinear {
		return timestampFromLinear(0)
	}
	return timestampFromLinear(tsLinear - otherLinear)
}

func (ts Timestamp) toLinear() uint32 {
	return uint32(ts.Second)<<16 + uint32(ts.Fraction)
}

func timestampFromLinear(value uint32) Timestamp {
	return Timestamp{
		Second:   uint8(value >> 16),
		Fraction: uint16(value),
	}
}
