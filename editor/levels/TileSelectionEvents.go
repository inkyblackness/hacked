package levels

// TileSelectionSetEvent notifies about the current set of selected tiles.
type TileSelectionSetEvent struct {
	tiles []MapPosition
}

// TileSelectionAddEvent notifies about added tiles to the current selection.
type TileSelectionAddEvent struct {
	tiles []MapPosition
}

// TileSelectionRemoveEvent notifies about removed tiles from the current selection.
type TileSelectionRemoveEvent struct {
	tiles []MapPosition
}
