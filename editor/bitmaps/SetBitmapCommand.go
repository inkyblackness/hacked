package bitmaps

import (
	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
)

type setBitmapCommand struct {
	model *viewModel

	displayKey resource.Key

	resourceKey resource.Key
	oldData     []byte
	newData     []byte
}

func (cmd setBitmapCommand) Do(trans cmd.Transaction) error {
	return cmd.perform(trans, cmd.newData)
}

func (cmd setBitmapCommand) Undo(trans cmd.Transaction) error {
	return cmd.perform(trans, cmd.oldData)
}

func (cmd setBitmapCommand) perform(trans cmd.Transaction, data []byte) error {
	trans.SetResourceBlock(cmd.resourceKey.Lang, cmd.resourceKey.ID, cmd.resourceKey.Index, data)

	cmd.model.restoreFocus = true
	cmd.model.currentKey = cmd.displayKey
	return nil
}
