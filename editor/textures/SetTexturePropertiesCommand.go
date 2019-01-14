package textures

import (
	"github.com/inkyblackness/hacked/ss1/content/texture"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
)

type setTexturePropertiesCommand struct {
	model *viewModel

	textureIndex int

	oldProperties texture.Properties
	newProperties texture.Properties
}

func (cmd setTexturePropertiesCommand) Do(trans cmd.Transaction) error {
	return cmd.perform(trans, cmd.newProperties)
}

func (cmd setTexturePropertiesCommand) Undo(trans cmd.Transaction) error {
	return cmd.perform(trans, cmd.oldProperties)
}

func (cmd setTexturePropertiesCommand) perform(trans cmd.Transaction, properties texture.Properties) error {
	trans.SetTextureProperties(cmd.textureIndex, properties)

	cmd.model.restoreFocus = true
	cmd.model.currentIndex = cmd.textureIndex
	return nil
}
