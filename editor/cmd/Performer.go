package cmd

// Performer tries ot execute the given command.
type Performer interface {
	Perform(command Command) error
}
