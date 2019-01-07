package texts

import (
	"github.com/inkyblackness/hacked/ss1/edit"
	"github.com/inkyblackness/hacked/ss1/resource"
)

type viewModel struct {
	windowOpen   bool
	restoreFocus bool
	currentKey   resource.Key
}

func freshViewModel() viewModel {
	return viewModel{
		currentKey: resource.KeyOf(edit.KnownTexts()[0].ID, resource.LangDefault, 0),
	}
}
