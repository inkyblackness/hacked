package textures

import "github.com/inkyblackness/hacked/ss1/resource"

type viewModel struct {
	windowOpen   bool
	restoreFocus bool

	currentLang  resource.Language
	currentIndex int
}

func freshViewModel() viewModel {
	return viewModel{
		currentIndex: 0,
		currentLang:  resource.LangDefault,
	}
}
