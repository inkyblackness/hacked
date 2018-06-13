package text_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/resource"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type LineCacheSuite struct {
	suite.Suite

	localizedResources resource.LocalizedResourcesList

	cp       text.Codepage
	instance *text.LineCache
}

func TestResourceLineCacheSuite(t *testing.T) {
	suite.Run(t, new(LineCacheSuite))
}

func (suite *LineCacheSuite) SetupTest() {
	suite.cp = text.DefaultCodepage()
	suite.instance = text.NewLineCache(suite.cp, suite)
}

func (suite *LineCacheSuite) TestLineReturnsValueIfOK() {
	suite.whenResourcesAre(
		suite.someLocalizedResources(resource.LangGerman,
			suite.storing(0x1000, "test")))
	suite.thenLineShouldReturn("test", resource.KeyOf(0x1000, resource.LangGerman, 0))
}

func (suite *LineCacheSuite) TestLineReturnsErrorIfResourceNotExisting() {
	suite.whenResourcesAre()
	suite.thenLineShouldReturnError(resource.KeyOf(0x1000, resource.LangGerman, 0))
}

func (suite *LineCacheSuite) TestLineReturnsValueIfBlockNotExisting() {
	suite.whenResourcesAre(
		suite.someLocalizedResources(resource.LangGerman,
			suite.storing(0x1000, "zero", "one")))
	suite.thenLineShouldReturnError(resource.KeyOf(0x1000, resource.LangGerman, 2))
}

func (suite *LineCacheSuite) TestLineReturnsCachedValueIfPreviouslyRetrieved() {
	key := resource.KeyOf(0x1000, resource.LangGerman, 0)
	suite.givenResourcesAre(
		suite.someLocalizedResources(resource.LangGerman,
			suite.storing(0x1000, "test")))
	suite.givenLineWasRetrieved(key)
	suite.whenResourcesAre()
	suite.thenLineShouldReturn("test", key)
}

func (suite *LineCacheSuite) TestLineTriesToReloadWhenCacheInvalidated() {
	key := resource.KeyOf(0x1000, resource.LangGerman, 0)
	suite.givenResourcesAre(
		suite.someLocalizedResources(resource.LangGerman,
			suite.storing(0x1000, "test")))
	suite.givenLineWasRetrieved(key)
	suite.givenResourcesAre()
	suite.whenCacheResourcesAreInvalidated(0x1000)
	suite.thenLineShouldReturnError(key)
}

func (suite *LineCacheSuite) TestInvalidationConsidersOnlyAffectedIDs() {
	key := resource.KeyOf(0x1000, resource.LangGerman, 0)
	suite.givenResourcesAre(
		suite.someLocalizedResources(resource.LangGerman,
			suite.storing(0x1000, "test"), suite.storing(0x1001, "other")))
	suite.givenLineWasRetrieved(key)
	suite.givenResourcesAre()
	suite.whenCacheResourcesAreInvalidated(0x1001)
	suite.thenLineShouldReturn("test", key)
}

func (suite *LineCacheSuite) TestLineReturnsErrorIfResourceIsNotAText() {
	suite.whenResourcesAre(suite.someLocalizedResources(resource.LangDefault,
		suite.storingNonText(0x1000)))
	suite.thenLineShouldReturnError(resource.KeyOf(0x1000, resource.LangDefault, 0))
}

func (suite *LineCacheSuite) givenResourcesAre(resources ...resource.LocalizedResources) {
	suite.localizedResources = resources
}

func (suite *LineCacheSuite) givenLineWasRetrieved(key resource.Key) {
	suite.instance.Line(key)
}

func (suite *LineCacheSuite) whenResourcesAre(resources ...resource.LocalizedResources) {
	suite.localizedResources = resources
}

func (suite *LineCacheSuite) whenCacheResourcesAreInvalidated(ids ...resource.ID) {
	suite.instance.InvalidateResources(ids)
}

func (suite *LineCacheSuite) thenLineShouldReturn(expected string, key resource.Key) {
	line, err := suite.instance.Line(key)
	require.Nil(suite.T(), err, "No error expected")
	assert.Equal(suite.T(), expected, line, "Lines don't match")
}

func (suite *LineCacheSuite) thenLineShouldReturnError(key resource.Key) {
	_, err := suite.instance.Line(key)
	require.NotNil(suite.T(), err, "Error expected")
}

func (suite *LineCacheSuite) someLocalizedResources(lang resource.Language, modifiers ...func(*resource.Store)) resource.LocalizedResources {
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

func (suite *LineCacheSuite) storing(id int, lines ...string) func(*resource.Store) {
	blocks := make([][]byte, len(lines))
	for i, line := range lines {
		blocks[i] = suite.cp.Encode(line)
	}
	return func(store *resource.Store) {
		store.Put(resource.ID(id), &resource.Resource{
			ContentType:   resource.Text,
			Compound:      len(blocks) != 1,
			BlockProvider: resource.MemoryBlockProvider(blocks),
		})
	}
}

func (suite *LineCacheSuite) storingNonText(id int) func(*resource.Store) {
	return func(store *resource.Store) {
		store.Put(resource.ID(id), &resource.Resource{
			ContentType:   resource.Sound,
			Compound:      false,
			BlockProvider: resource.MemoryBlockProvider([][]byte{{}}),
		})
	}
}

func (suite *LineCacheSuite) LocalizedResources(lang resource.Language) resource.Selector {
	return resource.Selector{
		From: suite.localizedResources,
		Lang: lang,
	}
}
