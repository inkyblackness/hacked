package level

import "fmt"

// TileFlag describes simple properties of a map tile.
type TileFlag uint32

// MusicIndex returns the music identifier. Range: [0..15].
func (flag TileFlag) MusicIndex() int {
	return int((flag & 0x0000F000) >> 12)
}

// WithMusicIndex returns a new flag value with the given music index set. Values beyond allowed range are ignored.
func (flag TileFlag) WithMusicIndex(value int) TileFlag {
	if (value < 0) || (value > 15) {
		return flag
	}
	return TileFlag(uint32(flag&^0x0000F000) | (uint32(value) << 12))
}

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

// String returns the textual representation of the value.
func (ctrl TileSlopeControl) String() string {
	switch ctrl {
	case TileSlopeControlCeilingInverted:
		return "CeilingInverted"
	case TileSlopeControlCeilingMirrored:
		return "CeilingMirrored"
	case TileSlopeControlCeilingFlat:
		return "CeilingFlat"
	case TileSlopeControlFloorFlat:
		return "FloorFlat"
	default:
		return fmt.Sprintf("Unknown%02X", int(ctrl))
	}
}

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

// TileSlopeControls returns all control values.
func TileSlopeControls() []TileSlopeControl {
	return []TileSlopeControl{
		TileSlopeControlCeilingInverted,
		TileSlopeControlCeilingMirrored,
		TileSlopeControlCeilingFlat,
		TileSlopeControlFloorFlat,
	}
}
