package resource_test

import (
	"io/ioutil"
	"testing"

	"github.com/inkyblackness/hacked/ss1/resource"

	"github.com/stretchr/testify/assert"
)

func TestBlockCountReturnsZeroForNilInstance(t *testing.T) {
	var blocks resource.Blocks
	assert.Equal(t, 0, blocks.BlockCount())
}

func TestBlockCountReturnsLengthOfArray(t *testing.T) {
	blocks := resource.BlocksFrom(make([][]byte, 3))
	assert.Equal(t, 3, blocks.BlockCount())
}

func TestBlockReturnsArrayEntries(t *testing.T) {
	blocks := resource.BlocksFrom([][]byte{{0x01}, {0x02, 0x03}})
	verifyBlock := func(index int) {
		reader, err := blocks.Block(index)
		assert.Nil(t, err)
		assert.NotNil(t, reader)
	}
	verifyBlock(0)
	verifyBlock(1)
}

func TestBlockReturnsErrorForInvalidIndex(t *testing.T) {
	blocks := resource.BlocksFrom([][]byte{{0x01}, {0x02, 0x03}})
	verifyError := func(index int) {
		_, err := blocks.Block(index)
		assert.NotNil(t, err, "Error expected for index ", index)
	}
	verifyError(-1)
	verifyError(2)
	verifyError(3)
}

func TestBlockSetting(t *testing.T) {
	var blocks resource.Blocks

	blocks.Set(make([][]byte, 3))
	assert.Equal(t, 3, blocks.BlockCount(), "block count should have been set")

	blocks.SetBlock(1, []byte{0x01, 0x02})
	block, err := blocks.Block(1)
	assert.Nil(t, err, "block should be known")
	storedData, err := ioutil.ReadAll(block)
	assert.Nil(t, err, "data should be read")
	assert.Equal(t, []byte{0x01, 0x02}, storedData, "stored data mismatch")

	blocks.SetBlock(5, []byte{0xAA})
	assert.Equal(t, 6, blocks.BlockCount(), "block count should have been updated")
}
