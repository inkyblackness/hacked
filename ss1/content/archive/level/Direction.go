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

// Offset returns a direction that you reach by offsetting this one by the given value.
// The direction goes in a circle, and both positive and negative offsets are allowed.
func (dir Direction) Offset(offset int) Direction {
	return Direction((int(dir) + ((8 + (offset % 8)) % 8)) % 8)
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
