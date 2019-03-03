package objects

import (
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
)

type setObjectBitmapCommand struct {
	model *viewModel

	triple object.Triple
	bitmap int

	resourceKey resource.Key
	oldData     []byte
	newData     []byte
}

func (command setObjectBitmapCommand) Do(modder world.Modder) error {
	return command.perform(modder, command.newData)
}

func (command setObjectBitmapCommand) Undo(modder world.Modder) error {
	return command.perform(modder, command.oldData)
}

func (command setObjectBitmapCommand) perform(modder world.Modder, data []byte) error {
	modder.SetResourceBlock(command.resourceKey.Lang, command.resourceKey.ID, command.resourceKey.Index, data)

	command.model.restoreFocus = true
	command.model.currentObject = command.triple
	command.model.currentBitmap = command.bitmap
	return nil
}
