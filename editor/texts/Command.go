package texts

import (
	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/ss1/cyber"
	"github.com/inkyblackness/hacked/ss1/resource"
)

type command struct {
	model *viewModel

	key resource.Key

	forward  func(setter cyber.AugmentedTextBlockSetter)
	backward func(setter cyber.AugmentedTextBlockSetter)
}

func (cmd command) Do(trans cmd.Transaction) error {
	return cmd.perform(trans, cmd.forward)
}

func (cmd command) Undo(trans cmd.Transaction) error {
	return cmd.perform(trans, cmd.backward)
}

func (cmd command) perform(trans cmd.Transaction, callback func(setter cyber.AugmentedTextBlockSetter)) error {
	callback(trans)

	cmd.model.restoreFocus = true
	cmd.model.currentKey = cmd.key
	return nil
}
