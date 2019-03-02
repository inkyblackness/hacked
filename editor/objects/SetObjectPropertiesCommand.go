package objects

import (
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/world"
)

type setObjectPropertiesCommand struct {
	model *viewModel

	triple object.Triple
	bitmap int

	oldProperties object.Properties
	newProperties object.Properties
}

func (command setObjectPropertiesCommand) Do(modder world.Modder) error {
	return command.perform(modder, command.newProperties)
}

func (command setObjectPropertiesCommand) Undo(modder world.Modder) error {
	return command.perform(modder, command.oldProperties)
}

func (command setObjectPropertiesCommand) perform(modder world.Modder, properties object.Properties) error {
	modder.SetObjectProperties(command.triple, properties)

	command.model.restoreFocus = true
	command.model.currentObject = command.triple
	command.model.currentBitmap = command.bitmap
	return nil
}
