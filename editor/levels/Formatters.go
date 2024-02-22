package levels

import (
	"fmt"

	"github.com/inkyblackness/hacked/ss1/content/archive/level"
)

func tileHeightFormatterFor(levelHeight level.HeightShift) func(int) string {
	return func(value int) string {
		tileHeight, err := levelHeight.ValueFromTileHeight(level.TileHeightUnit(value))
		tileHeightString := hintUnknown
		if err == nil {
			tileHeightString = fmt.Sprintf("%2.3f", tileHeight)
		}
		return tileHeightString + " tile(s) - raw: %d"
	}
}

func objectHeightFormatterFor(levelHeight level.HeightShift) func(int) string {
	return func(value int) string {
		tileHeight, err := levelHeight.ValueFromObjectHeight(level.HeightUnit(value))
		tileHeightString := hintUnknown
		if err == nil {
			tileHeightString = fmt.Sprintf("%2.3f", tileHeight)
		}
		return tileHeightString + " tile(s) - raw: %d"
	}
}

func moveTileHeightFormatterFor(levelHeight level.HeightShift) func(int) string {
	normalFormatter := tileHeightFormatterFor(levelHeight)
	return func(value int) string {
		if (value >= 0) && (value <= int(level.TileHeightUnitMax)) {
			return normalFormatter(value)
		}
		return "Don't change  - raw: %d"
	}
}
