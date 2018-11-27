package undoable

import (
	"github.com/inkyblackness/hacked/ss1/edit"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
)

type AugmentedTextService struct {
	wrapped   edit.AugmentedTextService
	commander cmd.Commander
}

func NewAugmentedTextService(wrapped edit.AugmentedTextService, commander cmd.Commander) AugmentedTextService {
	return AugmentedTextService{
		wrapped:   wrapped,
		commander: commander,
	}
}
