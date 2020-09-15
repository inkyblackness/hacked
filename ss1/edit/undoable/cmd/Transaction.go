package cmd

import (
	"fmt"

	"github.com/inkyblackness/hacked/ss1/world"
)

// Task represents a modification on given modder.
type Task func(modder world.Modder) error

// TransactionError is returned from a transaction execution if a task failed.
type TransactionError struct {
	// Name is that of the transaction, if it was set.
	Name string
	// Index identifies the list index of tasks.
	Index int
	// Nested is the error returned from the task. In case of nested commands, this may be another TransactionError.
	Nested error
}

// Error returns the path of the transaction task, combined with the nested error text.
func (err TransactionError) Error() string {
	return fmt.Sprintf("%s[%d]: %v", err.Name, err.Index, err.Nested)
}

// Unwrap returns the nested error, to support diving into nested errors.
func (err TransactionError) Unwrap() error {
	return err.Nested
}

// Transaction is a named list of forward and reverse tasks.
// Transaction provides execution of the tasks as a Command.
type Transaction struct {
	Name    string
	Forward []Task
	Reverse []Task
}

// Do executes the forward tasks in series. The first task to return an error aborts the iteration.
// A TransactionError is returned in case of aborted iteration.
func (txn Transaction) Do(modder world.Modder) error {
	return txn.perform(txn.Forward, modder)
}

// Do executes the reverse tasks in series. The first task to return an error aborts the iteration.
// A TransactionError is returned in case of aborted iteration.
func (txn Transaction) Undo(modder world.Modder) error {
	return txn.perform(txn.Reverse, modder)
}

func (txn Transaction) perform(tasks []Task, modder world.Modder) error {
	for i, task := range tasks {
		err := task(modder)
		if err != nil {
			return TransactionError{
				Name:   txn.Name,
				Index:  i,
				Nested: err,
			}
		}
	}
	return nil
}

// TransactionModifier is a function to change a transaction while it is being built.
type TransactionModifier func(*Transaction)

// TransactionBuilder provides a way to register a command at a commander with
// possibly nested actions.
type TransactionBuilder struct {
	Commander Commander
	active    *Transaction
}

// Register creates a new transaction based on given modifier.
func (logger *TransactionBuilder) Register(modifier ...TransactionModifier) {
	previous := logger.active
	current := &Transaction{}

	logger.active = current
	for _, mod := range modifier {
		mod(current)
	}
	logger.active = previous

	// invert reverse task for easier iteration
	for left, right := 0, len(current.Reverse)-1; left < right; left, right = left+1, right-1 {
		current.Reverse[left], current.Reverse[right] = current.Reverse[right], current.Reverse[left]
	}

	logger.Queue(current)
}

// Queue implements the Commander interface to add given command to a currently built transaction.
// If there is no outer transaction being built, then the command is immediately forwarded to the commander.
// Avoid building a loop between TransactionBuilder instances, as this will result in a stack overflow.
func (logger *TransactionBuilder) Queue(command Command) {
	if logger.active == nil {
		logger.Commander.Queue(command)
		return
	}

	logger.active.Forward = append(logger.active.Forward, command.Do)
	logger.active.Reverse = append(logger.active.Reverse, command.Undo)
}

// Named allows to name the current level of transaction.
func Named(name string) TransactionModifier {
	return func(txn *Transaction) {
		txn.Name = name
	}
}

// Forward creates a modifier for a task that is called when doing a command.
// If a nested command is added to the transaction after the forward function, then
// the nested task is performed after the given task.
// If a nested command is added to the transaction before the forward function, then
// the nested task is performed before the given task.
func Forward(task Task) TransactionModifier {
	return func(txn *Transaction) {
		txn.Forward = append(txn.Forward, task)
	}
}

// Reverse creates a modifier for a task that is called when undoing a command.
// If a nested command is added to the transaction after the reverse function, then
// the nested task is performed before the given task.
// If a nested command is added to the transaction before the reverse function, then
// the nested task is performed after the given task.
func Reverse(task Task) TransactionModifier {
	return func(txn *Transaction) {
		txn.Reverse = append(txn.Reverse, task)
	}
}

// Nested adds a further level to the currently registered transaction.
// Depending on the order of the modifier in the register, the corresponding tasks will be called
// in sequence. See other modifier.
func Nested(nested func()) TransactionModifier {
	return func(txn *Transaction) {
		nested()
	}
}
