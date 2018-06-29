package level

// Direction is an enumeration of cardinal and ordinal directions.
// Cardinal directions describe sides, ordinal directions describe corners.
type Direction byte

// Plus is a shortcut for dir.AsMask().Plus(other).
func (dir Direction) Plus(other Direction) DirectionMask {
	return dir.AsMask().Plus(other)
}

// AsMask returns the direction as a mask
func (dir Direction) AsMask() DirectionMask {
	return DirectionMask(1 << dir)
}

// DirectionMask is a bitfield combination of directions.
type DirectionMask byte

// Plus adds the mask of the given direction to this mask and returns the result.
func (mask DirectionMask) Plus(dir Direction) DirectionMask {
	return mask | dir.AsMask()
}

// Direction constants
const (
	DirNorth Direction = iota
	DirNorthEast
	DirEast
	DirSouthEast
	DirSouth
	DirSouthWest
	DirWest
	DirNorthWest
)

// SlopeFactors is a list of multipliers for each direction of a tile.
type SlopeFactors [8]float32

// TileTypeInfo is the meta information about a tile type.
type TileTypeInfo struct {
	// SolidSides is a bitfield of cardinal directions describing solid walls.
	SolidSides DirectionMask

	// SlopeFloorFactors defines how a slope affects the floor in each direction of a tile [0.0 .. 1.0].
	SlopeFloorFactors SlopeFactors

	// SlopeMirrorType is the type that mirrors the slope to form a solid tile if merged.
	// Types that have no slope mirror themselves.
	SlopeMirrorType TileType
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
