package movies

import (
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

type viewModel struct {
	windowOpen   bool
	restoreFocus bool

	currentKey          resource.Key
	currentSubtitleLang resource.Language
	currentScene        int
	currentFrame        int
}

func freshViewModel() viewModel {
	return viewModel{
		currentKey:          resource.KeyOf(ids.MovieIntro, resource.LangDefault, 0),
		currentSubtitleLang: resource.LangDefault,
	}
}
