package levels

import (
	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
)

type stateRestorer func()

type patchData struct {
	id         resource.ID
	blockIndex int

	dataOffset int
	data       []byte
}

type patchLevelDataCommand struct {
	restoreState stateRestorer

	oldData []patchData
	newData []patchData
}

func (cmd patchLevelDataCommand) Do(trans cmd.Transaction) error {
	return cmd.perform(trans, cmd.newData)
}

func (cmd patchLevelDataCommand) Undo(trans cmd.Transaction) error {
	return cmd.perform(trans, cmd.oldData)
}

func (cmd patchLevelDataCommand) perform(trans cmd.Transaction, patches []patchData) error {
	for _, patch := range patches {
		trans.PatchResourceBlock(resource.LangAny, patch.id, patch.blockIndex, patch.dataOffset, patch.data)
	}
	cmd.restoreState()
	return nil
}
