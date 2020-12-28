package level

import (
	"math"
)

var tileHeights = []float64{32.0, 16.0, 8.0, 4.0, 2.0, 1.0, 0.5, 0.25}

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

// ValueToObjectHeight returns an object height value based on given floating point value.
func (shift HeightShift) ValueToObjectHeight(value float32) HeightUnit {
	return HeightUnit(math.Min(0xFF, math.Max(0, shift.valueToScale(value, float64(0x100)))))
}

func (shift HeightShift) valueFromScale(raw float32, scale float64) (float32, error) {
	if (shift < 0) || (int(shift) >= len(tileHeights)) {
		return 0.0, errInvalidHeightShift
	}
	return float32((float64(raw) * tileHeights[int(shift)]) / scale), nil
}

func (shift HeightShift) valueToScale(value float32, scale float64) float64 {
	if (shift < 0) || (int(shift) >= len(tileHeights)) {
		return 0.0
	}
	return (float64(value) * scale) / tileHeights[shift]
}
