package compression

// PackedControlWord contains a ControlWord together with a count value of occurrences.
type PackedControlWord uint32

// Times returns the amount of occurrences of the contained control word.
func (packed PackedControlWord) Times() int {
	return int((packed >> 24) & 0xFF)
}

// Value returns the control word value.
func (packed PackedControlWord) Value() ControlWord {
	return ControlWord(uint32(packed) & 0x00FFFFFF)
}
