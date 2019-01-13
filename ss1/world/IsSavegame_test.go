package world_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/archive"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/hacked/ss1/world/ids"

	"github.com/stretchr/testify/assert"
)

func TestIsSavegameTrueForActualSavegame(t *testing.T) {
	stateData := make([]byte, archive.GameStateSize)
	stateData[0x009C] = 0x80
	store := resource.NewProviderBackedStore(resource.NullProvider())
	store.Put(ids.GameState, resource.Resource{
		Properties: resource.Properties{
			Compressed:  false,
			ContentType: resource.Archive,
			Compound:    false,
		},
		Blocks: resource.BlocksFrom([][]byte{stateData}),
	}.ToView())

	result := world.IsSavegame(store)
	assert.True(t, result)
}

func TestIsSavegameFalseForMissingStateData(t *testing.T) {
	store := resource.NewProviderBackedStore(resource.NullProvider())

	result := world.IsSavegame(store)
	assert.False(t, result)
}

func TestIsSavegameFalseForWrongResourceContent(t *testing.T) {
	store := resource.NewProviderBackedStore(resource.NullProvider())
	store.Put(ids.GameState, resource.Resource{
		Properties: resource.Properties{
			Compressed:  false,
			ContentType: resource.Archive,
			Compound:    true,
		},
		Blocks: resource.BlocksFrom([][]byte{}),
	}.ToView())

	result := world.IsSavegame(store)
	assert.False(t, result)
}

func TestIsSavegameFalseForTooShortData(t *testing.T) {
	store := resource.NewProviderBackedStore(resource.NullProvider())
	store.Put(ids.GameState, resource.Resource{
		Properties: resource.Properties{
			Compressed:  false,
			ContentType: resource.Archive,
			Compound:    true,
		},
		Blocks: resource.BlocksFrom([][]byte{make([]byte, 0x10)}),
	}.ToView())

	result := world.IsSavegame(store)
	assert.False(t, result)
}

func TestIsSavegameFalseForZeroData(t *testing.T) {
	stateData := make([]byte, archive.GameStateSize)
	store := resource.NewProviderBackedStore(resource.NullProvider())
	store.Put(ids.GameState, resource.Resource{
		Properties: resource.Properties{
			Compressed:  false,
			ContentType: resource.Archive,
			Compound:    true,
		},
		Blocks: resource.BlocksFrom([][]byte{stateData}),
	}.ToView())

	result := world.IsSavegame(store)
	assert.False(t, result)
}
