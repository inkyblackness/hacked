package levels

type tileCoordinate struct {
	x int
	y int
}

type tilesViewModel struct {
	selectedTiles []tileCoordinate

	restoreFocus bool
	windowOpen   bool
}

func freshTilesViewModel() tilesViewModel {
	return tilesViewModel{}
}
