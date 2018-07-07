package levels

type tilesViewModel struct {
	selectedTiles  tileCoordinates
	textureDisplay TextureDisplay

	restoreFocus bool
	windowOpen   bool
}

func freshTilesViewModel() tilesViewModel {
	return tilesViewModel{
		textureDisplay: TextureDisplayFloor,
	}
}
