package texts

import (
	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
)

type textCommand struct {
	model *viewModel

	key resource.Key

	forward  func(trans cmd.Transaction)
	backward func(trans cmd.Transaction)
}

func (cmd textCommand) Do(trans cmd.Transaction) error {
	return cmd.perform(trans, cmd.forward)
}

func (cmd textCommand) Undo(trans cmd.Transaction) error {
	return cmd.perform(trans, cmd.backward)
}

func (cmd textCommand) perform(trans cmd.Transaction, callback func(trans cmd.Transaction)) error {
	callback(trans)

	cmd.model.restoreFocus = true
	cmd.model.currentKey = cmd.key
	return nil
}
