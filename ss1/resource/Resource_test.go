package resource_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/inkyblackness/hacked/ss1/resource"
)

func verifyBlockContent(t *testing.T, provider resource.BlockProvider, index int, expected []byte) {
	t.Helper()
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
	t.Helper()
	_, err := provider.Block(index)
	assert.NotNil(t, err, "Error expected for index %d", index)
}

func TestResourceRefersToBlockProviderByDefault(t *testing.T) {
	res := &resource.Resource{Blocks: resource.BlocksFrom([][]byte{{0x01}, {0x02, 0x02}})}

	assert.Equal(t, 2, res.BlockCount())
	verifyBlockContent(t, res, 0, []byte{0x01})
	verifyBlockContent(t, res, 1, []byte{0x02, 0x02})
}

func TestResourceBlockReturnsErrorOnInvalidIndex(t *testing.T) {
	res := &resource.Resource{Blocks: resource.BlocksFrom(nil)}

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
