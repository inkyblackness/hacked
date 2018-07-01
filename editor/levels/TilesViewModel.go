package levels

type tilesViewModel struct {
	selectedTiles tileCoordinates

	restoreFocus bool
	windowOpen   bool
}

func freshTilesViewModel() tilesViewModel {
	return tilesViewModel{}
}
