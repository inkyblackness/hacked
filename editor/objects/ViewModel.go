package objects

import (
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/resource"
)

type viewModel struct {
	windowOpen   bool
	restoreFocus bool

	currentObject object.Triple
	currentLang   resource.Language
}

func freshViewModel() viewModel {
	return viewModel{}
}
