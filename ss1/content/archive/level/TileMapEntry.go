package level

// TileMapEntry describes one tile of the map.
type TileMapEntry struct {
	// Type indicates what kind of tile this is.
	Type TileType
	// Floor describes floor properties.
	Floor FloorInfo
	// Ceiling describes ceiling properties.
	Ceiling CeilingInfo
	// SlopeHeight indicates for non-flat tiles the height offset.
	SlopeHeight TileHeightUnit

	// FirstObjectIndex points into the level object cross reference table to the first object in this tile.
	FirstObjectIndex int16
	// TextureInfo describes tile texturing.
	TextureInfo TileTextureInfo
	// Flags contains arbitrary additional information.
	Flags TileFlag

	// SubClip is always 0xFF.
	SubClip byte
	_       [2]byte
	// LightDelta describes the light offsets of this tile.
	LightDelta byte
}
