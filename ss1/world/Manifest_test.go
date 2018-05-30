package world_test

import (
	"fmt"
	"testing"

	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ManifestSuite struct {
	suite.Suite
	manifest *world.Manifest
}

func TestManifestSuite(t *testing.T) {
	suite.Run(t, new(ManifestSuite))
}

func (suite *ManifestSuite) SetupTest() {
	suite.manifest = new(world.Manifest)
	suite.manifest.Modified = suite.onManifestModified
}

func (suite *ManifestSuite) TestEntriesCanBeAdded() {
	suite.whenEntryIsInserted(0, "id1")
	suite.thenEntryCountShouldBe(1)
	suite.thenEntryAtShouldBe(0, "id1")
}

func (suite *ManifestSuite) TestEntryReturnsErrorOnInvalidInput() {
	suite.givenEntryWasInserted(0, "id1")

	assertErrorAt := func(at int, message string) {
		_, err := suite.manifest.Entry(at)
		assert.Error(suite.T(), err, message)
	}
	assertErrorAt(-1, "Error expected with negative index")
	assertErrorAt(1, "Error expected with index beyond limit")
}

func (suite *ManifestSuite) TestEntriesCanBeInserted() {
	suite.givenEntryWasInserted(0, "id1")
	suite.givenEntryWasInserted(1, "id2")
	suite.whenEntryIsInserted(1, "id3")
	suite.thenEntryCountShouldBe(3)
	suite.thenEntryAtShouldBe(0, "id1")
	suite.thenEntryAtShouldBe(1, "id3")
	suite.thenEntryAtShouldBe(2, "id2")
}

func (suite *ManifestSuite) TestInsertReturnsErrorsOnInvalidData() {
	assert.Error(suite.T(), suite.manifest.InsertEntry(-1, suite.aSimpleEntry("-1")), "Error expected with negative index")
	assert.Error(suite.T(), suite.manifest.InsertEntry(1, suite.aSimpleEntry("1")), "Error expected with index beyond count")
	assert.Error(suite.T(), suite.manifest.InsertEntry(0, nil), "Error expected for nil entry")
}

func (suite *ManifestSuite) TestEntriesCanBeRemoved() {
	suite.givenEntryWasInserted(0, "id1")
	suite.givenEntryWasInserted(1, "id2")
	suite.givenEntryWasInserted(2, "id3")
	suite.whenEntryIsRemoved(1)
	suite.thenEntryCountShouldBe(2)
	suite.thenEntryAtShouldBe(0, "id1")
	suite.thenEntryAtShouldBe(1, "id3")

	suite.whenEntryIsRemoved(1)
	suite.thenEntryCountShouldBe(1)
	suite.thenEntryAtShouldBe(0, "id1")
}

func (suite *ManifestSuite) TestRemoveReturnsErrorOnInvalidData() {
	suite.givenEntryWasInserted(0, "id1")

	assert.Error(suite.T(), suite.manifest.RemoveEntry(-1), "Error expected with negative index")
	assert.Error(suite.T(), suite.manifest.RemoveEntry(1), "Error expected with index beyond count")
}

func (suite *ManifestSuite) TestEntriesCanBeReplaced() {
	suite.givenEntryWasInserted(0, "id1")
	suite.whenEntryIsReplaced(0, "id1-new")
	suite.thenEntryAtShouldBe(0, "id1-new")
}

func (suite *ManifestSuite) TestReplaceReturnsErrorOnInvalidData() {
	suite.givenEntryWasInserted(0, "id1")

	assert.Error(suite.T(), suite.manifest.ReplaceEntry(-1, suite.aSimpleEntry("-1")), "Error expected with negative index")
	assert.Error(suite.T(), suite.manifest.ReplaceEntry(1, suite.aSimpleEntry("1")), "Error expected with index beyond count")
	assert.Error(suite.T(), suite.manifest.ReplaceEntry(0, nil), "Error expected for nil entry")
}

func (suite *ManifestSuite) TestEntriesCanBeMovedBackward() {
	suite.givenEntryWasInserted(0, "id1")
	suite.givenEntryWasInserted(1, "id2")
	suite.givenEntryWasInserted(2, "id3")
	suite.givenEntryWasInserted(3, "id4")
	suite.givenEntryWasInserted(4, "id5")
	suite.whenEntryIsMoved(1, 3)
	suite.thenEntryOrderShouldBe("id1", "id4", "id2", "id3", "id5")
}

func (suite *ManifestSuite) TestEntriesCanBeMovedForward() {
	suite.givenEntryWasInserted(0, "id1")
	suite.givenEntryWasInserted(1, "id2")
	suite.givenEntryWasInserted(2, "id3")
	suite.givenEntryWasInserted(3, "id4")
	suite.givenEntryWasInserted(4, "id5")
	suite.whenEntryIsMoved(3, 1)
	suite.thenEntryOrderShouldBe("id1", "id3", "id4", "id2", "id5")
}

func (suite *ManifestSuite) TestMoveEntryToSamePlaceDoesNothing() {
	suite.givenEntryWasInserted(0, "id1")
	suite.givenEntryWasInserted(1, "id2")
	suite.givenEntryWasInserted(2, "id3")
	suite.whenEntryIsMoved(1, 1)
	suite.thenEntryOrderShouldBe("id1", "id2", "id3")
}

func (suite *ManifestSuite) TestMoveEntryToEnd() {
	suite.givenEntryWasInserted(0, "id1")
	suite.givenEntryWasInserted(1, "id2")
	suite.givenEntryWasInserted(2, "id3")
	suite.whenEntryIsMoved(2, 0)
	suite.thenEntryOrderShouldBe("id2", "id3", "id1")
}

func (suite *ManifestSuite) TestMoveEntryReturnsErrorOnInvalidInput() {
	suite.givenEntryWasInserted(0, "id1")
	suite.givenEntryWasInserted(1, "id2")
	suite.givenEntryWasInserted(2, "id3")

	assert.Error(suite.T(), suite.manifest.MoveEntry(0, -1), "Error expected moving from negative index")
	assert.Error(suite.T(), suite.manifest.MoveEntry(0, 3), "Error expected moving from index beyond limit")
	assert.Error(suite.T(), suite.manifest.MoveEntry(-1, 1), "Error expected moving to negative index")
	assert.Error(suite.T(), suite.manifest.MoveEntry(3, 1), "Error expected moving to index beyond limit")
}

func (suite *ManifestSuite) onManifestModified(modifiedIDs []resource.ID, failedIDs []resource.ID) {

}

func (suite *ManifestSuite) aSimpleEntry(id string) *world.ManifestEntry {
	return &world.ManifestEntry{
		ID: id,
	}
}

func (suite *ManifestSuite) givenEntryWasInserted(at int, id string) {
	suite.whenEntryIsInserted(at, id)
}

func (suite *ManifestSuite) whenEntryIsInserted(at int, id string) {
	err := suite.manifest.InsertEntry(at, suite.aSimpleEntry(id))
	require.Nil(suite.T(), err, fmt.Sprintf("No error expected inserting entry at %d", at))
}

func (suite *ManifestSuite) whenEntryIsRemoved(at int) {
	err := suite.manifest.RemoveEntry(at)
	require.Nil(suite.T(), err, fmt.Sprintf("No error expected removing entry at %d", at))
}

func (suite *ManifestSuite) whenEntryIsReplaced(at int, id string) {
	err := suite.manifest.ReplaceEntry(at, suite.aSimpleEntry(id))
	require.Nil(suite.T(), err, fmt.Sprintf("No error expected replacing entry at %d", at))
}

func (suite *ManifestSuite) whenEntryIsMoved(to int, from int) {
	err := suite.manifest.MoveEntry(to, from)
	require.Nil(suite.T(), err, fmt.Sprintf("No error expected moving entry to %d from %d", to, from))
}

func (suite *ManifestSuite) thenEntryCountShouldBe(expected int) {
	assert.Equal(suite.T(), expected, suite.manifest.EntryCount())
}

func (suite *ManifestSuite) thenEntryAtShouldBe(at int, id string) {
	entry, err := suite.manifest.Entry(at)
	assert.Nil(suite.T(), err, fmt.Sprintf("No error expected retrieving entry at %d", at))
	require.NotNil(suite.T(), entry, fmt.Sprintf("Entry expected at %d", at))
	assert.Equal(suite.T(), id, entry.ID, fmt.Sprintf("Wrong entry found at %d", at))
}

func (suite *ManifestSuite) thenEntryOrderShouldBe(expected ...string) {
	expectedLen := len(expected)
	require.Equal(suite.T(), expectedLen, suite.manifest.EntryCount(), "Invalid number of entries")
	currentIDs := make([]string, expectedLen)
	for index := 0; index < expectedLen; index++ {
		entry, err := suite.manifest.Entry(index)
		require.Nil(suite.T(), err, fmt.Sprintf("No error expected retrieving entry at %d", index))
		currentIDs[index] = entry.ID
	}
	assert.Equal(suite.T(), expected, currentIDs, "IDs don't not match")
}
