package level

// SlopeFactors is a list of multipliers for each direction of a tile.
type SlopeFactors [8]float32

// Negated returns the factors multiplied by -1.
func (factors SlopeFactors) Negated() SlopeFactors {
	return SlopeFactors{
		factors[0] * -1, factors[1] * -1, factors[2] * -1, factors[3] * -1,
		factors[4] * -1, factors[5] * -1, factors[6] * -1, factors[7] * -1,
	}
}

// TileTypeInfo is the meta information about a tile type.
type TileTypeInfo struct {
	// Name is the textual representation of the tile type.
	Name string

	// SolidSides is a bitfield of cardinal directions describing solid walls.
	SolidSides DirectionMask

	// SlopeFloorFactors defines how a slope affects the floor in each direction of a tile [0.0 .. 1.0].
	SlopeFloorFactors SlopeFactors

	// SlopeInvertedType is the type that inverts the slope to form a solid tile if merged (e.g. floor & ceiling).
	// Types that have no slope invert themselves.
	SlopeInvertedType TileType
}

var tileTypeInfoList = []TileTypeInfo{
	{"Solid", DirNorth.Plus(DirEast).Plus(DirSouth).Plus(DirWest), SlopeFactors{0, 0, 0, 0, 0, 0, 0, 0}, TileTypeSolid},
	{"Open", 0, SlopeFactors{0, 0, 0, 0, 0, 0, 0, 0}, TileTypeOpen},

	{"DiagonalOpenSouthEast", DirNorth.Plus(DirWest), SlopeFactors{0, 0, 0, 0, 0, 0, 0, 0}, TileTypeDiagonalOpenSouthEast},
	{"DiagonalOpenSouthWest", DirNorth.Plus(DirEast), SlopeFactors{0, 0, 0, 0, 0, 0, 0, 0}, TileTypeDiagonalOpenSouthWest},
	{"DiagonalOpenNorthWest", DirSouth.Plus(DirEast), SlopeFactors{0, 0, 0, 0, 0, 0, 0, 0}, TileTypeDiagonalOpenNorthWest},
	{"DiagonalOpenNorthEast", DirSouth.Plus(DirWest), SlopeFactors{0, 0, 0, 0, 0, 0, 0, 0}, TileTypeDiagonalOpenNorthEast},

	{"SlopeSouthToNorth", 0, SlopeFactors{1.0, 1.0, 0.5, 0.0, 0.0, 0.0, 0.5, 1.0}, TileTypeSlopeNorthToSouth},
	{"SlopeWestToEast", 0, SlopeFactors{0.5, 1.0, 1.0, 1.0, 0.5, 0.0, 0.0, 0.0}, TileTypeSlopeEastToWest},
	{"SlopeNorthToSouth", 0, SlopeFactors{0.0, 0.0, 0.5, 1.0, 1.0, 1.0, 0.5, 0.0}, TileTypeSlopeSouthToNorth},
	{"SlopeEastToWest", 0, SlopeFactors{0.5, 0.0, 0.0, 0.0, 0.5, 1.0, 1.0, 1.0}, TileTypeSlopeWestToEast},

	{"ValleySouthEastToNorthWest", 0, SlopeFactors{1.0, 1.0, 0.5, 0.0, 0.5, 1.0, 1.0, 1.0}, TileTypeRidgeNorthWestToSouthEast},
	{"ValleySouthWestToNorthEast", 0, SlopeFactors{1.0, 1.0, 1.0, 1.0, 0.5, 0.0, 0.5, 1.0}, TileTypeRidgeNorthEastToSouthWest},
	{"ValleyNorthWestToSouthEast", 0, SlopeFactors{0.5, 1.0, 1.0, 1.0, 1.0, 1.0, 0.5, 0.0}, TileTypeRidgeSouthEastToNorthWest},
	{"ValleyNorthEastToSouthWest", 0, SlopeFactors{0.5, 0.0, 0.5, 1.0, 1.0, 1.0, 1.0, 1.0}, TileTypeRidgeSouthWestToNorthEast},

	{"RidgeNorthWestToSouthEast", 0, SlopeFactors{0.0, 0.0, 0.5, 1.0, 0.5, 0.0, 0.0, 0.0}, TileTypeValleySouthEastToNorthWest},
	{"RidgeNorthEastToSouthWest", 0, SlopeFactors{0.0, 0.0, 0.0, 0.0, 0.5, 1.0, 0.5, 0.0}, TileTypeValleySouthWestToNorthEast},
	{"RidgeSouthEastToNorthWest", 0, SlopeFactors{0.5, 0.0, 0.0, 0.0, 0.0, 0.0, 0.5, 1.0}, TileTypeValleyNorthWestToSouthEast},
	{"RidgeSouthWestToNorthEast", 0, SlopeFactors{0.5, 1.0, 0.5, 0.0, 0.0, 0.0, 0.0, 0.0}, TileTypeValleyNorthEastToSouthWest},
}
