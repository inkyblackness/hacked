package levels

type controlViewModel struct {
	selectedAtlasIndex              int
	selectedSurveillanceObjectIndex int
	selectedTextureAnimationIndex   int

	restoreFocus bool
	windowOpen   bool
}

func freshControlViewModel() controlViewModel {
	return controlViewModel{
		selectedTextureAnimationIndex: 1,
	}
}
