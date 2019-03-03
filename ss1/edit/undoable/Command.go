package undoable

import (
	"github.com/inkyblackness/hacked/ss1/world"
)

type command struct {
	forward func(world.Modder)
	reverse func(world.Modder)
	restore func()
}

func (c command) Do(modder world.Modder) error {
	return c.perform(modder, c.forward)
}

func (c command) Undo(modder world.Modder) error {
	return c.perform(modder, c.reverse)
}

func (c command) perform(modder world.Modder, callback func(setter world.Modder)) error {
	callback(modder)
	c.restore()
	return nil
}
