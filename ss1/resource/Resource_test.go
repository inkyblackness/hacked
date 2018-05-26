package resource_test

import (
	"io/ioutil"
	"testing"

	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/stretchr/testify/assert"
)

func verifyBlockContent(t *testing.T, provider resource.BlockProvider, index int, expected []byte) {
	reader, err := provider.Block(index)
	assert.Nil(t, err, "No error expected for index %d", index)
	assert.NotNil(t, reader, "Reader expected for index %d", index)
	if reader != nil {
		data, dataErr := ioutil.ReadAll(reader)
		assert.Nil(t, dataErr, "No error expected reading data from index %d", index)
		assert.Equal(t, expected, data, "Proper data expected from index %d", index)
	}
}

func verifyBlockError(t *testing.T, provider resource.BlockProvider, index int) {
	_, err := provider.Block(index)
	assert.NotNil(t, err, "Error expected for index %d", index)
}

func TestResourceRefersToBlockProviderByDefault(t *testing.T) {
	res := &resource.Resource{BlockProvider: resource.MemoryBlockProvider([][]byte{{0x01}, {0x02, 0x02}})}

	assert.Equal(t, 2, res.BlockCount())
	verifyBlockContent(t, res, 0, []byte{0x01})
	verifyBlockContent(t, res, 1, []byte{0x02, 0x02})
}

func TestResourceBlockReturnsErrorOnInvalidIndex(t *testing.T) {
	res := &resource.Resource{BlockProvider: resource.MemoryBlockProvider(nil)}

	verifyBlockError(t, res, -1)
	verifyBlockError(t, res, 0)
	verifyBlockError(t, res, 1)
	verifyBlockError(t, res, 2)
}

func TestResourceBlockReturnsErrorForDefaultObject(t *testing.T) {
	var res resource.Resource

	assert.Equal(t, 0, res.BlockCount())
	verifyBlockError(t, res, -1)
	verifyBlockError(t, res, 0)
	verifyBlockError(t, res, 1)
}

func TestResourceCanBeExtendedWithBlocks(t *testing.T) {
	var res resource.Resource

	res.SetBlock(0, []byte{0x10})
	res.SetBlock(2, []byte{0x20, 0x20})
	assert.Equal(t, 3, res.BlockCount())
	verifyBlockContent(t, res, 0, []byte{0x10})
	verifyBlockContent(t, res, 1, []byte{})
	verifyBlockContent(t, res, 2, []byte{0x20, 0x20})
}

func TestResourceDefaultsToProviderWhenNoExtensionOverridesIt(t *testing.T) {
	res := &resource.Resource{BlockProvider: resource.MemoryBlockProvider([][]byte{{0x01}, {0x02, 0x02}})}

	res.SetBlock(0, []byte{0xA0})
	assert.Equal(t, 2, res.BlockCount())
	verifyBlockContent(t, res, 0, []byte{0xA0})
	verifyBlockContent(t, res, 1, []byte{0x02, 0x02})
}

func TestResourcePanicsForNegativeBlockIndex(t *testing.T) {
	var res resource.Resource

	assert.Panics(t, func() { res.SetBlock(-1, []byte{}) }, "Panic expected")
}
