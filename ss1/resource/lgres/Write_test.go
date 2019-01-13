package lgres_test

import (
	"bytes"
	"testing"

	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/resource/lgres"
	"github.com/inkyblackness/hacked/ss1/serial"

	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	target := serial.NewByteStore()
	provider := resource.NewProviderBackedStore(resource.NullProvider())
	aResource := func(compressed bool, contentType resource.ContentType, compound bool, blocks resource.Blocks) resource.View {
		return resource.Resource{
			Properties: resource.Properties{
				Compressed:  compressed,
				ContentType: contentType,
				Compound:    compound,
			},
			Blocks: blocks,
		}
	}

	provider.Put(resource.ID(1), aResource(false, resource.Bitmap, false, resource.BlocksFrom([][]byte{{0x11}})))
	provider.Put(resource.ID(3), aResource(false, resource.Font, true, resource.BlocksFrom([][]byte{{0x21}, {0x22, 0x23}})))
	provider.Put(resource.ID(2), aResource(true, resource.Geometry, false, resource.BlocksFrom([][]byte{{0x31}})))
	provider.Put(resource.ID(4), aResource(true, resource.Archive, true, resource.BlocksFrom([][]byte{{0x41}, {0x42, 0x43}})))

	errWrite := lgres.Write(target, provider)
	if errWrite != nil {
		assert.Nil(t, errWrite, "no error expected writing")
	}

	reader, errReader := lgres.ReaderFrom(bytes.NewReader(target.Data()))
	if errReader != nil {
		assert.Nil(t, errReader, "no error expected reading")
	}

	assert.Equal(t, []resource.ID{resource.ID(1), resource.ID(3), resource.ID(2), resource.ID(4)}, reader.IDs())
}
