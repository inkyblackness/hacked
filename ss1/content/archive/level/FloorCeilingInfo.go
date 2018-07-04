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

// TextureRotations returns the rotation steps to apply for the floor texture. Valid range: [0..3].
func (info FloorInfo) TextureRotations() int {
	return int(info&0x60) >> 5
}

// WithTextureRotations returns a floor info that specifies the given the rotation steps.
// The provided value is normalized to the valid range.
func (info FloorInfo) WithTextureRotations(value int) FloorInfo {
	normalized := (4 + (value % 4)) % 4
	return FloorInfo(byte(info&^0x60) | byte(normalized)<<5)
}

// HasHazard returns true if the hazard flag is set.
func (info FloorInfo) HasHazard() bool {
	return (byte(info) & 0x80) != 0
}

// WithHazard returns a floor info with the hazard flag set according to given value.
func (info FloorInfo) WithHazard(value bool) FloorInfo {
	var flag byte
	if value {
		flag = 0x80
	}
	return FloorInfo(byte(info&^0x80) | flag)
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

// TextureRotations returns the rotation steps to apply for the ceiling texture. Valid range: [0..3].
func (info CeilingInfo) TextureRotations() int {
	return int(info&0x60) >> 5
}

// WithTextureRotations returns a ceiling info that specifies the given the rotation steps.
// The provided value is normalized to the valid range.
func (info CeilingInfo) WithTextureRotations(value int) CeilingInfo {
	normalized := (4 + (value % 4)) % 4
	return CeilingInfo(byte(info&^0x60) | byte(normalized)<<5)
}

// HasHazard returns true if the hazard flag is set.
func (info CeilingInfo) HasHazard() bool {
	return (byte(info) & 0x80) != 0
}

// WithHazard returns a ceiling info with the hazard flag set according to given value.
func (info CeilingInfo) WithHazard(value bool) CeilingInfo {
	var flag byte
	if value {
		flag = 0x80
	}
	return CeilingInfo(byte(info&^0x80) | flag)
}
