package project

import "github.com/inkyblackness/hacked/editor/cmd"

type manifestEntryMover interface {
	MoveEntry(to, from int) error
}

type moveManifestEntryCommand struct {
	mover manifestEntryMover
	model *viewModel
	from  int
	to    int
}

func (cmd moveManifestEntryCommand) Do(trans cmd.Transaction) error {
	return cmd.move(cmd.to, cmd.from)
}

func (cmd moveManifestEntryCommand) Undo(trans cmd.Transaction) error {
	return cmd.move(cmd.from, cmd.to)
}

func (cmd moveManifestEntryCommand) move(target, source int) error {
	err := cmd.mover.MoveEntry(target, source)
	if err != nil {
		return err
	}
	cmd.model.selectedManifestEntry = target
	cmd.model.restoreFocus = true
	return nil
}
