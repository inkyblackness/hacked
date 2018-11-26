package texts

import (
	"github.com/inkyblackness/hacked/ss1/resource"
)

type viewModel struct {
	windowOpen   bool
	restoreFocus bool
	currentKey   resource.Key
}

func freshViewModel() viewModel {
	return viewModel{
		currentKey: resource.KeyOf(knownTextTypes[0].id, resource.LangDefault, 0),
	}
}
