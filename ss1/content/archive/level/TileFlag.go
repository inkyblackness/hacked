package level

// TileFlag describes simple properties of a map tile.
type TileFlag uint32

// SlopeControl returns the slope control as per flags.
func (flag TileFlag) SlopeControl() TileSlopeControl {
	return TileSlopeControl((flag & 0x00000C00) >> 10)
}

// WithSlopeControl returns a new flag value with the given slope control set.
func (flag TileFlag) WithSlopeControl(ctrl TileSlopeControl) TileFlag {
	return TileFlag(uint32(flag&^0x00000C00) | (uint32(ctrl&0x3) << 10))
}

// TileSlopeControl defines how the floor and ceiling of a sloped tile should be processed.
type TileSlopeControl byte

// FloorSlopeFactors returns the slope factors for the given tile type as per control constant.
func (ctrl TileSlopeControl) FloorSlopeFactors(tileType TileType) SlopeFactors {
	if ctrl == TileSlopeControlFloorFlat {
		return SlopeFactors{}
	}
	return tileType.Info().SlopeFloorFactors
}

// CeilingSlopeFactors returns the slope factors for the given tile type as per control constant.
func (ctrl TileSlopeControl) CeilingSlopeFactors(tileType TileType) SlopeFactors {
	if ctrl == TileSlopeControlCeilingFlat {
		return SlopeFactors{}
	}
	if ctrl == TileSlopeControlCeilingMirrored {
		return tileType.Info().SlopeFloorFactors
	}
	return tileType.Info().SlopeInvertedType.Info().SlopeFloorFactors
}

// TileSlopeControl constants.
const (
	TileSlopeControlCeilingInverted TileSlopeControl = 0
	TileSlopeControlCeilingMirrored TileSlopeControl = 1
	TileSlopeControlCeilingFlat     TileSlopeControl = 2
	TileSlopeControlFloorFlat       TileSlopeControl = 3
)
