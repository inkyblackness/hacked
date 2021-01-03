package levels

import "github.com/inkyblackness/hacked/ss1/content/archive/level"

// MapPosition describes a specific two-dimensional point on the map.
type MapPosition struct {
	X level.Coordinate
	Y level.Coordinate
}

// Tile returns the tile position of the exact position.
func (pos MapPosition) Tile() level.TilePosition {
	return level.TilePosition{X: pos.X.Tile(), Y: pos.Y.Tile()}
}
