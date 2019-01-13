package world_test

import (
	"fmt"
	"sort"
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
	selector *resource.Selector

	lastModifiedIDs []resource.ID
	lastFailedIDs   []resource.ID
}

func TestManifestSuite(t *testing.T) {
	suite.Run(t, new(ManifestSuite))
}

func (suite *ManifestSuite) SetupTest() {
	suite.manifest = world.NewManifest(suite.onManifestModified)
	suite.selector = nil

	suite.lastModifiedIDs = nil
	suite.lastFailedIDs = nil
}

func (suite *ManifestSuite) TestEntriesCanBeAdded() {
	suite.whenEntryIsInserted(0, "id1")
	suite.thenEntryCountShouldBe(1)
	suite.thenEntryAtShouldBe(0, "id1")
}

func (suite *ManifestSuite) TestMultipleEntriesCanBeAdded() {
	suite.whenEntriesAreInserted(0, "id1", "id2", "id3")
	suite.thenEntryCountShouldBe(3)
	suite.thenEntryAtShouldBe(0, "id1")
	suite.thenEntryAtShouldBe(1, "id2")
	suite.thenEntryAtShouldBe(2, "id3")
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

func (suite *ManifestSuite) TestMultipleEntriesCanBeInserted() {
	suite.givenEntryWasInserted(0, "id1")
	suite.givenEntryWasInserted(1, "id2")
	suite.whenEntriesAreInserted(1, "id3", "id4")
	suite.thenEntryCountShouldBe(4)
	suite.thenEntryAtShouldBe(0, "id1")
	suite.thenEntryAtShouldBe(1, "id3")
	suite.thenEntryAtShouldBe(2, "id4")
	suite.thenEntryAtShouldBe(3, "id2")
}

func (suite *ManifestSuite) TestInsertReturnsErrorsOnInvalidData() {
	assert.Error(suite.T(), suite.manifest.InsertEntry(-1, suite.aSimpleEntry("-1")), "Error expected with negative index")
	assert.Error(suite.T(), suite.manifest.InsertEntry(1, suite.aSimpleEntry("1")), "Error expected with index beyond count")
	assert.Error(suite.T(), suite.manifest.InsertEntry(0, nil), "Error expected for nil entry")
	assert.Error(suite.T(), suite.manifest.InsertEntry(0, suite.aSimpleEntry("1"), nil, suite.aSimpleEntry("2")), "Error expected for nil entry")
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

func (suite *ManifestSuite) TestLocalizedResourcesCanBeRetrieved() {
	suite.givenEntryWasInsertedWith(0, "id1",
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x1234, [][]byte{{0xAA}})))
	suite.givenEntryWasInserted(1, "id2")
	suite.whenResourcesAreQueriedFor(resource.LangAny)
	suite.thenResourcesCanBeSelected(0x1234)
}

func (suite *ManifestSuite) TestModifiedCallbackOnInsertFromEmptyListsNewIDs() {
	suite.whenEntryIsInsertedWith(0, "id1",
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x1000, [][]byte{{0x11}}),
			suite.storing(0x2000, [][]byte{{0x22}})))
	suite.thenModifiedResourcesShouldBe([]int{0x1000, 0x2000})
}

func (suite *ManifestSuite) TestModifiedCallbackOnInsertFromEmptyListsNewIDsMultipleEntries() {
	suite.whenEntriesAreInsertedWith(0, 10,
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x1000, [][]byte{{0x11}}),
			suite.storing(0x2000, [][]byte{{0x22}})),
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x3000, [][]byte{{0x11}}),
			suite.storing(0x4000, [][]byte{{0x22}})))
	suite.thenModifiedResourcesShouldBe([]int{0x1000, 0x2000, 0x3000, 0x4000})
}

func (suite *ManifestSuite) TestModifiedCallbackOnInsertWithNewOnesListsNewIDs() {
	suite.givenEntryWasInsertedWith(0, "id1",
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0800, [][]byte{{0xAA}})))
	suite.whenEntryIsInsertedWith(1, "id2",
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0900, [][]byte{{0x11}}),
			suite.storing(0x0A00, [][]byte{{0x22}})))
	suite.thenModifiedResourcesShouldBe([]int{0x0900, 0x0A00})
}

func (suite *ManifestSuite) TestModifiedCallbackOnInsertIgnoresIdenticalData() {
	suite.givenEntryWasInsertedWith(0, "id1",
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0800, [][]byte{{0xAA}})))
	suite.whenEntryIsInsertedWith(1, "id2",
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0800, [][]byte{{0xAA}})))
	suite.thenModifiedResourcesShouldBe([]int{})
}

