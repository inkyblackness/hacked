package level

import "fmt"

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
