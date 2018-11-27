package textures

import (
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
)

type setTextureTextCommand struct {
	model *viewModel

	key resource.Key

	oldData []byte
	newData []byte
}

func (command setTextureTextCommand) Do(trans cmd.Transaction) error {
	return command.perform(trans, command.newData)
}

func (command setTextureTextCommand) Undo(trans cmd.Transaction) error {
	return command.perform(trans, command.oldData)
}

func (command setTextureTextCommand) perform(trans cmd.Transaction, data []byte) error {
	trans.SetResourceBlock(command.key.Lang, command.key.ID, command.key.Index, data)
	command.model.restoreFocus = true
	command.model.currentIndex = command.key.Index
	command.model.currentLang = command.key.Lang
	return nil
}
