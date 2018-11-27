package levels

import (
	"github.com/inkyblackness/hacked/editor/model"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
)

type stateRestorer func(forward bool)

type patchLevelDataCommand struct {
	restoreState stateRestorer

	patches []model.BlockPatch
}

func (cmd patchLevelDataCommand) Do(trans cmd.Transaction) error {
	cmd.perform(trans, cmd.patches, func(p *model.BlockPatch) []byte { return p.ForwardData })
	cmd.restoreState(true)
	return nil
}

func (cmd patchLevelDataCommand) Undo(trans cmd.Transaction) error {
	cmd.perform(trans, cmd.patches, func(p *model.BlockPatch) []byte { return p.ReverseData })
	cmd.restoreState(false)
	return nil
}

func (cmd patchLevelDataCommand) perform(trans cmd.Transaction, patches []model.BlockPatch, dataResolver func(*model.BlockPatch) []byte) {
	for _, patch := range patches {
		trans.PatchResourceBlock(resource.LangAny, patch.ID, patch.BlockIndex, patch.BlockLength, dataResolver(&patch))
	}
}
