package cmd

import "github.com/inkyblackness/hacked/ss1/world"

// Command describes an action that can be performed, undone, and redone.
type Command interface {
	// Do performs the action.
	// An error is returned if the action could not be performed.
	Do(modder world.Modder) error

	// Undo reverses its action to restore the environment to the previous state.
	// An error is returned if the action could not be successfully undone. The
	// environment may not be in the state as before in an error occurred.
	Undo(modder world.Modder) error
}
