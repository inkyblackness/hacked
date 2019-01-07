package undoable

import (
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
)

type command struct {
	forward func(cmd.Transaction)
	reverse func(cmd.Transaction)
	restore func()
}

func (c command) Do(trans cmd.Transaction) error {
	return c.perform(trans, c.forward)
}

func (c command) Undo(trans cmd.Transaction) error {
	return c.perform(trans, c.reverse)
}

func (c command) perform(trans cmd.Transaction, callback func(setter cmd.Transaction)) error {
	callback(trans)
	c.restore()
	return nil
}
