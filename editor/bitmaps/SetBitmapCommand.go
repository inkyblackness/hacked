package bitmaps

import (
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
)

type setBitmapCommand struct {
	model *viewModel

	displayKey resource.Key

	resourceKey resource.Key
	oldData     []byte
	newData     []byte
}

func (cmd setBitmapCommand) Do(modder world.Modder) error {
	return cmd.perform(modder, cmd.newData)
}

func (cmd setBitmapCommand) Undo(modder world.Modder) error {
	return cmd.perform(modder, cmd.oldData)
}

func (cmd setBitmapCommand) perform(modder world.Modder, data []byte) error {
	modder.SetResourceBlock(cmd.resourceKey.Lang, cmd.resourceKey.ID, cmd.resourceKey.Index, data)

	cmd.model.restoreFocus = true
	cmd.model.currentKey = cmd.displayKey
	return nil
}
