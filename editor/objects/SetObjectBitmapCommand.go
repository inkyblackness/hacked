package objects

import (
	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/resource"
)

type setObjectBitmapCommand struct {
	model *viewModel

	triple object.Triple
	bitmap int

	resourceKey resource.Key
	oldData     []byte
	newData     []byte
}

func (command setObjectBitmapCommand) Do(trans cmd.Transaction) error {
	return command.perform(trans, command.newData)
}

func (command setObjectBitmapCommand) Undo(trans cmd.Transaction) error {
	return command.perform(trans, command.oldData)
}

func (command setObjectBitmapCommand) perform(trans cmd.Transaction, data []byte) error {
	trans.SetResourceBlock(command.resourceKey.Lang, command.resourceKey.ID, command.resourceKey.Index, data)

	command.model.restoreFocus = true
	command.model.currentObject = command.triple
	command.model.currentBitmap = command.bitmap
	return nil
}
