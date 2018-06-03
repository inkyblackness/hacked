package resource_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/resource"

	"github.com/stretchr/testify/assert"
)

func TestIsSavegameTrueForActualSavegame(t *testing.T) {
	stateData := make([]byte, 0x0100) // TODO: use proper length
	stateData[0x009C] = 0x80
	store := resource.NewProviderBackedStore(resource.NullProvider())
	store.Put(0x0FA1, &resource.Resource{
		Compressed:    false,
		ContentType:   resource.Archive,
		Compound:      false,
		BlockProvider: resource.MemoryBlockProvider([][]byte{stateData}),
	})

	result := resource.IsSavegame(store)
	assert.True(t, result)
}

func TestIsSavegameFalseForMissingStateData(t *testing.T) {
	store := resource.NewProviderBackedStore(resource.NullProvider())

	result := resource.IsSavegame(store)
	assert.False(t, result)
}

func TestIsSavegameFalseForWrongResourceContent(t *testing.T) {
	store := resource.NewProviderBackedStore(resource.NullProvider())
	store.Put(0x0FA1, &resource.Resource{
		Compressed:    false,
		ContentType:   resource.Archive,
		Compound:      true,
		BlockProvider: resource.MemoryBlockProvider([][]byte{}),
	})

	result := resource.IsSavegame(store)
	assert.False(t, result)
}

func TestIsSavegameFalseForTooShortData(t *testing.T) {
	store := resource.NewProviderBackedStore(resource.NullProvider())
	store.Put(0x0FA1, &resource.Resource{
		Compressed:    false,
		ContentType:   resource.Archive,
		Compound:      true,
		BlockProvider: resource.MemoryBlockProvider([][]byte{make([]byte, 0x10)}),
	})

	result := resource.IsSavegame(store)
	assert.False(t, result)
}

func TestIsSavegameFalseForZeroData(t *testing.T) {
	stateData := make([]byte, 0x0100) // TODO: use proper length
	store := resource.NewProviderBackedStore(resource.NullProvider())
	store.Put(0x0FA1, &resource.Resource{
		Compressed:    false,
		ContentType:   resource.Archive,
		Compound:      true,
		BlockProvider: resource.MemoryBlockProvider([][]byte{stateData}),
	})

	result := resource.IsSavegame(store)
	assert.False(t, result)
}
