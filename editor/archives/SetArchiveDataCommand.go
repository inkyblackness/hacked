package archives

import (
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
)

type setArchiveDataCommand struct {
	model *viewModel

	selectedLevel int

	oldData map[resource.ID][]byte
	newData map[resource.ID][]byte
}

func (command setArchiveDataCommand) Do(trans cmd.Transaction) error {
	command.delResources(trans, command.oldData)
	return command.perform(trans, command.newData)
}

func (command setArchiveDataCommand) Undo(trans cmd.Transaction) error {
	command.delResources(trans, command.newData)
	return command.perform(trans, command.oldData)
}

func (command setArchiveDataCommand) delResources(trans cmd.Transaction, data map[resource.ID][]byte) {
	for id := range data {
		trans.DelResource(resource.LangAny, id)
	}
}

func (command setArchiveDataCommand) perform(trans cmd.Transaction, data map[resource.ID][]byte) error {
	for id, blockData := range data {
		if len(blockData) > 0 {
			trans.SetResourceBlocks(resource.LangAny, id, [][]byte{blockData})
		} else {
			trans.DelResource(resource.LangAny, id)
		}
	}
	command.model.restoreFocus = true
	command.model.selectedLevel = command.selectedLevel
	return nil
}
