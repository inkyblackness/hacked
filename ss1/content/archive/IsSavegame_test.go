package archive_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/inkyblackness/hacked/ss1/content/archive"
	"github.com/inkyblackness/hacked/ss1/resource"
)

func TestIsSavegameTrueForActualSavegame(t *testing.T) {
	stateData := make([]byte, archive.GameStateSize)
	stateData[0x001C+2] = 0x10
	res := resource.Resource{
		Properties: resource.Properties{
			Compressed:  false,
			ContentType: resource.Archive,
			Compound:    false,
		},
		Blocks: resource.BlocksFrom([][]byte{stateData}),
	}

	result := archive.IsSavegame(res)
	assert.True(t, result)
}

func TestIsSavegameFalseForWrongResourceContent(t *testing.T) {
	res := resource.Resource{
		Properties: resource.Properties{
			Compressed:  false,
			ContentType: resource.Archive,
			Compound:    true,
		},
		Blocks: resource.BlocksFrom([][]byte{}),
	}

	result := archive.IsSavegame(res)
	assert.False(t, result)
}

func TestIsSavegameFalseForTooShortData(t *testing.T) {
	res := resource.Resource{
		Properties: resource.Properties{
			Compressed:  false,
			ContentType: resource.Archive,
			Compound:    true,
		},
		Blocks: resource.BlocksFrom([][]byte{make([]byte, 0x10)}),
	}

	result := archive.IsSavegame(res)
	assert.False(t, result)
}

func TestIsSavegameFalseForZeroData(t *testing.T) {
	stateData := make([]byte, archive.GameStateSize)
	res := resource.Resource{
		Properties: resource.Properties{
			Compressed:  false,
			ContentType: resource.Archive,
			Compound:    true,
		},
		Blocks: resource.BlocksFrom([][]byte{stateData}),
	}

	result := archive.IsSavegame(res)
	assert.False(t, result)
}
