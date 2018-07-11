package levels

import (
	"github.com/inkyblackness/hacked/editor/event"
)

type tileCoordinates struct {
	list []MapPosition
}

func (coords tileCoordinates) contains(pos MapPosition) bool {
	for _, entry := range coords.list {
		if entry == pos {
			return true
		}
	}
	return false
}

func (coords *tileCoordinates) registerAt(registry event.Registry) {
	registry.RegisterHandler(coords.onTileSelectionSetEvent)
	registry.RegisterHandler(coords.onTileSelectionAddEvent)
	registry.RegisterHandler(coords.onTileSelectionRemoveEvent)
}

func (coords *tileCoordinates) onTileSelectionSetEvent(evt TileSelectionSetEvent) {
	coords.list = evt.tiles
}

func (coords *tileCoordinates) onTileSelectionAddEvent(evt TileSelectionAddEvent) {
	coords.list = append(coords.list, evt.tiles...)
}

func (coords *tileCoordinates) onTileSelectionRemoveEvent(evt TileSelectionRemoveEvent) {
	newList := make([]MapPosition, 0, len(coords.list))
	for _, oldEntry := range coords.list {
		keep := true
		for _, removedEntry := range evt.tiles {
			if oldEntry == removedEntry {
				keep = false
			}
		}
		if keep {
			newList = append(newList, oldEntry)
		}
	}
	coords.list = newList
}
