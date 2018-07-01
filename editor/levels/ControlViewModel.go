package levels

import "github.com/inkyblackness/hacked/ss1/world"

type controlViewModel struct {
	selectedLevel int

	restoreFocus bool
	windowOpen   bool
}

func freshControlViewModel() controlViewModel {
	return controlViewModel{
		selectedLevel: world.StartingLevel,
	}
}
