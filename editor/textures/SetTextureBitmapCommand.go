package textures

import (
	"fmt"

	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

type setTextureBitmapCommand struct {
	model *viewModel

	textureIndex int
	id           resource.ID

	oldData []byte
	newData []byte
}

func (command setTextureBitmapCommand) Do(modder world.Modder) error {
	return command.perform(modder, command.newData)
}

func (command setTextureBitmapCommand) Undo(modder world.Modder) error {
	return command.perform(modder, command.oldData)
}

func (command setTextureBitmapCommand) perform(modder world.Modder, data []byte) error {
	info, existing := ids.Info(command.id)
	if !existing {
		panic(fmt.Sprintf("unknown identifier for bitmap resource: %v", command.id))
	}
	resourceID := command.id
	blockIndex := command.textureIndex
	if !info.List {
		resourceID = resourceID.Plus(blockIndex)
		blockIndex = 0
	}

	if (len(data) > 0) || info.List {
		modder.SetResourceBlock(resource.LangAny, resourceID, blockIndex, data)
	} else {
		modder.DelResource(resource.LangAny, resourceID)
	}

	command.model.restoreFocus = true
	command.model.currentIndex = command.textureIndex
	return nil
}
