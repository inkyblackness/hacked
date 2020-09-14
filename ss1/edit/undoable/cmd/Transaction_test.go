package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/world"
)

func TestTransactionBuilder(t *testing.T) {
	suite.Run(t, new(TransactionBuilderSuite))
}

type TransactionBuilderSuite struct {
	suite.Suite

	builder cmd.TransactionBuilder
	command cmd.Command
	result  []string
}

func (suite *TransactionBuilderSuite) SetupTest() {
	suite.builder = cmd.TransactionBuilder{Commander: suite}
	suite.command = nil
	suite.result = nil
}

func (suite *TransactionBuilderSuite) TestTransactionCallsForwardTasksFromLeftToRight() {
	suite.givenRegisteredTransaction(
		cmd.Forward(suite.aTask("fa")),
		cmd.Forward(suite.aTask("fb")),
		cmd.Reverse(suite.aTask("ra")),
	)
	suite.whenCommandIsDone()
	suite.thenResultShouldBe([]string{"fa", "fb"})
}

func (suite *TransactionBuilderSuite) TestTransactionCallsReverseTasksFromRightToLeft() {
	suite.givenRegisteredTransaction(
		cmd.Forward(suite.aTask("fa")),
		cmd.Reverse(suite.aTask("ra")),
		cmd.Reverse(suite.aTask("rb")),
	)
	suite.whenCommandIsUndone()
	suite.thenResultShouldBe([]string{"rb", "ra"})
}

func (suite *TransactionBuilderSuite) TestTransactionCallsNestedInForwardDirection() {
	suite.givenRegisteredTransaction(
		cmd.Forward(suite.aTask("fa")),
		cmd.Nested(func() { suite.builder.Register(cmd.Forward(suite.aTask("na"))) }),
		cmd.Forward(suite.aTask("fb")),
	)
	suite.whenCommandIsDone()
	suite.thenResultShouldBe([]string{"fa", "na", "fb"})
}

func (suite *TransactionBuilderSuite) TestTransactionCallsNestedInReverseDirection() {
	suite.givenRegisteredTransaction(
		cmd.Reverse(suite.aTask("ra")),
		cmd.Nested(func() { suite.builder.Register(cmd.Reverse(suite.aTask("na")), cmd.Reverse(suite.aTask("nb"))) }),
		cmd.Reverse(suite.aTask("rb")),
	)
	suite.whenCommandIsUndone()
	suite.thenResultShouldBe([]string{"rb", "nb", "na", "ra"})
}

func (suite *TransactionBuilderSuite) givenRegisteredTransaction(modifier ...cmd.TransactionModifier) {
	suite.builder.Register(modifier...)
}

func (suite *TransactionBuilderSuite) whenCommandIsDone() {
	_ = suite.command.Do(nil)
}

func (suite *TransactionBuilderSuite) whenCommandIsUndone() {
	_ = suite.command.Undo(nil)
}

func (suite *TransactionBuilderSuite) thenResultShouldBe(expected []string) {
	assert.Equal(suite.T(), expected, suite.result)
}

func (suite *TransactionBuilderSuite) Queue(command cmd.Command) {
	require.Nil(suite.T(), suite.command, "Not expecting a second command")
	suite.command = command
}

func (suite *TransactionBuilderSuite) aTask(text string) cmd.Task {
	return func(modder world.Modder) error {
		suite.result = append(suite.result, text)
		return nil
	}
}