func (suite *ManifestSuite) TestModifiedCallbackOnInsertWithChangedData() {
	suite.givenEntryWasInsertedWith(0, "id1",
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0800, [][]byte{{0xAA}})))
	suite.whenEntryIsInsertedWith(1, "id2",
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0800, [][]byte{{0xBB}})))
	suite.thenModifiedResourcesShouldBe([]int{0x0800})
}

func (suite *ManifestSuite) TestModifiedCallbackOnRemoveWithLostIDs() {
	suite.givenEntryWasInsertedWith(0, "id1",
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0800, [][]byte{{0xAA}})))
	suite.givenEntryWasInsertedWith(1, "id1",
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0900, [][]byte{{0xBB}})))
	suite.whenEntryIsRemoved(0)
	suite.thenModifiedResourcesShouldBe([]int{0x0800})
}

func (suite *ManifestSuite) TestModifiedCallbackOnRemoveIgnoresIdenticalData_RemovedBefore() {
	suite.givenEntryWasInsertedWith(0, "id1",
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0800, [][]byte{{0xAA}, {0xBB}})))
	suite.givenEntryWasInsertedWith(1, "id2",
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0800, [][]byte{{0xAA}, {0xBB}})))
	suite.whenEntryIsRemoved(0)
	suite.thenModifiedResourcesShouldBe([]int{})
}

func (suite *ManifestSuite) TestModifiedCallbackOnRemoveIgnoresIdenticalData_RemovedAfter() {
	suite.givenEntryWasInsertedWith(0, "id1",
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0800, [][]byte{{0xAA}, {0xBB}})))
	suite.givenEntryWasInsertedWith(1, "id2",
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0800, [][]byte{{0xAA}, {0xBB}})))
	suite.whenEntryIsRemoved(1)
	suite.thenModifiedResourcesShouldBe([]int{})
}

func (suite *ManifestSuite) TestModifiedCallbackOnRemoveWithChangedData() {
	suite.givenEntryWasInsertedWith(0, "id1",
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0800, [][]byte{{0xAA}})))
	suite.givenEntryWasInsertedWith(1, "id2",
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0800, [][]byte{{0xBB}})))
	suite.whenEntryIsRemoved(1)
	suite.thenModifiedResourcesShouldBe([]int{0x0800})
}

func (suite *ManifestSuite) TestModifiedCallbackOnReplaceWithChangedData() {
	suite.givenEntryWasInsertedWith(0, "id1",
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0800, [][]byte{{0xAA}})))
	suite.whenEntryIsReplacedWith(0, "id2",
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0800, [][]byte{{0xBB}})))
	suite.thenModifiedResourcesShouldBe([]int{0x0800})
}

func (suite *ManifestSuite) TestModifiedCallbackWithEmptyListOnIdenticalData() {
	suite.givenEntryWasInsertedWith(0, "id1",
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0800, [][]byte{{0xAA}})))
	suite.whenEntryIsReplacedWith(0, "id2",
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0800, [][]byte{{0xAA}})))
	suite.thenModifiedResourcesShouldBe([]int{})
}

func (suite *ManifestSuite) TestModifiedCallbackWithChangedIDsOnReplace() {
	suite.givenEntryWasInsertedWith(0, "id1",
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0800, [][]byte{{0xAA}})))
	suite.whenEntryIsReplacedWith(0, "id2",
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0900, [][]byte{{0xAA}})))
	suite.thenModifiedResourcesShouldBe([]int{0x0800, 0x0900})
}

func (suite *ManifestSuite) TestModifiedCallbackWithChangedIDsOnMove() {
	suite.givenEntryWasInsertedWith(0, "id1",
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0800, [][]byte{{0xAA}})))
	suite.givenEntryWasInsertedWith(1, "id2",
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0800, [][]byte{{0xBB}})))
	suite.givenEntryWasInsertedWith(2, "id3",
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0900, [][]byte{{0x11}})))
	suite.whenEntryIsMoved(2, 0)
	suite.thenModifiedResourcesShouldBe([]int{0x0800})
}

func (suite *ManifestSuite) onManifestModified(modifiedIDs []resource.ID, failedIDs []resource.ID) {
	suite.lastModifiedIDs = modifiedIDs
	suite.lastFailedIDs = failedIDs
}

