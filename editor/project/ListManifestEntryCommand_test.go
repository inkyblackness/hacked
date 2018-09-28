package project

import (
	"testing"

	"errors"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testingKeeper struct {
	lastAt    int
	lastEntry *world.ManifestEntry
	nextError error
}

func (keeper *testingKeeper) InsertEntry(at int, entry ...*world.ManifestEntry) error {
	keeper.lastAt = at
	keeper.lastEntry = entry[0]
	return keeper.nextError
}

func (keeper *testingKeeper) RemoveEntry(at int) error {
	keeper.lastAt = at
	return keeper.nextError
}

type ListManifestEntryCommandSuite struct {
	suite.Suite

	instance  listManifestEntryCommand
	model     viewModel
	keeper    testingKeeper
	lastError error
}

func TestListManifestEntryCommandSuite(t *testing.T) {
	suite.Run(t, new(ListManifestEntryCommandSuite))
}

func (suite *ListManifestEntryCommandSuite) SetupTest() {
	suite.keeper = testingKeeper{}
	suite.model = freshViewModel()
	suite.instance = listManifestEntryCommand{
		keeper: &suite.keeper,
		model:  &suite.model,

		at: -1,
	}
}

func (suite *ListManifestEntryCommandSuite) TestDoReturnsErrorIfKeeperDoes() {
	suite.keeper.nextError = errors.New("some error")
	suite.givenParameters(2, suite.someEntry(), true)
	suite.whenCommandIsDone()
	suite.thenLastErrorShouldBe(suite.keeper.nextError)
}

func (suite *ListManifestEntryCommandSuite) TestDoInsertsEntryIfAdder() {
	entry := suite.someEntry()
	suite.givenParameters(2, entry, true)
	suite.whenCommandIsDone()
	suite.thenLastErrorShouldBeNil()
	suite.thenEntryShouldHaveBeenInserted(2, entry)
}

func (suite *ListManifestEntryCommandSuite) TestUndoRemovesEntryIfAdder() {
	entry := suite.someEntry()
	suite.givenParameters(4, entry, true)
	suite.whenCommandIsUndone()
	suite.thenLastErrorShouldBeNil()
	suite.thenEntryShouldHaveBeenRemoved(4)
}

func (suite *ListManifestEntryCommandSuite) TestDoRemovesEntryIfNotAdder() {
	entry := suite.someEntry()
	suite.givenParameters(3, entry, false)
	suite.whenCommandIsDone()
	suite.thenLastErrorShouldBeNil()
	suite.thenEntryShouldHaveBeenRemoved(3)
}

func (suite *ListManifestEntryCommandSuite) TestUndoInsertsEntryIfNotAdder() {
	entry := suite.someEntry()
	suite.givenParameters(1, entry, false)
	suite.whenCommandIsUndone()
	suite.thenLastErrorShouldBeNil()
	suite.thenEntryShouldHaveBeenInserted(1, entry)
}

func (suite *ListManifestEntryCommandSuite) TestCommandRequestsFocus() {
	suite.givenParameters(2, suite.someEntry(), true)
	suite.whenCommandIsDone()
	suite.thenRestoreFocusShouldBe(true)
}

func (suite *ListManifestEntryCommandSuite) TestCommandSetsSelectedEntryWhenInsertingOnAdd() {
	suite.givenParameters(0, suite.someEntry(), true)
	suite.whenCommandIsDone()
	suite.thenCurrentlySelectedItemShouldBe(0)
}

func (suite *ListManifestEntryCommandSuite) TestCommandSetsSelectedEntryWhenInsertingOnRemove() {
	suite.givenParameters(0, suite.someEntry(), false)
	suite.whenCommandIsUndone()
	suite.thenCurrentlySelectedItemShouldBe(0)
}

func (suite *ListManifestEntryCommandSuite) TestCommandClearsSelectedEntryWhenRemovingOnAdd() {
	suite.givenParameters(1, suite.someEntry(), true)
	suite.whenCommandIsUndone()
	suite.thenCurrentlySelectedItemShouldBe(-1)
}

func (suite *ListManifestEntryCommandSuite) TestCommandClearsSelectedEntryWhenRemovingOnRemove() {
	suite.givenParameters(1, suite.someEntry(), false)
	suite.whenCommandIsDone()
	suite.thenCurrentlySelectedItemShouldBe(-1)
}

func (suite *ListManifestEntryCommandSuite) givenParameters(at int, entry *world.ManifestEntry, adder bool) {
	suite.instance.at = at
	suite.instance.entry = entry
	suite.instance.adder = adder
}

func (suite *ListManifestEntryCommandSuite) whenCommandIsDone() {
	suite.lastError = suite.instance.Do(nil)
}

func (suite *ListManifestEntryCommandSuite) whenCommandIsUndone() {
	suite.lastError = suite.instance.Undo(nil)
}

func (suite *ListManifestEntryCommandSuite) thenLastErrorShouldBeNil() {
	assert.Nil(suite.T(), suite.lastError, "No error expected")
}

func (suite *ListManifestEntryCommandSuite) thenLastErrorShouldBe(expected error) {
	assert.Equal(suite.T(), expected, suite.lastError, "Error expected")
}

func (suite *ListManifestEntryCommandSuite) thenEntryShouldHaveBeenInserted(at int, entry *world.ManifestEntry) {
	assert.Equal(suite.T(), at, suite.keeper.lastAt, "AT mismatch")
	assert.Equal(suite.T(), entry, suite.keeper.lastEntry, "ENTRY mismatch")
}

func (suite *ListManifestEntryCommandSuite) thenEntryShouldHaveBeenRemoved(at int) {
	assert.Equal(suite.T(), at, suite.keeper.lastAt, "AT mismatch")
}

func (suite *ListManifestEntryCommandSuite) thenRestoreFocusShouldBe(expected bool) {
	assert.Equal(suite.T(), expected, suite.model.restoreFocus, "focus not as expected")
}

func (suite *ListManifestEntryCommandSuite) thenCurrentlySelectedItemShouldBe(expected int) {
	assert.Equal(suite.T(), expected, suite.model.selectedManifestEntry, "SELECTED mismatch")
}

func (suite *ListManifestEntryCommandSuite) someEntry() *world.ManifestEntry {
	return &world.ManifestEntry{
		ID: "someEntry",
	}
}
