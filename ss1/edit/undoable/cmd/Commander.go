package cmd

// Commander tries to execute the given command.
type Commander interface {
	Queue(command Command)
}
