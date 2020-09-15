package project

import (
	"github.com/inkyblackness/hacked/ss1/world"
)

type manifestEntryKeeper interface {
	RemoveEntry(at int) error
	InsertEntry(at int, entries ...*world.ManifestEntry) error
}

type listManifestEntryCommand struct {
	keeper manifestEntryKeeper

	at    int
	entry *world.ManifestEntry

	adder bool
}

func (cmd listManifestEntryCommand) Do(modder world.Modder) error {
	return cmd.perform(true)
}

func (cmd listManifestEntryCommand) Undo(modder world.Modder) error {
	return cmd.perform(false)
}

func (cmd listManifestEntryCommand) perform(insert bool) error {
	err := cmd.callKeeper(insert)
	if err != nil {
		return err
	}
	return nil
}

func (cmd listManifestEntryCommand) callKeeper(insert bool) error {
	if cmd.adder == insert {
		return cmd.keeper.InsertEntry(cmd.at, cmd.entry)
	}
	return cmd.keeper.RemoveEntry(cmd.at)
}
