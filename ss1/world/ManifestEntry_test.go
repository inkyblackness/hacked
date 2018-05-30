package world_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/world"

	"fmt"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ManifestEntrySuite struct {
	suite.Suite
	entry  *world.ManifestEntry
	finder world.ResourceSelector
}

func TestManifestEntrySuite(t *testing.T) {
	suite.Run(t, new(ManifestEntrySuite))
}

func (suite *ManifestEntrySuite) SetupTest() {
	suite.entry = new(world.ManifestEntry)
}

func (suite *ManifestEntrySuite) TestResourceReturnsErrorIfNothingFound() {
	suite.whenViewingResource(world.LangAny)
	suite.thenResourceViewShouldNotBeAvailable(1234)
}

func (suite *ManifestEntrySuite) TestResourceReturnsViewIfMatched() {
	suite.givenLocalizedResources(world.LangAny, 1000, [][]byte{{0x01}})
	suite.whenViewingResource(world.LangAny)
	suite.thenResourceViewShouldBeAvailable(1000)
}

func (suite *ManifestEntrySuite) TestResourceConsidersResourcesWithLanguageAnyWhenLookingForSpecific() {
	suite.givenLocalizedResources(world.LangAny, 1000, [][]byte{{0x01}})
	suite.whenViewingResource(world.LangDefault)
	suite.thenResourceViewShouldBeAvailable(1000)
}

func (suite *ManifestEntrySuite) TestResourceConsidersResourcesWithSpecificLanguage() {
	suite.givenLocalizedResources(world.LangFrench, 1000, [][]byte{{0x01}})
	suite.whenViewingResource(world.LangFrench)
	suite.thenResourceViewShouldBeAvailable(1000)
}

func (suite *ManifestEntrySuite) TestResourceIgnoresResourcesWithSpecificLanguageWhenLookingForAny() {
	suite.givenLocalizedResources(world.LangFrench, 1000, [][]byte{{0x01}})
	suite.whenViewingResource(world.LangAny)
	suite.thenResourceViewShouldNotBeAvailable(1000)
}

func (suite *ManifestEntrySuite) TestResourceIgnoresResourcesWithSpecificLanguageWhenLookingForDifferentLanguage() {
	suite.givenLocalizedResources(world.LangFrench, 1000, [][]byte{{0x01}})
	suite.whenViewingResource(world.LangDefault)
	suite.thenResourceViewShouldNotBeAvailable(1000)

	suite.whenViewingResource(world.LangGerman)
	suite.thenResourceViewShouldNotBeAvailable(1000)
}

func (suite *ManifestEntrySuite) givenLocalizedResources(lang world.Language, id int, blocks [][]byte) {
	store := resource.NewProviderBackedStore(resource.NullProvider())
	store.Put(resource.ID(id), &resource.Resource{
		BlockProvider: resource.MemoryBlockProvider(blocks),
	})
	suite.entry.Resources = append(suite.entry.Resources, world.LocalizedResources{
		Language: lang,
		Provider: store,
	})
}

func (suite *ManifestEntrySuite) whenViewingResource(lang world.Language) {
	suite.finder = suite.entry.LocalizedResources(lang)
}

func (suite *ManifestEntrySuite) thenResourceViewShouldNotBeAvailable(id int) {
	view, err := suite.finder.Select(resource.ID(id))
	assert.Error(suite.T(), err, fmt.Sprintf("Error expected finding resource %v", id))
	assert.Nil(suite.T(), view, fmt.Sprintf("No view expected finding resource %v", id))
}

func (suite *ManifestEntrySuite) thenResourceViewShouldBeAvailable(id int) {
	view, err := suite.finder.Select(resource.ID(id))
	require.Nil(suite.T(), err, fmt.Sprintf("No error expected finding resource %v", id))
	assert.NotNil(suite.T(), view, fmt.Sprintf("ResourceView expected for resource %v", id))
}
