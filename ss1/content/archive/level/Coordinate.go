package level

// FineCoordinatesPerTileSide is the number of possible positions along an axis within a tile.
const FineCoordinatesPerTileSide = 256

// Coordinate describes a tile position in one dimension.
type Coordinate uint16

// CoordinateAt returns a coordinate with given parameters.
func CoordinateAt(tile, fine byte) Coordinate {
	return Coordinate(uint16(tile)<<8 | uint16(fine))
}

// Tile returns the tile part of the coordinate.
func (coord Coordinate) Tile() byte {
	return byte(coord >> 8)
}

// Fine returns the fine part of the coordinate.
func (coord Coordinate) Fine() byte {
	return byte(coord & 0xFF)
}
