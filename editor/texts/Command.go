package texts

import (
	"github.com/inkyblackness/hacked/ss1/edit"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
)

type command struct {
	model *viewModel

	key resource.Key

	forward  func(setter edit.AugmentedTextBlockSetter)
	backward func(setter edit.AugmentedTextBlockSetter)
}

func (cmd command) Do(trans cmd.Transaction) error {
	return cmd.perform(trans, cmd.forward)
}

func (cmd command) Undo(trans cmd.Transaction) error {
	return cmd.perform(trans, cmd.backward)
}

func (cmd command) perform(trans cmd.Transaction, callback func(setter edit.AugmentedTextBlockSetter)) error {
	callback(trans)

	cmd.model.restoreFocus = true
	cmd.model.currentKey = cmd.key
	return nil
}
