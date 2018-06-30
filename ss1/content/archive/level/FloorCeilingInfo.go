package level

// FloorInfo describes the properties of a floor
type FloorInfo byte

// AbsoluteHeight returns the floor height in range of [0..TileHeightUnitMax-1].
func (info FloorInfo) AbsoluteHeight() TileHeightUnit {
	return TileHeightUnit(info & 0x1F)
}

// WithAbsoluteHeight returns a floor info that specifies the given absolute height.
// The allowed range is [0..TileHeightUnitMax-1].
func (info FloorInfo) WithAbsoluteHeight(value TileHeightUnit) FloorInfo {
	return FloorInfo(byte(info&0xE0) | (byte(value) & 0x1F))
}

// CeilingInfo describes the properties of a ceiling
type CeilingInfo byte

// AbsoluteHeight returns the height (from minimum floor height zero) of the ceiling.
// Internally, this value is stored "from maximum ceiling height 32" down.
// The range is [1..TileHeightUnitMax].
func (info CeilingInfo) AbsoluteHeight() TileHeightUnit {
	value := info & 0x1F
	return TileHeightUnit(byte(TileHeightUnitMax) - byte(value))
}

// WithAbsoluteHeight returns a ceiling info that specifies the given absolute height.
// Internally, this value is stored "from maximum ceiling height 32" down.
// The allowed range is [1..TileHeightUnitMax].
func (info CeilingInfo) WithAbsoluteHeight(value TileHeightUnit) CeilingInfo {
	return CeilingInfo(byte(info&0xE0) | (byte(TileHeightUnitMax-value) & 0x1F))
}
