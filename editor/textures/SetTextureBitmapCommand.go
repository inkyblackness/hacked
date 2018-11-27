package textures

import (
	"errors"

	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

type setTextureBitmapCommand struct {
	model *viewModel

	textureIndex int
	id           resource.ID

	oldData []byte
	newData []byte
}

func (command setTextureBitmapCommand) Do(trans cmd.Transaction) error {
	return command.perform(trans, command.newData)
}

func (command setTextureBitmapCommand) Undo(trans cmd.Transaction) error {
	return command.perform(trans, command.oldData)
}

func (command setTextureBitmapCommand) perform(trans cmd.Transaction, data []byte) error {
	info, existing := ids.Info(command.id)
	if !existing {
		return errors.New("unknown identifier")
	}
	resourceID := command.id
	blockIndex := command.textureIndex
	if !info.List {
		resourceID = resourceID.Plus(blockIndex)
		blockIndex = 0
	}

	if (len(data) > 0) || info.List {
		trans.SetResourceBlock(resource.LangAny, resourceID, blockIndex, data)
	} else {
		trans.DelResource(resource.LangAny, resourceID)
	}

	command.model.restoreFocus = true
	command.model.currentIndex = command.textureIndex
	return nil
}
