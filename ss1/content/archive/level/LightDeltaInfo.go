package level

// LightDeltaInfo describes the light offsets for a tile.
type LightDeltaInfo byte

// OfCeiling returns the delta value for the ceiling.
func (info LightDeltaInfo) OfCeiling() int {
	return int(info >> 4)
}

// WithCeiling returns a new instance with the given ceiling value.
func (info LightDeltaInfo) WithCeiling(value int) LightDeltaInfo {
	if (value < 0) || (value > 15) {
		return info
	}
	return (info & 0x0F) | LightDeltaInfo(value<<4)
}

// OfFloor returns the delta value for the floor.
func (info LightDeltaInfo) OfFloor() int {
	return int(info & 0x0F)
}

// WithFloor returns a new instance with the given floor value.
func (info LightDeltaInfo) WithFloor(value int) LightDeltaInfo {
	if (value < 0) || (value > 15) {
		return info
	}
	return (info & 0xF0) | LightDeltaInfo(value)
}
