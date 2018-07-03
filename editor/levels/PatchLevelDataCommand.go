package levels

import (
	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
)

type stateRestorer func()

type patchEntry struct {
	id   resource.ID
	data []byte
}

type patchLevelDataCommand struct {
	restoreState stateRestorer

	forwardPatches []patchEntry
	reversePatches []patchEntry
}

func (cmd patchLevelDataCommand) Do(trans cmd.Transaction) error {
	return cmd.perform(trans, cmd.forwardPatches)
}

func (cmd patchLevelDataCommand) Undo(trans cmd.Transaction) error {
	return cmd.perform(trans, cmd.reversePatches)
}

func (cmd patchLevelDataCommand) perform(trans cmd.Transaction, patches []patchEntry) error {
	for _, patch := range patches {
		trans.PatchResourceBlock(resource.LangAny, patch.id, 0, patch.data)
	}
	cmd.restoreState()
	return nil
}
