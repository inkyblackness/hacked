package levels

import (
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
)

type stateRestorer func(forward bool)

type patchLevelDataCommand struct {
	restoreState stateRestorer

	patches []world.BlockPatch
}

func (cmd patchLevelDataCommand) Do(trans cmd.Transaction) error {
	cmd.perform(trans, cmd.patches, func(p *world.BlockPatch) []byte { return p.ForwardData })
	cmd.restoreState(true)
	return nil
}

func (cmd patchLevelDataCommand) Undo(trans cmd.Transaction) error {
	cmd.perform(trans, cmd.patches, func(p *world.BlockPatch) []byte { return p.ReverseData })
	cmd.restoreState(false)
	return nil
}

func (cmd patchLevelDataCommand) perform(trans cmd.Transaction, patches []world.BlockPatch, dataResolver func(*world.BlockPatch) []byte) {
	for _, patch := range patches {
		trans.PatchResourceBlock(resource.LangAny, patch.ID, patch.BlockIndex, patch.BlockLength, dataResolver(&patch)) // nolint: scopelint
	}
}
