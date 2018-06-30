package level

// Coordinate describes a tile position in one dimension.
type Coordinate uint16

// Tile returns the tile part of the coordinate.
func (coord Coordinate) Tile() byte {
	return byte(coord >> 8)
}

// Fine returns the fine part of the coordinate.
func (coord Coordinate) Fine() byte {
	return byte(coord & 0xFF)
}
