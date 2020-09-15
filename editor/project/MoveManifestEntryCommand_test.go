package project

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testingMover struct {
	lastTo    int
	lastFrom  int
	nextError error
}

func (mover *testingMover) MoveEntry(to, from int) error {
	mover.lastTo = to
	mover.lastFrom = from
	return mover.nextError
}

type MoveManifestEntryCommandSuite struct {
	suite.Suite

	instance  moveManifestEntryCommand
	mover     testingMover
	lastError error
}

func TestMoveManifestEntryCommandSuite(t *testing.T) {
	suite.Run(t, new(MoveManifestEntryCommandSuite))
}

func (suite *MoveManifestEntryCommandSuite) SetupTest() {
	suite.mover = testingMover{}
	suite.instance = moveManifestEntryCommand{
		from:  -1,
		to:    -1,
		mover: &suite.mover,
	}
}

func (suite *MoveManifestEntryCommandSuite) TestDoCallsMoverForward() {
	suite.givenParameters(20, 10)
	suite.whenCommandIsDone()
	suite.thenLastErrorShouldBeNil()
	suite.thenMoverShouldHaveBeenCalledWith(20, 10)
}

func (suite *MoveManifestEntryCommandSuite) TestUndoCallsMoverBackward() {
	suite.givenParameters(20, 10)
	suite.whenCommandIsUndone()
	suite.thenLastErrorShouldBeNil()
	suite.thenMoverShouldHaveBeenCalledWith(10, 20)
}

func (suite *MoveManifestEntryCommandSuite) TestCommandReturnsErrorIfMoverDoes() {
	suite.mover.nextError = errors.New("some error")
	suite.whenCommandIsDone()
	suite.thenLastErrorShouldBe(suite.mover.nextError)
}

func (suite *MoveManifestEntryCommandSuite) givenParameters(to int, from int) {
	suite.instance.to = to
	suite.instance.from = from
}

func (suite *MoveManifestEntryCommandSuite) whenCommandIsDone() {
	suite.lastError = suite.instance.Do(nil)
}

func (suite *MoveManifestEntryCommandSuite) whenCommandIsUndone() {
	suite.lastError = suite.instance.Undo(nil)
}

func (suite *MoveManifestEntryCommandSuite) thenLastErrorShouldBeNil() {
	assert.Nil(suite.T(), suite.lastError, "No error expected")
}

func (suite *MoveManifestEntryCommandSuite) thenLastErrorShouldBe(expected error) {
	assert.Equal(suite.T(), expected, suite.lastError, "Error expected")
}

func (suite *MoveManifestEntryCommandSuite) thenMoverShouldHaveBeenCalledWith(to int, from int) {
	assert.Equal(suite.T(), to, suite.mover.lastTo, "TO mismatch")
	assert.Equal(suite.T(), from, suite.mover.lastFrom, "FROM mismatch")
}
