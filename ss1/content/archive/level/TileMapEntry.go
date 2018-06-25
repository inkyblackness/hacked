package level

import "github.com/inkyblackness/hacked/ss1/serial"

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

// TileMap is a square set of tiles.
// The first index is the Y-axis, the second index the X-axis. This way the map can be serialized quicker.
type TileMap [][]TileMapEntry

// NewTileMap returns a new, initialized map.
func NewTileMap(width, height int) TileMap {
	m := make([][]TileMapEntry, height)
	for y := 0; y < height; y++ {
		m[y] = make([]TileMapEntry, width)
		for x := 0; x < width; x++ {
			m[y][x].SubClip = 0xFF
		}
	}
	return m
}

// Tile returns a pointer to the tile within the map for given position.
func (m TileMap) Tile(x, y int) *TileMapEntry {
	return &m[y][x]
}

// Code serializes the map.
func (m TileMap) Code(coder serial.Coder) {
	for y := 0; y < len(m); y++ {
		coder.Code(m[y])
	}
}
