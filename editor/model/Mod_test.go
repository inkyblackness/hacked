package model_test

import (
	"fmt"
	"io/ioutil"
	"sort"
	"testing"

	"github.com/inkyblackness/hacked/editor/model"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ModSuite struct {
	suite.Suite
	mod *model.Mod

	selector *resource.Selector

	lastModifiedIDs []resource.ID
	lastFailedIDs   []resource.ID
}

func TestModSuite(t *testing.T) {
	suite.Run(t, new(ModSuite))
}

func (suite *ModSuite) SetupTest() {
	suite.mod = model.NewMod(suite.onResourcesModified)

	suite.selector = nil

	suite.lastModifiedIDs = nil
	suite.lastFailedIDs = nil
}

func (suite *ModSuite) onResourcesModified(modifiedIDs []resource.ID, failedIDs []resource.ID) {
	suite.lastModifiedIDs = modifiedIDs
	suite.lastFailedIDs = failedIDs
}

func (suite *ModSuite) TestResourcesCanBeRetrievedFromTheWorld() {
	suite.givenWorldHas(
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0800, [][]byte{{0xAA}})))
	suite.whenResourcesAreQueriedFor(resource.LangAny)
	suite.thenResourcesCanBeSelected(0x0800)
}

func (suite *ModSuite) TestResourcesCanBeModified() {
	suite.givenWorldHas(
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0800, [][]byte{{0xAA}})))
	suite.whenModifyingBy(func(trans *model.ModTransaction) {
		trans.SetResourceBlock(resource.LangAny, 0x0800, 0, []byte{0xBB})
	})
	suite.thenResourceBlockShouldBe(resource.LangAny, 0x0800, 0, []byte{0xBB})
}

func (suite *ModSuite) TestResourcesCanBeExtended() {
	suite.givenWorldHas(
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0800, [][]byte{{0xAA}})))
	suite.whenModifyingBy(func(trans *model.ModTransaction) {
		trans.SetResourceBlock(resource.LangAny, 0x0800, 2, []byte{0xBB})
	})
	suite.thenResourceBlockShouldBe(resource.LangAny, 0x0800, 2, []byte{0xBB})
}

func (suite *ModSuite) TestResourceMetaCanBeChanged() {
	suite.givenWorldHas(
		suite.someLocalizedResources(resource.LangGerman,
			func(store *resource.Store) {
				store.Put(0x1000, &resource.Resource{
					Compound:      false,
					ContentType:   resource.Sound,
					Compressed:    false,
					BlockProvider: resource.MemoryBlockProvider(nil),
				})
			}),
		suite.someLocalizedResources(resource.LangFrench,
			func(store *resource.Store) {
				store.Put(0x1000, &resource.Resource{
					Compound:      false,
					ContentType:   resource.Palette,
					Compressed:    false,
					BlockProvider: resource.MemoryBlockProvider(nil),
				})
			}))
	suite.whenModifyingBy(func(trans *model.ModTransaction) {
		trans.SetResource(0x1000, true, resource.Movie, true)
	})
	suite.thenResourceMetaShouldBe(resource.LangFrench, 0x1000, true, resource.Movie, true)
	suite.thenResourceMetaShouldBe(resource.LangGerman, 0x1000, true, resource.Movie, true)
}

func (suite *ModSuite) TestResourcesCanBeRemoved() {
	suite.givenWorldHas(
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0800, [][]byte{{0xAA}})))
	suite.givenModifiedBy(func(trans *model.ModTransaction) {
		trans.SetResourceBlock(resource.LangAny, 0x0800, 0, []byte{0xBB})
	})
	suite.whenModifyingBy(func(trans *model.ModTransaction) {
		trans.DelResource(resource.LangAny, 0x0800)
	})
	suite.thenResourceBlockShouldBe(resource.LangAny, 0x0800, 0, []byte{0xAA})
}

func (suite *ModSuite) TestMetaModificationIsNotified() {
	suite.givenWorldHas(
		suite.someLocalizedResources(resource.LangGerman,
			func(store *resource.Store) {
				store.Put(0x1000, &resource.Resource{
					Compound:      false,
					ContentType:   resource.Sound,
					Compressed:    false,
					BlockProvider: resource.MemoryBlockProvider(nil),
				})
			}))
	suite.whenModifyingBy(func(trans *model.ModTransaction) {
		trans.SetResource(0x1000, true, resource.Movie, true)
	})
	suite.thenModifiedResourcesShouldBe(0x1000)
}

func (suite *ModSuite) TestAdditionsAreNotified() {
	suite.givenWorldHas(
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0800, [][]byte{{0xAA}})))
	suite.whenModifyingBy(func(trans *model.ModTransaction) {
		trans.SetResourceBlock(resource.LangAny, 0x0800, 2, []byte{0xBB})
	})
	suite.thenModifiedResourcesShouldBe(0x0800)
}

