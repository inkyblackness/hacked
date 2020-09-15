package project

import (
	"github.com/inkyblackness/hacked/ss1/world"
)

type manifestEntryMover interface {
	MoveEntry(to, from int) error
}

type moveManifestEntryCommand struct {
	mover manifestEntryMover
	from  int
	to    int
}

func (cmd moveManifestEntryCommand) Do(modder world.Modder) error {
	return cmd.move(cmd.to, cmd.from)
}

func (cmd moveManifestEntryCommand) Undo(modder world.Modder) error {
	return cmd.move(cmd.from, cmd.to)
}

func (cmd moveManifestEntryCommand) move(target, source int) error {
	err := cmd.mover.MoveEntry(target, source)
	if err != nil {
		return err
	}
	return nil
}
