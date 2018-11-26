package bitmap_test

import (
	"bytes"
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	baseAnimationResourceID = 0x0A4C
)

type AnimationCacheSuite struct {
	suite.Suite

	localizedResources resource.LocalizedResourcesList

	instance *bitmap.AnimationCache
}

func TestAnimationCacheSuite(t *testing.T) {
	suite.Run(t, new(AnimationCacheSuite))
}

func (suite *AnimationCacheSuite) SetupTest() {
	suite.instance = nil
}

func (suite *AnimationCacheSuite) TestAnimationReturnsValueIfOK() {
	anim := suite.someAnimation(0)
	suite.givenAnInstance()
	suite.whenResourcesAre(
		suite.someLocalizedResources(
			suite.storing(0, anim)))
	suite.thenAnimationShouldReturn(anim, 0)
}

func (suite *AnimationCacheSuite) TestAnimationReturnsErrorIfResourceNotExisting() {
	suite.givenAnInstance()
	suite.whenResourcesAre()
	suite.thenAnimationShouldReturnError(0)
}

func (suite *AnimationCacheSuite) TestTextTriesToReloadWhenCacheInvalidated() {
	suite.givenAnInstance()
	suite.givenResourcesAre(
		suite.someLocalizedResources(
			suite.storing(0, suite.someAnimation(0))))
	suite.givenAnimationWasRetrieved(0)
	suite.givenResourcesAre()
	suite.whenCacheResourcesAreInvalidated(baseAnimationResourceID)
	suite.thenAnimationShouldReturnError(0)
}

func (suite *AnimationCacheSuite) TestInvalidationConsidersOnlyAffectedIDs() {
	suite.givenAnInstance()
	suite.givenResourcesAre(
		suite.someLocalizedResources(
			suite.storing(0, suite.someAnimation(0)), suite.storing(1, suite.someAnimation(1))))
	suite.givenAnimationWasRetrieved(0)
	suite.givenResourcesAre()
	suite.whenCacheResourcesAreInvalidated(0x02BD)
	suite.thenAnimationShouldReturn(suite.someAnimation(0), 0)
}

func (suite *AnimationCacheSuite) TestAnimationReturnsErrorIfResourceIsNotAnAnimation() {
	suite.givenAnInstance()
	suite.whenResourcesAre(suite.someLocalizedResources(suite.storingNonAnimation(0)))
	suite.thenAnimationShouldReturnError(0)
}

func (suite *AnimationCacheSuite) givenAnInstance() {
	suite.instance = bitmap.NewAnimationCache(suite)
}

func (suite *AnimationCacheSuite) givenResourcesAre(resources ...resource.LocalizedResources) {
	suite.localizedResources = resources
}

func (suite *AnimationCacheSuite) givenAnimationWasRetrieved(index int) {
	_, err := suite.instance.Animation(suite.keyed(index))
	require.Nil(suite.T(), err, "No error expected setting up")
}

func (suite *AnimationCacheSuite) whenResourcesAre(resources ...resource.LocalizedResources) {
	suite.localizedResources = resources
}

func (suite *AnimationCacheSuite) whenCacheResourcesAreInvalidated(ids ...resource.ID) {
	suite.instance.InvalidateResources(ids)
}

func (suite *AnimationCacheSuite) thenAnimationShouldReturn(expected bitmap.Animation, index int) {
	result, err := suite.instance.Animation(suite.keyed(index))
	require.Nil(suite.T(), err, "No error expected")
	assert.Equal(suite.T(), expected, result, "Animations don't match")
}

func (suite *AnimationCacheSuite) thenAnimationShouldReturnError(index int) {
	_, err := suite.instance.Animation(suite.keyed(index))
	require.NotNil(suite.T(), err, "Error expected")
}

func (suite AnimationCacheSuite) keyed(index int) resource.Key {
	return resource.KeyOf(baseAnimationResourceID, resource.LangAny, index)
}

func (suite *AnimationCacheSuite) someLocalizedResources(modifiers ...func(*resource.Store)) resource.LocalizedResources {
	store := resource.NewProviderBackedStore(resource.NullProvider())
	for _, modifier := range modifiers {
		modifier(store)
	}
	return resource.LocalizedResources{
		ID:       "unnamed",
		Language: resource.LangAny,
		Provider: store,
	}
}

func (suite *AnimationCacheSuite) storing(id int, anim bitmap.Animation) func(*resource.Store) {
	buf := bytes.NewBuffer(nil)
	_ = bitmap.WriteAnimation(buf, anim)
	return func(store *resource.Store) {
		store.Put(resource.ID(baseAnimationResourceID).Plus(id), &resource.Resource{
			ContentType:   resource.Animation,
			Compound:      true,
			BlockProvider: resource.MemoryBlockProvider([][]byte{buf.Bytes()}),
		})
	}
}

func (suite *AnimationCacheSuite) someAnimation(seed uint8) bitmap.Animation {
	anim := bitmap.Animation{
		Width:      int16(seed) * 4,
		Height:     int16(seed) * 3,
		IntroFlag:  uint16(seed % 2),
		ResourceID: resource.ID(seed),
		Entries:    make([]bitmap.AnimationEntry, seed+1),
	}

	for i := 0; i < len(anim.Entries); i++ {
		anim.Entries[i] = bitmap.AnimationEntry{
			FirstFrame: byte(i),
			LastFrame:  byte(i),
			FrameTime:  100,
		}
	}
	return anim
}

func (suite *AnimationCacheSuite) storingNonAnimation(id int) func(*resource.Store) {
	return func(store *resource.Store) {
		store.Put(resource.ID(baseAnimationResourceID).Plus(id), &resource.Resource{
			ContentType:   resource.Sound,
			Compound:      true,
			BlockProvider: resource.MemoryBlockProvider([][]byte{{}}),
		})
	}
}

func (suite *AnimationCacheSuite) LocalizedResources(lang resource.Language) resource.Selector {
	return resource.Selector{
		From: suite.localizedResources,
		Lang: lang,
	}
}
