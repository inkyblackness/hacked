package objects

import (
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
)

type setObjectTextCommand struct {
	model *viewModel

	triple object.Triple
	bitmap int
	key    resource.Key

	oldData []byte
	newData []byte
}

func (command setObjectTextCommand) Do(modder world.Modder) error {
	return command.perform(modder, command.newData)
}

func (command setObjectTextCommand) Undo(modder world.Modder) error {
	return command.perform(modder, command.oldData)
}

func (command setObjectTextCommand) perform(modder world.Modder, data []byte) error {
	modder.SetResourceBlock(command.key.Lang, command.key.ID, command.key.Index, data)
	command.model.restoreFocus = true
	command.model.currentObject = command.triple
	command.model.currentBitmap = command.bitmap
	command.model.currentLang = command.key.Lang
	return nil
}
