package cmd

import "github.com/inkyblackness/hacked/ss1/world"

// List is a sequence of commands, which are executed as one step.
type List []Command

// Do performs the entries in the list in ascending order.
// If an entry returns an error, the iteration is aborted and that error is returned.
func (list List) Do(modder world.Modder) error {
	for _, entry := range list {
		err := entry.Do(modder)
		if err != nil {
			return err
		}
	}
	return nil
}

// Undo performs the entries in the list in descending order.
// If an entry returns an error, the iteration is aborted and that error is returned.
func (list List) Undo(modder world.Modder) error {
	for i := len(list) - 1; i >= 0; i-- {
		err := list[i].Undo(modder)
		if err != nil {
			return err
		}
	}
	return nil
}
