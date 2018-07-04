package level

// TileTextureInfo describes the textures used for a map tile.
type TileTextureInfo uint16

// WallTextureIndex returns the texture index into the texture atlas for the walls.
// Valid range [0..63].
// This property is only valid in real world.
func (info TileTextureInfo) WallTextureIndex() int {
	return int(info & 0x003F)
}

// WithWallTextureIndex returns an info with given index set.
// Values outside valid range are ignored.
func (info TileTextureInfo) WithWallTextureIndex(value int) TileTextureInfo {
	if (value < 0) || (value >= 64) {
		return info
	}
	return TileTextureInfo(uint16(info&^0x003F) | uint16(value&0x003F))
}

// CeilingTextureIndex returns the texture index into the texture atlas for the ceiling.
// Valid range [0..31].
// This property is only valid in real world.
func (info TileTextureInfo) CeilingTextureIndex() int {
	return int(info&0x07C0) >> 6
}

// WithCeilingTextureIndex returns an info with given index set.
// Values outside valid range are ignored.
func (info TileTextureInfo) WithCeilingTextureIndex(value int) TileTextureInfo {
	if (value < 0) || (value >= FloorCeilingTextureLimit) {
		return info
	}
	return TileTextureInfo(uint16(info&^0x07C0) | (uint16(value&0x001F) << 6))
}

// FloorTextureIndex returns the texture index into the texture atlas for the floor.
// Valid range [0..31].
// This property is only valid in real world.
func (info TileTextureInfo) FloorTextureIndex() int {
	return int(info&0xF800) >> 11
}

// WithFloorTextureIndex returns an info with given index set.
// Values outside valid range are ignored.
func (info TileTextureInfo) WithFloorTextureIndex(value int) TileTextureInfo {
	if (value < 0) || (value >= FloorCeilingTextureLimit) {
		return info
	}
	return TileTextureInfo(uint16(info&^0xF800) | (uint16(value&0x001F) << 11))
}

// FloorPaletteIndex returns the palette index for the floor in cyberspace.
func (info TileTextureInfo) FloorPaletteIndex() byte {
	return byte(info)
}

// WithFloorPaletteIndex returns an info with given index set.
func (info TileTextureInfo) WithFloorPaletteIndex(value byte) TileTextureInfo {
	return TileTextureInfo(uint16(info&0xFF00) | uint16(value))
}

// CeilingPaletteIndex returns the palette index for the ceiling in cyberspace.
func (info TileTextureInfo) CeilingPaletteIndex() byte {
	return byte(info >> 8)
}

// WithCeilingPaletteIndex returns an info with given index set.
func (info TileTextureInfo) WithCeilingPaletteIndex(value byte) TileTextureInfo {
	return TileTextureInfo(uint16(info&0x00FF) | (uint16(value) << 8))
}