func (suite *ManifestSuite) givenEntryWasInserted(at int, id string) {
	suite.whenEntryIsInserted(at, id)
}

func (suite *ManifestSuite) givenEntryWasInsertedWith(at int, id string, res ...resource.LocalizedResources) {
	suite.whenEntryIsInsertedWith(at, id, res...)
	suite.lastModifiedIDs = nil
	suite.lastFailedIDs = nil
}

func (suite *ManifestSuite) whenEntryIsInsertedWith(at int, id string, res ...resource.LocalizedResources) {
	err := suite.manifest.InsertEntry(at, suite.anEntryWithResources(id, res...))
	require.Nil(suite.T(), err, fmt.Sprintf("No error expected inserting entry at %d", at))
}

func (suite *ManifestSuite) whenEntriesAreInsertedWith(at int, idBase int, resList ...resource.LocalizedResources) {
	entries := make([]*world.ManifestEntry, len(resList))
	for index, res := range resList {
		entries[index] = suite.anEntryWithResources(fmt.Sprintf("id%d", idBase+index), res)
	}
	err := suite.manifest.InsertEntry(at, entries...)
	require.Nil(suite.T(), err, fmt.Sprintf("No error expected inserting entry at %d", at))
}

func (suite *ManifestSuite) whenEntryIsInserted(at int, id string) {
	err := suite.manifest.InsertEntry(at, suite.aSimpleEntry(id))
	require.Nil(suite.T(), err, fmt.Sprintf("No error expected inserting entry at %d", at))
}

func (suite *ManifestSuite) whenEntriesAreInserted(at int, ids ...string) {
	entries := make([]*world.ManifestEntry, len(ids))
	for index, id := range ids {
		entries[index] = suite.aSimpleEntry(id)
	}
	err := suite.manifest.InsertEntry(at, entries...)
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

func (suite *ManifestSuite) whenEntryIsReplacedWith(at int, id string, res ...resource.LocalizedResources) {
	err := suite.manifest.ReplaceEntry(at, suite.anEntryWithResources(id, res...))
	require.Nil(suite.T(), err, fmt.Sprintf("No error expected replacing entry at %d", at))
}

func (suite *ManifestSuite) whenEntryIsMoved(to int, from int) {
	err := suite.manifest.MoveEntry(to, from)
	require.Nil(suite.T(), err, fmt.Sprintf("No error expected moving entry to %d from %d", to, from))
}

func (suite *ManifestSuite) whenResourcesAreQueriedFor(lang resource.Language) {
	selector := suite.manifest.LocalizedResources(lang)
	suite.selector = &selector
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

func (suite *ManifestSuite) thenResourcesCanBeSelected(id int) {
	view, err := suite.selector.Select(resource.ID(id))
	assert.Nil(suite.T(), err, "No error expected")
	assert.NotNil(suite.T(), view, "View expected")
}

func (suite *ManifestSuite) aSimpleEntry(id string) *world.ManifestEntry {
	return &world.ManifestEntry{
		ID: id,
	}
}

func (suite *ManifestSuite) someLocalizedResources(lang resource.Language, modifiers ...func(*resource.Store)) resource.LocalizedResources {
	store := resource.NewProviderBackedStore(resource.NullProvider())
	for _, modifier := range modifiers {
		modifier(store)
	}
	return resource.LocalizedResources{
		ID:       "unnamed",
		Language: lang,
		Provider: store,
	}
}

func (suite *ManifestSuite) anEntryWithResources(id string, res ...resource.LocalizedResources) *world.ManifestEntry {
	return &world.ManifestEntry{
		ID:        id,
		Resources: res,
	}
}

func (suite *ManifestSuite) storing(id int, data [][]byte) func(*resource.Store) {
	return func(store *resource.Store) {
		store.Put(resource.ID(id), resource.Resource{
			Blocks: resource.BlocksFrom(data),
		}.ToView())
	}
}

func (suite *ManifestSuite) thenModifiedResourcesShouldBe(expected []int) {
	identified := make([]resource.ID, len(expected))
	for index, id := range expected {
		identified[index] = resource.ID(id)
	}
	suite.sortIDs(identified)
	suite.sortIDs(suite.lastModifiedIDs)

	assert.Equal(suite.T(), identified, suite.lastModifiedIDs, "Modified IDs don't match")
}

func (suite *ManifestSuite) sortIDs(ids []resource.ID) {
	sort.Slice(ids, func(a, b int) bool { return ids[a] < ids[b] })
}
