package bitmap_test

import (
	"bytes"
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/serial"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	basePaletteResourceID = 0x02BC
)

type PaletteCacheSuite struct {
	suite.Suite

	localizedResources resource.LocalizedResourcesList

	instance *bitmap.PaletteCache
}

func TestPaletteCacheSuite(t *testing.T) {
	suite.Run(t, new(PaletteCacheSuite))
}

func (suite *PaletteCacheSuite) SetupTest() {
	suite.instance = nil
}

func (suite *PaletteCacheSuite) TestPaletteReturnsValueIfOK() {
	pal := suite.somePalette(0)
	suite.givenAnInstance()
	suite.whenResourcesAre(
		suite.someLocalizedResources(
			suite.storing(0, pal)))
	suite.thenPaletteShouldReturn(pal, 0)
}

func (suite *PaletteCacheSuite) TestPaletteReturnsErrorIfResourceNotExisting() {
	suite.givenAnInstance()
	suite.whenResourcesAre()
	suite.thenPaletteShouldReturnError(0)
}

func (suite *PaletteCacheSuite) TestTextTriesToReloadWhenCacheInvalidated() {
	suite.givenAnInstance()
	suite.givenResourcesAre(
		suite.someLocalizedResources(
			suite.storing(0, suite.somePalette(0))))
	suite.givenPaletteWasRetrieved(0)
	suite.givenResourcesAre()
	suite.whenCacheResourcesAreInvalidated(basePaletteResourceID)
	suite.thenPaletteShouldReturnError(0)
}

func (suite *PaletteCacheSuite) TestInvalidationConsidersOnlyAffectedIDs() {
	suite.givenAnInstance()
	suite.givenResourcesAre(
		suite.someLocalizedResources(
			suite.storing(0, suite.somePalette(0)), suite.storing(1, suite.somePalette(1))))
	suite.givenPaletteWasRetrieved(0)
	suite.givenResourcesAre()
	suite.whenCacheResourcesAreInvalidated(0x02BD)
	suite.thenPaletteShouldReturn(suite.somePalette(0), 0)
}

func (suite *PaletteCacheSuite) TestPaletteReturnsErrorIfResourceIsNotAPalette() {
	suite.givenAnInstance()
	suite.whenResourcesAre(suite.someLocalizedResources(suite.storingNonPalette(0)))
	suite.thenPaletteShouldReturnError(0)
}

func (suite *PaletteCacheSuite) givenAnInstance() {
	suite.instance = bitmap.NewPaletteCache(suite)
}

func (suite *PaletteCacheSuite) givenResourcesAre(resources ...resource.LocalizedResources) {
	suite.localizedResources = resources
}

func (suite *PaletteCacheSuite) givenPaletteWasRetrieved(index int) {
	_, err := suite.instance.Palette(suite.keyed(index))
	require.Nil(suite.T(), err, "No error expected setting up")
}

func (suite *PaletteCacheSuite) whenResourcesAre(resources ...resource.LocalizedResources) {
	suite.localizedResources = resources
}

func (suite *PaletteCacheSuite) whenCacheResourcesAreInvalidated(ids ...resource.ID) {
	suite.instance.InvalidateResources(ids)
}

func (suite *PaletteCacheSuite) thenPaletteShouldReturn(expected bitmap.Palette, index int) {
	result, err := suite.instance.Palette(suite.keyed(index))
	require.Nil(suite.T(), err, "No error expected")
	assert.Equal(suite.T(), expected, result, "Palettes don't match")
}

func (suite *PaletteCacheSuite) thenPaletteShouldReturnError(index int) {
	_, err := suite.instance.Palette(suite.keyed(index))
	require.NotNil(suite.T(), err, "Error expected")
}

func (suite PaletteCacheSuite) keyed(index int) resource.Key {
	return resource.KeyOf(basePaletteResourceID, resource.LangAny, index)
}

func (suite *PaletteCacheSuite) someLocalizedResources(modifiers ...func(*resource.Store)) resource.LocalizedResources {
	store := resource.NewStore()
	for _, modifier := range modifiers {
		modifier(store)
	}
	return resource.LocalizedResources{
		ID:       "unnamed",
		Language: resource.LangAny,
		Viewer:   store,
	}
}

func (suite *PaletteCacheSuite) storing(id int, palette bitmap.Palette) func(*resource.Store) {
	buf := bytes.NewBuffer(nil)
	encoder := serial.NewEncoder(buf)
	encoder.Code(palette)
	return func(store *resource.Store) {
		store.Put(resource.ID(basePaletteResourceID).Plus(id), resource.Resource{
			Properties: resource.Properties{
				ContentType: resource.Palette,
				Compound:    false,
			},
			Blocks: resource.BlocksFrom([][]byte{buf.Bytes()}),
		})
	}
}

func (suite *PaletteCacheSuite) somePalette(seed uint8) bitmap.Palette {
	var pal bitmap.Palette
	for i := 0; i < len(pal); i++ {
		pal[i].Red = seed + uint8(i)
		pal[i].Green = seed + uint8(255-i)
		pal[i].Blue = seed + uint8(i*2)
	}
	return pal
}

func (suite *PaletteCacheSuite) storingNonPalette(id int) func(*resource.Store) {
	return func(store *resource.Store) {
		store.Put(resource.ID(basePaletteResourceID).Plus(id), resource.Resource{
			Properties: resource.Properties{
				ContentType: resource.Sound,
				Compound:    false,
			},
			Blocks: resource.BlocksFrom([][]byte{{}}),
		})
	}
}

func (suite *PaletteCacheSuite) LocalizedResources(lang resource.Language) resource.Selector {
	return resource.Selector{
		From: suite.localizedResources,
		Lang: lang,
	}
}
