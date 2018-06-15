package cmd

// Command describes an action that can be performed, undone, and redone.
type Command interface {
	// Do performs the action.
	// An error is returned if the action could not be performed.
	Do(trans Transaction) error

	// Undo reverses its action to restore the environment to the previous state.
	// An error is returned if the action could not be successfully undone. The
	// environment may not be in the state as before in an error occurred.
	Undo(trans Transaction) error
}
