package text_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/resource"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CacheSuite struct {
	suite.Suite

	localizedResources resource.LocalizedResourcesList

	cp       text.Codepage
	instance *text.Cache
}

func TestCacheSuite(t *testing.T) {
	suite.Run(t, new(CacheSuite))
}

func (suite *CacheSuite) SetupTest() {
	suite.cp = text.DefaultCodepage()
	suite.instance = nil
}

func (suite *CacheSuite) TestTextReturnsValueIfOKForLineCache() {
	suite.givenALineCache()
	suite.whenResourcesAre(
		suite.someLocalizedResources(resource.LangGerman,
			suite.storing(0x1000, "test")))
	suite.thenTextShouldReturn("test", resource.KeyOf(0x1000, resource.LangGerman, 0))
}

func (suite *CacheSuite) TestTextReturnsValueIfOKForPageCache() {
	suite.givenAPageCache()
	suite.whenResourcesAre(
		suite.someLocalizedResources(resource.LangGerman,
			suite.storing(0x1000, "this is", " separated", " over lines")))
	suite.thenTextShouldReturn("this is separated over lines", resource.KeyOf(0x1000, resource.LangGerman, 0))
}

func (suite *CacheSuite) TestTextReturnsErrorIfResourceNotExistingForLineCache() {
	suite.givenALineCache()
	suite.whenResourcesAre()
	suite.thenTextShouldReturnError(resource.KeyOf(0x1000, resource.LangGerman, 0))
}

func (suite *CacheSuite) TestTextReturnsErrorIfResourceNotExistingForPageCache() {
	suite.givenAPageCache()
	suite.whenResourcesAre()
	suite.thenTextShouldReturnError(resource.KeyOf(0x1000, resource.LangGerman, 0))
}

func (suite *CacheSuite) TestTextReturnsErrorIfBlockNotExistingForLineCache() {
	suite.givenALineCache()
	suite.whenResourcesAre(
		suite.someLocalizedResources(resource.LangGerman,
			suite.storing(0x1000, "zero", "one")))
	suite.thenTextShouldReturnError(resource.KeyOf(0x1000, resource.LangGerman, 2))
}

func (suite *CacheSuite) TestTextReturnsCachedValueIfPreviouslyRetrieved() {
	key := resource.KeyOf(0x1000, resource.LangGerman, 0)
	suite.givenALineCache()
	suite.givenResourcesAre(
		suite.someLocalizedResources(resource.LangGerman,
			suite.storing(0x1000, "test")))
	suite.givenTextWasRetrieved(key)
	suite.whenResourcesAre()
	suite.thenTextShouldReturn("test", key)
}

func (suite *CacheSuite) TestTextTriesToReloadWhenCacheInvalidated() {
	key := resource.KeyOf(0x1000, resource.LangGerman, 0)
	suite.givenALineCache()
	suite.givenResourcesAre(
		suite.someLocalizedResources(resource.LangGerman,
			suite.storing(0x1000, "test")))
	suite.givenTextWasRetrieved(key)
	suite.givenResourcesAre()
	suite.whenCacheResourcesAreInvalidated(0x1000)
	suite.thenTextShouldReturnError(key)
}

func (suite *CacheSuite) TestInvalidationConsidersOnlyAffectedIDs() {
	suite.givenALineCache()
	key := resource.KeyOf(0x1000, resource.LangGerman, 0)
	suite.givenResourcesAre(
		suite.someLocalizedResources(resource.LangGerman,
			suite.storing(0x1000, "test"), suite.storing(0x1001, "other")))
	suite.givenTextWasRetrieved(key)
	suite.givenResourcesAre()
	suite.whenCacheResourcesAreInvalidated(0x1001)
	suite.thenTextShouldReturn("test", key)
}

func (suite *CacheSuite) TestTextReturnsErrorIfResourceIsNotATextForLineCache() {
	suite.givenALineCache()
	suite.whenResourcesAre(suite.someLocalizedResources(resource.LangDefault,
		suite.storingNonText(0x1000)))
	suite.thenTextShouldReturnError(resource.KeyOf(0x1000, resource.LangDefault, 0))
}

func (suite *CacheSuite) TestTextReturnsErrorIfResourceIsNotATextForPageCache() {
	suite.givenAPageCache()
	suite.whenResourcesAre(suite.someLocalizedResources(resource.LangDefault,
		suite.storingNonText(0x1000)))
	suite.thenTextShouldReturnError(resource.KeyOf(0x1000, resource.LangDefault, 0))
}

func (suite *CacheSuite) givenALineCache() {
	suite.instance = text.NewLineCache(suite.cp, suite)
}

func (suite *CacheSuite) givenAPageCache() {
	suite.instance = text.NewPageCache(suite.cp, suite)
}

func (suite *CacheSuite) givenResourcesAre(resources ...resource.LocalizedResources) {
	suite.localizedResources = resources
}

func (suite *CacheSuite) givenTextWasRetrieved(key resource.Key) {
	_, err := suite.instance.Text(key)
	assert.Nil(suite.T(), err, "no error expected")
}

func (suite *CacheSuite) whenResourcesAre(resources ...resource.LocalizedResources) {
	suite.localizedResources = resources
}

func (suite *CacheSuite) whenCacheResourcesAreInvalidated(ids ...resource.ID) {
	suite.instance.InvalidateResources(ids)
}

func (suite *CacheSuite) thenTextShouldReturn(expected string, key resource.Key) {
	result, err := suite.instance.Text(key)
	require.Nil(suite.T(), err, "No error expected")
	assert.Equal(suite.T(), expected, result, "Texts don't match")
}

func (suite *CacheSuite) thenTextShouldReturnError(key resource.Key) {
	_, err := suite.instance.Text(key)
	require.NotNil(suite.T(), err, "Error expected")
}

func (suite *CacheSuite) someLocalizedResources(lang resource.Language, modifiers ...func(*resource.Store)) resource.LocalizedResources {
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

func (suite *CacheSuite) storing(id int, lines ...string) func(*resource.Store) {
	data := make([][]byte, len(lines))
	for i, line := range lines {
		data[i] = suite.cp.Encode(line)
	}
	return func(store *resource.Store) {
		store.Put(resource.ID(id), &resource.Resource{
			ContentType:   resource.Text,
			Compound:      len(data) != 1,
			BlockProvider: resource.BlocksFrom(data),
		})
	}
}

func (suite *CacheSuite) storingNonText(id int) func(*resource.Store) {
	return func(store *resource.Store) {
		store.Put(resource.ID(id), &resource.Resource{
			ContentType:   resource.Sound,
			Compound:      false,
			BlockProvider: resource.BlocksFrom([][]byte{{}}),
		})
	}
}

func (suite *CacheSuite) LocalizedResources(lang resource.Language) resource.Selector {
	return resource.Selector{
		From: suite.localizedResources,
		Lang: lang,
	}
}
