package animations

import (
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

type viewModel struct {
	windowOpen   bool
	restoreFocus bool

	currentKey   resource.Key
	currentFrame int
}

func freshViewModel() viewModel {
	return viewModel{
		currentKey: resource.KeyOf(ids.VideoMailAnimationsStart, resource.LangAny, 0),
	}
}
