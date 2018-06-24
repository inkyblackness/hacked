package level

// TileType describes the general type of a map tile.
type TileType byte

// Tiles come in different forms:
// Solid tiles can not be entered, Open tiles are regular tiles with a flat floor and a flat ceiling.
// DiagonalOpen tiles are those with flat floors and ceilings, and two walls cut off by one diagonal wall.
// Slope tiles have a sloped floor (or ceiling). Valley tiles have one floor vertex lower while Ridge tiles have one
// floor vertex higher than the other three.
const (
	Solid TileType = 0x00
	Open  TileType = 0x01

	DiagonalOpenSouthEast TileType = 0x02
	DiagonalOpenSouthWest TileType = 0x03
	DiagonalOpenNorthWest TileType = 0x04
	DiagonalOpenNorthEast TileType = 0x05

	SlopeSouthToNorth TileType = 0x06
	SlopeWestToEast   TileType = 0x07
	SlopeNorthToSouth TileType = 0x08
	SlopeEastToWest   TileType = 0x09

	ValleySouthEastToNorthWest TileType = 0x0A
	ValleySouthWestToNorthEast TileType = 0x0B
	ValleyNorthWestToSouthEast TileType = 0x0C
	ValleyNorthEastToSouthWest TileType = 0x0D

	RidgeNorthWestToSouthEast TileType = 0x0E
	RidgeNorthEastToSouthWest TileType = 0x0F
	RidgeSouthEastToNorthWest TileType = 0x10
	RidgeSouthWestToNorthEast TileType = 0x11
)
