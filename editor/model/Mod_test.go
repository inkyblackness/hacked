package model_test

import (
	"fmt"
	"testing"

	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"

	"github.com/inkyblackness/hacked/editor/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
)

type ModSuite struct {
	suite.Suite
	mod *model.Mod

	selector *world.ResourceSelector

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
		suite.someLocalizedResources(world.LangAny,
			suite.storing(0x0800, [][]byte{{0xAA}})))
	suite.whenResourcesAreQueriedFor(world.LangAny)
	suite.thenResourcesCanBeSelected(0x0800)
}

/*
func (suite *ModSuite) TestResourcesCanBeModified() {
	suite.givenWorldHas(
		suite.someLocalizedResources(world.LangAny,
			suite.storing(0x0800, [][]byte{{0xAA}})))
	suite.whenModifyingBy(func(trans *model.ModTransaction) {
		trans.SetResourceBlock(world.LangAny, 0x0800, 0, []byte{0xBB})
	})
	suite.thenResourceBlockShouldBe(world.LangAny, 0x0800, 0, []byte{0xBB})
}
*/

func (suite *ModSuite) givenWorldHas(res ...world.LocalizedResources) {
	suite.whenWorldIsExtendedWith(res...)
	suite.lastModifiedIDs = nil
	suite.lastFailedIDs = nil
}

func (suite *ModSuite) whenWorldIsExtendedWith(res ...world.LocalizedResources) {
	manifest := suite.mod.World()
	at := manifest.EntryCount()
	id := fmt.Sprintf("entry-%d", at)
	err := suite.mod.World().InsertEntry(at, suite.anEntryWithResources(id, res...))
	require.Nil(suite.T(), err, fmt.Sprintf("No error expected inserting entry at %d", at))
}

func (suite *ModSuite) whenResourcesAreQueriedFor(lang world.Language) {
	selector := suite.mod.LocalizedResources(lang)
	suite.selector = &selector
}

func (suite *ModSuite) thenResourcesCanBeSelected(id int) {
	view, err := suite.selector.Select(resource.ID(id))
	assert.Nil(suite.T(), err, "No error expected")
	assert.NotNil(suite.T(), view, "View expected")
}

func (suite *ModSuite) someLocalizedResources(lang world.Language, modifiers ...func(*resource.Store)) world.LocalizedResources {
	store := resource.NewProviderBackedStore(resource.NullProvider())
	for _, modifier := range modifiers {
		modifier(store)
	}
	return world.LocalizedResources{
		ID:       "unnamed",
		Language: lang,
		Provider: store,
	}
}

func (suite *ModSuite) anEntryWithResources(id string, res ...world.LocalizedResources) *world.ManifestEntry {
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

func (suite *ModSuite) whenModifyingBy(modifier func(*model.ModTransaction)) {
	suite.mod.Modify(modifier)
}

func (suite *ModSuite) thenResourceBlockShouldBe(lang world.Language, id int, blockIndex int, expected []byte) {
	view, viewErr := suite.mod.LocalizedResources(lang).Select(resource.ID(id))
	require.Nil(suite.T(), viewErr, "No error expected selecting resource")
	reader, blockErr := view.Block(blockIndex)
	require.Nil(suite.T(), blockErr, "No error expected retrieving block")
	data, dataErr := ioutil.ReadAll(reader)
	require.Nil(suite.T(), dataErr, "No error expected reading data")
	assert.Equal(suite.T(), expected, data, "Data mismatch in block")
}
