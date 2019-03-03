package textures

import (
	"github.com/inkyblackness/hacked/ss1/content/texture"
	"github.com/inkyblackness/hacked/ss1/world"
)

type setTexturePropertiesCommand struct {
	model *viewModel

	textureIndex int

	oldProperties texture.Properties
	newProperties texture.Properties
}

func (cmd setTexturePropertiesCommand) Do(modder world.Modder) error {
	return cmd.perform(modder, cmd.newProperties)
}

func (cmd setTexturePropertiesCommand) Undo(modder world.Modder) error {
	return cmd.perform(modder, cmd.oldProperties)
}

func (cmd setTexturePropertiesCommand) perform(modder world.Modder, properties texture.Properties) error {
	modder.SetTextureProperties(cmd.textureIndex, properties)

	cmd.model.restoreFocus = true
	cmd.model.currentIndex = cmd.textureIndex
	return nil
}
