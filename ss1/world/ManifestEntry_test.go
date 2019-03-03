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

type ManifestEntrySuite struct {
	suite.Suite
	entry  *world.ManifestEntry
	finder resource.Selector
}

func TestManifestEntrySuite(t *testing.T) {
	suite.Run(t, new(ManifestEntrySuite))
}

func (suite *ManifestEntrySuite) SetupTest() {
	suite.entry = new(world.ManifestEntry)
}

func (suite *ManifestEntrySuite) TestResourceReturnsErrorIfNothingFound() {
	suite.whenViewingResource(resource.LangAny)
	suite.thenResourceViewShouldNotBeAvailable(1234)
}

func (suite *ManifestEntrySuite) TestResourceReturnsViewIfMatched() {
	suite.givenLocalizedResources(resource.LangAny, 1000, [][]byte{{0x01}})
	suite.whenViewingResource(resource.LangAny)
	suite.thenResourceViewShouldBeAvailable(1000)
}

func (suite *ManifestEntrySuite) TestResourceConsidersResourcesWithLanguageAnyWhenLookingForSpecific() {
	suite.givenLocalizedResources(resource.LangAny, 1000, [][]byte{{0x01}})
	suite.whenViewingResource(resource.LangDefault)
	suite.thenResourceViewShouldBeAvailable(1000)
}

func (suite *ManifestEntrySuite) TestResourceConsidersResourcesWithSpecificLanguage() {
	suite.givenLocalizedResources(resource.LangFrench, 1000, [][]byte{{0x01}})
	suite.whenViewingResource(resource.LangFrench)
	suite.thenResourceViewShouldBeAvailable(1000)
}

func (suite *ManifestEntrySuite) TestResourceIgnoresResourcesWithSpecificLanguageWhenLookingForAny() {
	suite.givenLocalizedResources(resource.LangFrench, 1000, [][]byte{{0x01}})
	suite.whenViewingResource(resource.LangAny)
	suite.thenResourceViewShouldNotBeAvailable(1000)
}

func (suite *ManifestEntrySuite) TestResourceIgnoresResourcesWithSpecificLanguageWhenLookingForDifferentLanguage() {
	suite.givenLocalizedResources(resource.LangFrench, 1000, [][]byte{{0x01}})
	suite.whenViewingResource(resource.LangDefault)
	suite.thenResourceViewShouldNotBeAvailable(1000)

	suite.whenViewingResource(resource.LangGerman)
	suite.thenResourceViewShouldNotBeAvailable(1000)
}

func (suite *ManifestEntrySuite) givenLocalizedResources(lang resource.Language, id int, data [][]byte) {
	var store resource.Store
	_ = store.Put(resource.ID(id), resource.Resource{
		Blocks: resource.BlocksFrom(data),
	})
	suite.entry.Resources = append(suite.entry.Resources, resource.LocalizedResources{
		Language: lang,
		Viewer:   store,
	})
}

func (suite *ManifestEntrySuite) whenViewingResource(lang resource.Language) {
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
	assert.NotNil(suite.T(), view, fmt.Sprintf("View expected for resource %v", id))
}