func (suite *ModSuite) TestModificationsToIdenticalDataAreNotNotified() {
	suite.givenWorldHas(
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0800, [][]byte{{0xAA}})))
	suite.whenModifyingBy(func(trans *model.ModTransaction) {
		trans.SetResourceBlock(resource.LangAny, 0x0800, 0, []byte{0xAA})
	})
	suite.thenModifiedResourcesShouldBe()
}

func (suite *ModSuite) TestDeletionIsNotified() {
	suite.givenWorldHas(
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x0800, [][]byte{{0xAA}})))
	suite.givenModifiedBy(func(trans *model.ModTransaction) {
		trans.SetResourceBlock(resource.LangAny, 0x0800, 0, []byte{0xBB})
	})
	suite.whenModifyingBy(func(trans *model.ModTransaction) {
		trans.DelResource(resource.LangAny, 0x0800)
	})
	suite.thenModifiedResourcesShouldBe(0x0800)
}

func (suite *ModSuite) givenWorldHas(res ...resource.LocalizedResources) {
	suite.whenWorldIsExtendedWith(res...)
	suite.lastModifiedIDs = nil
	suite.lastFailedIDs = nil
}

func (suite *ModSuite) whenWorldIsExtendedWith(res ...resource.LocalizedResources) {
	manifest := suite.mod.World()
	at := manifest.EntryCount()
	id := fmt.Sprintf("entry-%d", at)
	err := suite.mod.World().InsertEntry(at, suite.anEntryWithResources(id, res...))
	require.Nil(suite.T(), err, fmt.Sprintf("No error expected inserting entry at %d", at))
}

func (suite *ModSuite) whenResourcesAreQueriedFor(lang resource.Language) {
	selector := suite.mod.LocalizedResources(lang)
	suite.selector = &selector
}

func (suite *ModSuite) givenModifiedBy(modifier func(*model.ModTransaction)) {
	suite.mod.Modify(modifier)

	suite.lastModifiedIDs = nil
	suite.lastFailedIDs = nil
}

func (suite *ModSuite) whenModifyingBy(modifier func(*model.ModTransaction)) {
	suite.mod.Modify(modifier)
}

func (suite *ModSuite) thenResourceBlockShouldBe(lang resource.Language, id int, blockIndex int, expected []byte) {
	view, viewErr := suite.mod.LocalizedResources(lang).Select(resource.ID(id))
	require.Nil(suite.T(), viewErr, "No error expected selecting resource")
	reader, blockErr := view.Block(blockIndex)
	require.Nil(suite.T(), blockErr, "No error expected retrieving block")
	data, dataErr := ioutil.ReadAll(reader)
	require.Nil(suite.T(), dataErr, "No error expected reading data")
	assert.Equal(suite.T(), expected, data, "Data mismatch in block")
}

func (suite *ModSuite) thenResourceMetaShouldBe(lang resource.Language, id int,
	compound bool, contentType resource.ContentType, compressed bool) {
	view, viewErr := suite.mod.LocalizedResources(lang).Select(resource.ID(id))
	require.Nil(suite.T(), viewErr, "No error expected selecting resource")

	key := fmt.Sprintf("lang %v, res %v", lang, resource.ID(id))
	assert.Equal(suite.T(), compound, view.Compound(), "Compound property does not match for "+key)
	assert.Equal(suite.T(), contentType, view.ContentType(), "ContentType property does not match for "+key)
	assert.Equal(suite.T(), compressed, view.Compressed(), "Compressed property does not match for "+key)
}

func (suite *ModSuite) thenResourcesCanBeSelected(id int) {
	view, err := suite.selector.Select(resource.ID(id))
	assert.Nil(suite.T(), err, "No error expected")
	assert.NotNil(suite.T(), view, "View expected")
}

func (suite *ModSuite) thenModifiedResourcesShouldBe(expected ...int) {
	identified := make([]resource.ID, len(expected))
	for index, id := range expected {
		identified[index] = resource.ID(id)
	}
	suite.sortIDs(identified)
	suite.sortIDs(suite.lastModifiedIDs)

	assert.Equal(suite.T(), identified, suite.lastModifiedIDs, "Modified IDs don't match")
}

func (suite *ModSuite) sortIDs(ids []resource.ID) {
	sort.Slice(ids, func(a, b int) bool { return ids[a] < ids[b] })
}

func (suite *ModSuite) someLocalizedResources(lang resource.Language, modifiers ...func(*resource.Store)) resource.LocalizedResources {
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

func (suite *ModSuite) anEntryWithResources(id string, res ...resource.LocalizedResources) *world.ManifestEntry {
	return &world.ManifestEntry{
		ID:        id,
		Resources: res,
	}
}

func (suite *ModSuite) storing(id int, blocks [][]byte) func(*resource.Store) {
	return func(store *resource.Store) {
		store.Put(resource.ID(id), &resource.Resource{
			BlockProvider: resource.MemoryBlockProvider(blocks),
		})
	}
}
