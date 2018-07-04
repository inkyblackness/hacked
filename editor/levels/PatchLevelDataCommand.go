package levels

import (
	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/editor/model"
	"github.com/inkyblackness/hacked/ss1/resource"
)

type stateRestorer func()

type patchLevelDataCommand struct {
	restoreState stateRestorer

	patches []model.BlockPatch
}

func (cmd patchLevelDataCommand) Do(trans cmd.Transaction) error {
	return cmd.perform(trans, cmd.patches, func(p *model.BlockPatch) []byte { return p.ForwardData })
}

func (cmd patchLevelDataCommand) Undo(trans cmd.Transaction) error {
	return cmd.perform(trans, cmd.patches, func(p *model.BlockPatch) []byte { return p.ReverseData })
}

func (cmd patchLevelDataCommand) perform(trans cmd.Transaction, patches []model.BlockPatch, dataResolver func(*model.BlockPatch) []byte) error {
	for _, patch := range patches {
		trans.PatchResourceBlock(resource.LangAny, patch.ID, patch.BlockIndex, patch.BlockLength, dataResolver(&patch))
	}
	cmd.restoreState()
	return nil
}
