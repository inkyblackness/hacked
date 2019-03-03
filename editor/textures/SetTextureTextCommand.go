package textures

import (
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
)

type setTextureTextCommand struct {
	model *viewModel

	key resource.Key

	oldData []byte
	newData []byte
}

func (command setTextureTextCommand) Do(modder world.Modder) error {
	return command.perform(modder, command.newData)
}

func (command setTextureTextCommand) Undo(modder world.Modder) error {
	return command.perform(modder, command.oldData)
}

func (command setTextureTextCommand) perform(modder world.Modder, data []byte) error {
	modder.SetResourceBlock(command.key.Lang, command.key.ID, command.key.Index, data)
	command.model.restoreFocus = true
	command.model.currentIndex = command.key.Index
	command.model.currentLang = command.key.Lang
	return nil
}
