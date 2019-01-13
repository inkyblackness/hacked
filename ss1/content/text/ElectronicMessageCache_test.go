package text_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/resource"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ElectronicMessageCacheSuite struct {
	suite.Suite

	localizedResources resource.LocalizedResourcesList

	cp       text.Codepage
	instance *text.ElectronicMessageCache
}

func TestElectronicMessageCacheSuite(t *testing.T) {
	suite.Run(t, new(ElectronicMessageCacheSuite))
}

func (suite *ElectronicMessageCacheSuite) SetupTest() {
	suite.cp = text.DefaultCodepage()
	suite.instance = nil
}

func (suite *ElectronicMessageCacheSuite) TestMessageReturnsValueIfOK() {
	suite.givenACache()
	suite.whenResourcesAre(
		suite.someLocalizedResources(resource.LangGerman,
			suite.storing(0x1000, func(msg text.ElectronicMessage) text.ElectronicMessage {
				msg.VerboseText = "123"
				msg.TerseText = "terse"
				return msg
			})))
	suite.thenMessageShouldReturn(resource.KeyOf(0x1000, resource.LangGerman, 0), func(msg text.ElectronicMessage) {
		assert.Equal(suite.T(), "123", msg.VerboseText)
		assert.Equal(suite.T(), "terse", msg.TerseText)
	})
}

func (suite *ElectronicMessageCacheSuite) TestTextReturnsErrorIfResourceNotExisting() {
	suite.givenACache()
	suite.whenResourcesAre()
	suite.thenMessageShouldReturnError(resource.KeyOf(0x1000, resource.LangGerman, 0))
}

func (suite *ElectronicMessageCacheSuite) TestTextReturnsCachedValueIfPreviouslyRetrieved() {
	key := resource.KeyOf(0x1000, resource.LangGerman, 0)
	suite.givenACache()
	suite.givenResourcesAre(
		suite.someLocalizedResources(resource.LangGerman,
			suite.storing(0x1000, func(msg text.ElectronicMessage) text.ElectronicMessage {
				msg.VerboseText = "older text"
				msg.TerseText = "terse text"
				return msg
			})))
	suite.givenMessageWasRetrieved(key)
	suite.whenResourcesAre()
	suite.thenMessageShouldReturn(key, func(msg text.ElectronicMessage) {
		assert.Equal(suite.T(), "older text", msg.VerboseText)
		assert.Equal(suite.T(), "terse text", msg.TerseText)
	})
}

func (suite *ElectronicMessageCacheSuite) TestTextTriesToReloadWhenCacheInvalidated() {
	key := resource.KeyOf(0x1000, resource.LangGerman, 0)
	suite.givenACache()
	suite.givenResourcesAre(
		suite.someLocalizedResources(resource.LangGerman,
			suite.storing(0x1000, func(msg text.ElectronicMessage) text.ElectronicMessage {
				msg.VerboseText = "12"
				msg.TerseText = "20"
				return msg
			})))
	suite.givenMessageWasRetrieved(key)
	suite.givenResourcesAre()
	suite.whenCacheResourcesAreInvalidated(0x1000)
	suite.thenMessageShouldReturnError(key)
}

func (suite *ElectronicMessageCacheSuite) TestInvalidationConsidersOnlyAffectedIDs() {
	suite.givenACache()
	key := resource.KeyOf(0x1000, resource.LangGerman, 0)
	suite.givenResourcesAre(
		suite.someLocalizedResources(resource.LangGerman,
			suite.storing(0x1000, func(msg text.ElectronicMessage) text.ElectronicMessage {
				msg.VerboseText = "A"
				return msg
			}),
			suite.storing(0x1001, func(msg text.ElectronicMessage) text.ElectronicMessage {
				msg.VerboseText = "B"
				return msg
			})))
	suite.givenMessageWasRetrieved(key)
	suite.givenResourcesAre()
	suite.whenCacheResourcesAreInvalidated(0x1001)
	suite.thenMessageShouldReturn(key, func(msg text.ElectronicMessage) {
		assert.Equal(suite.T(), "A", msg.VerboseText)
	})
}

func (suite *ElectronicMessageCacheSuite) TestTextReturnsErrorIfResourceIsNotATextForLineCache() {
	suite.givenACache()
	suite.whenResourcesAre(suite.someLocalizedResources(resource.LangDefault,
		suite.storingNonText(0x1000)))
	suite.thenMessageShouldReturnError(resource.KeyOf(0x1000, resource.LangDefault, 0))
}

func (suite *ElectronicMessageCacheSuite) givenACache() {
	suite.instance = text.NewElectronicMessageCache(suite.cp, suite)
}

func (suite *ElectronicMessageCacheSuite) givenResourcesAre(resources ...resource.LocalizedResources) {
	suite.localizedResources = resources
}

func (suite *ElectronicMessageCacheSuite) givenMessageWasRetrieved(key resource.Key) {
	_, err := suite.instance.Message(key)
	assert.Nil(suite.T(), err, "no error expected")
}

func (suite *ElectronicMessageCacheSuite) whenResourcesAre(resources ...resource.LocalizedResources) {
	suite.localizedResources = resources
}

func (suite *ElectronicMessageCacheSuite) whenCacheResourcesAreInvalidated(ids ...resource.ID) {
	suite.instance.InvalidateResources(ids)
}

func (suite *ElectronicMessageCacheSuite) thenMessageShouldReturn(key resource.Key, validator func(msg text.ElectronicMessage)) {
	result, err := suite.instance.Message(key)
	require.Nil(suite.T(), err, "No error expected")
	validator(result)
}

func (suite *ElectronicMessageCacheSuite) thenMessageShouldReturnError(key resource.Key) {
	_, err := suite.instance.Message(key)
	require.NotNil(suite.T(), err, "Error expected")
}

func (suite *ElectronicMessageCacheSuite) someLocalizedResources(lang resource.Language, modifiers ...func(*resource.Store)) resource.LocalizedResources {
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

func (suite *ElectronicMessageCacheSuite) storing(id int, modifier func(msg text.ElectronicMessage) text.ElectronicMessage) func(*resource.Store) {
	data := modifier(text.EmptyElectronicMessage()).Encode(suite.cp)
	return func(store *resource.Store) {
		store.Put(resource.ID(id), resource.Resource{
			Properties: resource.Properties{
				ContentType: resource.Text,
				Compound:    true,
			},
			Blocks: resource.BlocksFrom(data),
		})
	}
}

func (suite *ElectronicMessageCacheSuite) storingNonText(id int) func(*resource.Store) {
	return func(store *resource.Store) {
		store.Put(resource.ID(id), resource.Resource{
			Properties: resource.Properties{
				ContentType: resource.Sound,
				Compound:    false,
			},
			Blocks: resource.BlocksFrom([][]byte{{}}),
		})
	}
}

func (suite *ElectronicMessageCacheSuite) LocalizedResources(lang resource.Language) resource.Selector {
	return resource.Selector{
		From: suite.localizedResources,
		Lang: lang,
	}
}
