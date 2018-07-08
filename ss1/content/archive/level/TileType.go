package level

import "fmt"

// TileType describes the general type of a map tile.
type TileType byte

// Info returns the information associated with the tile type.
func (t TileType) Info() TileTypeInfo {
	if int(t) < len(tileTypeInfoList) {
		return tileTypeInfoList[t]
	}
	info := tileTypeInfoList[TileTypeSolid]
	info.Name = t.String()
	info.SlopeInvertedType = t
	return info
}

// String returns the textual representation of the type.
func (t TileType) String() string {
	if int(t) < len(tileTypeInfoList) {
		return tileTypeInfoList[t].Name
	} else {
		return fmt.Sprintf("Unknown%02X", int(t))
	}
}

// Tiles come in different forms:
// Solid tiles can not be entered, Open tiles are regular tiles with a flat floor and a flat ceiling.
// DiagonalOpen tiles are those with flat floors and ceilings, and two walls cut off by one diagonal wall.
// Slope tiles have a sloped floor (or ceiling). Valley tiles have one floor vertex lower while Ridge tiles have one
// floor vertex higher than the other three.
const (
	TileTypeSolid TileType = 0x00
	TileTypeOpen  TileType = 0x01

	TileTypeDiagonalOpenSouthEast TileType = 0x02
	TileTypeDiagonalOpenSouthWest TileType = 0x03
	TileTypeDiagonalOpenNorthWest TileType = 0x04
	TileTypeDiagonalOpenNorthEast TileType = 0x05

	TileTypeSlopeSouthToNorth TileType = 0x06
	TileTypeSlopeWestToEast   TileType = 0x07
	TileTypeSlopeNorthToSouth TileType = 0x08
	TileTypeSlopeEastToWest   TileType = 0x09

	TileTypeValleySouthEastToNorthWest TileType = 0x0A
	TileTypeValleySouthWestToNorthEast TileType = 0x0B
	TileTypeValleyNorthWestToSouthEast TileType = 0x0C
	TileTypeValleyNorthEastToSouthWest TileType = 0x0D

	TileTypeRidgeNorthWestToSouthEast TileType = 0x0E
	TileTypeRidgeNorthEastToSouthWest TileType = 0x0F
	TileTypeRidgeSouthEastToNorthWest TileType = 0x10
	TileTypeRidgeSouthWestToNorthEast TileType = 0x11
)

// TileTypes returns a list of all known tile types.
func TileTypes() []TileType {
	return []TileType{
		TileTypeSolid, TileTypeOpen,
		TileTypeDiagonalOpenSouthEast, TileTypeDiagonalOpenSouthWest, TileTypeDiagonalOpenNorthWest, TileTypeDiagonalOpenNorthEast,
		TileTypeSlopeSouthToNorth, TileTypeSlopeWestToEast, TileTypeSlopeNorthToSouth, TileTypeSlopeEastToWest,
		TileTypeValleySouthEastToNorthWest, TileTypeValleySouthWestToNorthEast, TileTypeValleyNorthWestToSouthEast, TileTypeValleyNorthEastToSouthWest,
		TileTypeRidgeNorthWestToSouthEast, TileTypeRidgeNorthEastToSouthWest, TileTypeRidgeSouthEastToNorthWest, TileTypeRidgeSouthWestToNorthEast,
	}
}
