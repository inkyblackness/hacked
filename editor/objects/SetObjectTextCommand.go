package objects

import (
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
)

type setObjectTextCommand struct {
	model *viewModel

	triple object.Triple
	key    resource.Key

	oldData []byte
	newData []byte
}

func (command setObjectTextCommand) Do(trans cmd.Transaction) error {
	return command.perform(trans, command.newData)
}

func (command setObjectTextCommand) Undo(trans cmd.Transaction) error {
	return command.perform(trans, command.oldData)
}

func (command setObjectTextCommand) perform(trans cmd.Transaction, data []byte) error {
	trans.SetResourceBlock(command.key.Lang, command.key.ID, command.key.Index, data)
	command.model.restoreFocus = true
	command.model.currentObject = command.triple
	command.model.currentLang = command.key.Lang
	return nil
}
