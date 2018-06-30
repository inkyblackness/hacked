package level

// HeightShift indicates the vertical scale of a level.
type HeightShift int32

// ValueFromTileHeight returns a floating point value in tiles based on a tile height.
func (shift HeightShift) ValueFromTileHeight(raw TileHeightUnit) float32 {
	return shift.valueFromScale(float32(raw), float32(TileHeightUnitMax))
}

func (shift HeightShift) valueFromScale(raw float32, scale float32) float32 {
	tileHeights := []float32{32.0, 16.0, 8.0, 4.0, 2.0, 1.0, 0.5, 0.25}
	result := float32(0)
	if (shift >= 0) && (int(shift) < len(tileHeights)) {
		result = (raw * tileHeights[int(shift)]) / scale
	}
	return result
}
