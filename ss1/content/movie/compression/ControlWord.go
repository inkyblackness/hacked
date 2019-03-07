package compression

// ControlWord describes the current compression action.
type ControlWord uint32

// ControlWordOf returns a word with the given paramters.
func ControlWordOf(count int, controlType ControlType, param uint32) ControlWord {
	return (ControlWord(count) << 20) | (ControlWord(controlType&0x7) << 17) | ControlWord(param&0x1FFFF)
}

// Packed returns the control word packed in a PackedControlWord with the given number of times.
func (word ControlWord) Packed(times int) PackedControlWord {
	return PackedControlWord((uint32(word) & 0x00FFFFFF) | (uint32(times) << 24))
}

// Count returns the count value of the control
func (word ControlWord) Count() int {
	return int((uint32(word) >> 20) & 0xF)
}

// IsLongOffset returns true if Count() returns 0.
func (word ControlWord) IsLongOffset() bool {
	return word.Count() == 0
}

// LongOffset returns the long offset value. Only relevant if IsLongOffset() returns true.
func (word ControlWord) LongOffset() uint32 {
	return (uint32(word) >> 0) & 0xFFFFF
}

// Type returns the type of the control. Only relevant if IsLongOffset() returns false.
func (word ControlWord) Type() ControlType {
	return ControlType((uint32(word) >> 17) & 0x7)
}

// Parameter returns the type specific parameter value.
func (word ControlWord) Parameter() uint32 {
	return (uint32(word) >> 0) & 0x1FFFF
}
