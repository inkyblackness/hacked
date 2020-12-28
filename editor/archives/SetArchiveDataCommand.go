package archives

import (
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
)

type setArchiveDataCommand struct {
	model *viewModel

	selectedLevel int

	oldData map[resource.ID][]byte
	newData map[resource.ID][]byte
}

func (command setArchiveDataCommand) Do(modder world.Modder) error {
	command.delResources(modder, command.oldData)
	return command.perform(modder, command.newData)
}

func (command setArchiveDataCommand) Undo(modder world.Modder) error {
	command.delResources(modder, command.newData)
	return command.perform(modder, command.oldData)
}

func (command setArchiveDataCommand) delResources(modder world.Modder, data map[resource.ID][]byte) {
	for id := range data {
		modder.DelResource(resource.LangAny, id)
	}
}

// nolint: interfacer
func (command setArchiveDataCommand) perform(modder world.Modder, data map[resource.ID][]byte) error {
	for id, blockData := range data {
		if len(blockData) > 0 {
			modder.SetResourceBlocks(resource.LangAny, id, [][]byte{blockData})
		} else {
			modder.DelResource(resource.LangAny, id)
		}
	}
	command.model.restoreFocus = true
	command.model.selectedLevel = command.selectedLevel
	return nil
}
