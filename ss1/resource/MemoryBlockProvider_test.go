package resource_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/resource"

	"github.com/stretchr/testify/assert"
)

func TestBlockCountReturnsZeroForNilInstance(t *testing.T) {
	var provider resource.MemoryBlockProvider
	assert.Equal(t, 0, provider.BlockCount())
}

func TestBlockCountReturnsLengthOfArray(t *testing.T) {
	var provider resource.MemoryBlockProvider = make([][]byte, 3)
	assert.Equal(t, 3, provider.BlockCount())
}

func TestBlockReturnsArrayEntries(t *testing.T) {
	var provider resource.MemoryBlockProvider = [][]byte{{0x01}, {0x02, 0x03}}
	verifyBlock := func(index int) {
		reader, err := provider.Block(index)
		assert.Nil(t, err)
		assert.NotNil(t, reader)
	}
	verifyBlock(0)
	verifyBlock(1)
}

func TestBlockReturnsErrorForInvalidIndex(t *testing.T) {
	var provider resource.MemoryBlockProvider = [][]byte{{0x01}, {0x02, 0x03}}
	verifyError := func(index int) {
		_, err := provider.Block(index)
		assert.NotNil(t, err, "Error expected for index ", index)
	}
	verifyError(-1)
	verifyError(2)
	verifyError(3)
}
