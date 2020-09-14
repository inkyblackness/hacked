package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/world"
)

type CommanderFunc func(command cmd.Command)

func (f CommanderFunc) Queue(command cmd.Command) {
	f(command)
}

func TestTransactionCallsForwardTasksFromLeftToRight(t *testing.T) {
	var txn cmd.TransactionBuilder
	txn.Commander = CommanderFunc(func(command cmd.Command) { _ = command.Do(nil) })
	var result []string
	aTask := func(text string) cmd.Task {
		return func(modder world.Modder) error {
			result = append(result, text)
			return nil
		}
	}
	txn.Register(cmd.Forward(aTask("fa")), cmd.Forward(aTask("fb")), cmd.Reverse(aTask("ra")))

	assert.Equal(t, []string{"fa", "fb"}, result)
}

func TestTransactionCallsReverseTasksFromRightToLeft(t *testing.T) {
	var txn cmd.TransactionBuilder
	txn.Commander = CommanderFunc(func(command cmd.Command) { _ = command.Undo(nil) })
	var result []string
	aTask := func(text string) cmd.Task {
		return func(modder world.Modder) error {
			result = append(result, text)
			return nil
		}
	}
	txn.Register(cmd.Forward(aTask("fa")), cmd.Reverse(aTask("ra")), cmd.Reverse(aTask("rb")))

	assert.Equal(t, []string{"rb", "ra"}, result)
}

func TestTransactionCallsNestedInForwardDirection(t *testing.T) {
	var txn cmd.TransactionBuilder
	txn.Commander = CommanderFunc(func(command cmd.Command) { _ = command.Do(nil) })
	var result []string
	aTask := func(text string) cmd.Task {
		return func(modder world.Modder) error {
			result = append(result, text)
			return nil
		}
	}
	txn.Register(
		cmd.Forward(aTask("fa")),
		cmd.Nested(func() { txn.Register(cmd.Forward(aTask("na"))) }),
		cmd.Forward(aTask("fb")))

	assert.Equal(t, []string{"fa", "na", "fb"}, result)
}

func TestTransactionCallsNestedInReverseDirection(t *testing.T) {
	var txn cmd.TransactionBuilder
	txn.Commander = CommanderFunc(func(command cmd.Command) { _ = command.Undo(nil) })
	var result []string
	aTask := func(text string) cmd.Task {
		return func(modder world.Modder) error {
			result = append(result, text)
			return nil
		}
	}
	txn.Register(
		cmd.Reverse(aTask("ra")),
		cmd.Nested(func() { txn.Register(cmd.Reverse(aTask("na")), cmd.Reverse(aTask("nb"))) }),
		cmd.Reverse(aTask("rb")))

	assert.Equal(t, []string{"rb", "nb", "na", "ra"}, result)
}
