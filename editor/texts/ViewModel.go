package texts

import "github.com/inkyblackness/hacked/ss1/resource"

type viewModel struct {
	windowOpen   bool
	restoreFocus bool
	currentKey   resource.Key
}

func freshViewModel() viewModel {
	return viewModel{
		currentKey: resource.KeyOf(0x0867, resource.LangDefault, 0),
	}
}
