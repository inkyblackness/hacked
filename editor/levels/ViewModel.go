package levels

import "github.com/inkyblackness/hacked/ss1/world"

type viewModel struct {
	selectedLevel int

	restoreFocus bool
	windowOpen   bool
}

func freshViewModel() viewModel {
	return viewModel{
		selectedLevel: world.StartingLevel,
	}
}
