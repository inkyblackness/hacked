package texts

import (
	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
)

type setTextLineCommand struct {
	model *viewModel

	key resource.Key

	old []byte
	new []byte
}

func (cmd setTextLineCommand) Do(trans cmd.Transaction) error {
	return cmd.perform(trans, cmd.new)
}

func (cmd setTextLineCommand) Undo(trans cmd.Transaction) error {
	return cmd.perform(trans, cmd.old)
}

func (cmd setTextLineCommand) perform(trans cmd.Transaction, data []byte) error {
	trans.SetResourceBlock(cmd.key.Lang, cmd.key.ID, cmd.key.Index, data)
	cmd.model.restoreFocus = true
	cmd.model.currentKey = cmd.key
	return nil
}
