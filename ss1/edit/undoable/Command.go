package undoable

import (
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
)

type command struct {
	forward  func(cmd.Transaction)
	backward func(cmd.Transaction)
	restore  func()
}

func (cmd command) Do(trans cmd.Transaction) error {
	return cmd.perform(trans, cmd.forward)
}

func (cmd command) Undo(trans cmd.Transaction) error {
	return cmd.perform(trans, cmd.backward)
}

func (cmd command) perform(trans cmd.Transaction, callback func(setter cmd.Transaction)) error {
	callback(trans)
	cmd.restore()
	return nil
}
