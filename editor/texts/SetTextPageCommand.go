package texts

import (
	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
)

type setTextPageCommand struct {
	model *viewModel

	key resource.Key

	oldData [][]byte
	newData [][]byte
}

func (cmd setTextPageCommand) Do(trans cmd.Transaction) error {
	return cmd.perform(trans, cmd.newData)
}

func (cmd setTextPageCommand) Undo(trans cmd.Transaction) error {
	return cmd.perform(trans, cmd.oldData)
}

func (cmd setTextPageCommand) perform(trans cmd.Transaction, data [][]byte) error {
	id := cmd.key.ID.Plus(cmd.key.Index)
	if len(data) > 0 {
		trans.SetResourceBlocks(cmd.key.Lang, id, data)
	} else {
		trans.DelResource(cmd.key.Lang, id)
	}
	cmd.model.restoreFocus = true
	cmd.model.currentKey = cmd.key
	return nil
}
