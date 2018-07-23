package level

import "errors"

// HeightShift indicates the vertical scale of a level.
type HeightShift int32

// ValueFromTileHeight returns a floating point value in tiles based on a tile height.
func (shift HeightShift) ValueFromTileHeight(raw TileHeightUnit) (float32, error) {
	return shift.valueFromScale(float32(raw), float64(TileHeightUnitMax))
}

// ValueFromObjectHeight returns a floating point value in tiles based on an object height.
func (shift HeightShift) ValueFromObjectHeight(raw HeightUnit) (float32, error) {
	return shift.valueFromScale(float32(raw), float64(0x100))
}

func (shift HeightShift) valueFromScale(raw float32, scale float64) (float32, error) {
	tileHeights := []float64{32.0, 16.0, 8.0, 4.0, 2.0, 1.0, 0.5, 0.25}
	if (shift < 0) || (int(shift) >= len(tileHeights)) {
		return 0.0, errors.New("invalid height shift")
	}
	return float32((float64(raw) * tileHeights[int(shift)]) / scale), nil
}
