package texts

import (
	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
)

type setAudioCommand struct {
	model      *viewModel
	restoreKey resource.Key

	dataKey resource.Key
	oldData [][]byte
	newData [][]byte
}

func (cmd setAudioCommand) Do(trans cmd.Transaction) error {
	return cmd.perform(trans, cmd.newData)
}

func (cmd setAudioCommand) Undo(trans cmd.Transaction) error {
	return cmd.perform(trans, cmd.oldData)
}

func (cmd setAudioCommand) perform(trans cmd.Transaction, data [][]byte) error {
	if len(data) > 0 {
		trans.SetResourceBlocks(cmd.dataKey.Lang, cmd.dataKey.ID, data)
	} else {
		trans.DelResource(cmd.dataKey.Lang, cmd.dataKey.ID)
	}
	cmd.model.restoreFocus = true
	cmd.model.currentKey = cmd.restoreKey
	return nil
}
