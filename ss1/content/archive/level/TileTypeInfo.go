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
	// SolidSides is a bitfield of cardinal directions describing solid walls.
	SolidSides DirectionMask

	// SlopeFloorFactors defines how a slope affects the floor in each direction of a tile [0.0 .. 1.0].
	SlopeFloorFactors SlopeFactors

	// SlopeInvertedType is the type that inverts the slope to form a solid tile if merged (e.g. floor & ceiling).
	// Types that have no slope invert themselves.
	SlopeInvertedType TileType
}

var tileTypeInfoList = []TileTypeInfo{
	{DirNorth.Plus(DirEast).Plus(DirSouth).Plus(DirWest), SlopeFactors{0, 0, 0, 0, 0, 0, 0, 0}, TileTypeSolid},
	{0, SlopeFactors{0, 0, 0, 0, 0, 0, 0, 0}, TileTypeOpen},

	{DirNorth.Plus(DirWest), SlopeFactors{0, 0, 0, 0, 0, 0, 0, 0}, TileTypeDiagonalOpenSouthEast},
	{DirNorth.Plus(DirEast), SlopeFactors{0, 0, 0, 0, 0, 0, 0, 0}, TileTypeDiagonalOpenSouthWest},
	{DirSouth.Plus(DirEast), SlopeFactors{0, 0, 0, 0, 0, 0, 0, 0}, TileTypeDiagonalOpenNorthWest},
	{DirSouth.Plus(DirWest), SlopeFactors{0, 0, 0, 0, 0, 0, 0, 0}, TileTypeDiagonalOpenNorthEast},

	{0, SlopeFactors{1.0, 1.0, 0.5, 0.0, 0.0, 0.0, 0.5, 1.0}, TileTypeSlopeNorthToSouth},
	{0, SlopeFactors{0.5, 1.0, 1.0, 1.0, 0.5, 0.0, 0.0, 0.0}, TileTypeSlopeEastToWest},
	{0, SlopeFactors{0.0, 0.0, 0.5, 1.0, 1.0, 1.0, 0.5, 0.0}, TileTypeSlopeSouthToNorth},
	{0, SlopeFactors{0.5, 0.0, 0.0, 0.0, 0.5, 1.0, 1.0, 1.0}, TileTypeSlopeWestToEast},

	{0, SlopeFactors{1.0, 1.0, 0.5, 0.0, 0.5, 1.0, 1.0, 1.0}, TileTypeRidgeNorthWestToSouthEast},
	{0, SlopeFactors{1.0, 1.0, 1.0, 1.0, 0.5, 0.0, 0.5, 1.0}, TileTypeRidgeNorthEastToSouthWest},
	{0, SlopeFactors{0.5, 1.0, 1.0, 1.0, 1.0, 1.0, 0.5, 0.0}, TileTypeRidgeSouthEastToNorthWest},
	{0, SlopeFactors{0.5, 0.0, 0.5, 1.0, 1.0, 1.0, 1.0, 1.0}, TileTypeRidgeSouthWestToNorthEast},

	{0, SlopeFactors{0.0, 0.0, 0.5, 1.0, 0.5, 0.0, 0.0, 0.0}, TileTypeValleySouthEastToNorthWest},
	{0, SlopeFactors{0.0, 0.0, 0.0, 0.0, 0.5, 1.0, 0.5, 0.0}, TileTypeValleySouthWestToNorthEast},
	{0, SlopeFactors{0.5, 0.0, 0.0, 0.0, 0.0, 0.0, 0.5, 1.0}, TileTypeValleyNorthWestToSouthEast},
	{0, SlopeFactors{0.5, 1.0, 0.5, 0.0, 0.0, 0.0, 0.0, 0.0}, TileTypeValleyNorthEastToSouthWest},
}
