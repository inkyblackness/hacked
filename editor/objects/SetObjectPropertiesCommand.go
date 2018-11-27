package objects

import (
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
)

type setObjectPropertiesCommand struct {
	model *viewModel

	triple object.Triple

	oldProperties object.Properties
	newProperties object.Properties
}

func (command setObjectPropertiesCommand) Do(trans cmd.Transaction) error {
	return command.perform(trans, command.newProperties)
}

func (command setObjectPropertiesCommand) Undo(trans cmd.Transaction) error {
	return command.perform(trans, command.oldProperties)
}

func (command setObjectPropertiesCommand) perform(trans cmd.Transaction, properties object.Properties) error {
	trans.SetObjectProperties(command.triple, properties)

	command.model.restoreFocus = true
	command.model.currentObject = command.triple
	return nil
}
