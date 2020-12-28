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

// Undo executes the reverse tasks in series. The first task to return an error aborts the iteration.
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
type TransactionModifier func(*Transaction) error

// Registry allows to build commands as transactions with modifier functions.
type Registry interface {
	Commander
	Register(modifier ...TransactionModifier) error
}

// TransactionBuilder provides a way to register a command at a commander with
// possibly nested actions.
type TransactionBuilder struct {
	Commander Commander
	active    *Transaction
}

// Register creates a new transaction based on given modifier.
func (builder *TransactionBuilder) Register(modifier ...TransactionModifier) error {
	previous := builder.active
	current := &Transaction{}

	builder.active = current
	for _, mod := range modifier {
		err := mod(current)
		if err != nil {
			return err
		}
	}
	builder.active = previous

	// invert reverse task for easier iteration
	for left, right := 0, len(current.Reverse)-1; left < right; left, right = left+1, right-1 {
		current.Reverse[left], current.Reverse[right] = current.Reverse[right], current.Reverse[left]
	}

	builder.Queue(current)
	return nil
}

// Queue implements the Commander interface to add given command to a currently built transaction.
// If there is no outer transaction being built, then the command is immediately forwarded to the commander.
// Avoid building a loop between TransactionBuilder instances, as this will result in a stack overflow.
func (builder *TransactionBuilder) Queue(command Command) {
	if builder.active == nil {
		builder.Commander.Queue(command)
		return
	}

	builder.active.Forward = append(builder.active.Forward, command.Do)
	builder.active.Reverse = append(builder.active.Reverse, command.Undo)
}

// Named allows to name the current level of transaction.
func Named(name string) TransactionModifier {
	return func(txn *Transaction) error {
		txn.Name = name
		return nil
	}
}

// Forward creates a modifier for a task that is called when doing a command.
// If a nested command is added to the transaction after the forward function, then
// the nested task is performed after the given task.
// If a nested command is added to the transaction before the forward function, then
// the nested task is performed before the given task.
func Forward(task Task) TransactionModifier {
	return func(txn *Transaction) error {
		txn.Forward = append(txn.Forward, task)
		return nil
	}
}

// Reverse creates a modifier for a task that is called when undoing a command.
// If a nested command is added to the transaction after the reverse function, then
// the nested task is performed before the given task.
// If a nested command is added to the transaction before the reverse function, then
// the nested task is performed after the given task.
func Reverse(task Task) TransactionModifier {
	return func(txn *Transaction) error {
		txn.Reverse = append(txn.Reverse, task)
		return nil
	}
}

// Nested adds a further level to the currently registered transaction.
// Depending on the order of the modifier in the register, the corresponding tasks will be called
// in sequence. See other modifier.
func Nested(nested func() error) TransactionModifier {
	return func(txn *Transaction) error {
		return nested()
	}
}
