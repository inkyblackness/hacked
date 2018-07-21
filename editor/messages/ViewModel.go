package messages

import (
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

type viewModel struct {
	windowOpen   bool
	restoreFocus bool

	currentKey      resource.Key
	showVerboseText bool
}

func freshViewModel() viewModel {
	return viewModel{
		currentKey:      resource.KeyOf(ids.MailsStart, resource.LangDefault, 0),
		showVerboseText: true,
	}
}
