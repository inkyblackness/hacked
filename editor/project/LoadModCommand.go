package project

import (
	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/editor/model"
	"github.com/inkyblackness/hacked/ss1/resource"
)

type loadModCommand struct {
	model *viewModel

	oldModPath   string
	oldResources model.LocalizedResources

	newModPath   string
	newResources model.LocalizedResources
}

func (command loadModCommand) Do(trans cmd.Transaction) error {
	return command.perform(trans, command.newModPath, command.newResources)
}

func (command loadModCommand) Undo(trans cmd.Transaction) error {
	return command.perform(trans, command.oldModPath, command.oldResources)
}

func (command loadModCommand) perform(trans cmd.Transaction, modPath string, resources model.LocalizedResources) error {
	modifiedIDs := make(resource.IDMarkerMap)
	collectIDs := func(res model.LocalizedResources) {
		for _, resMap := range res {
			for id := range resMap {
				modifiedIDs.Add(id)
			}
		}
	}
	collectIDs(command.oldResources)
	collectIDs(command.newResources)

	trans.SetState(modPath, resources, modifiedIDs.ToList())
	command.model.restoreFocus = true
	return nil
}
