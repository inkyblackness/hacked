package level

// FloorInfo describes the properties of a floor
type FloorInfo byte

// AbsoluteHeight returns the floor height.
func (info FloorInfo) AbsoluteHeight() TileHeightUnit {
	return TileHeightUnit(info & 0x1F)
}

// CeilingInfo describes the properties of a ceiling
type CeilingInfo byte

// AbsoluteHeight returns the height (from minimum floor height zero) of the ceiling.
// Internally, this value is stored "from maximum ceiling height 32" down.
func (info CeilingInfo) AbsoluteHeight() TileHeightUnit {
	value := info & 0x1F
	return TileHeightUnit(byte(TileHeightUnitMax) - byte(value))
}
