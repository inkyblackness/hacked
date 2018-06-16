package resource

import (
	"bytes"
	"fmt"
	"io"
)

// Resource provides meta information as well as access to its contained blocks.
type Resource struct {
	// Compound tells whether the resource should be serialized with a directory.
	// Compound resources can have zero, one, or more blocks.
	// Simple resources always have exactly one block.
	Compound bool

	// ContentType describes how the block data shall be interpreted.
	ContentType ContentType

	// Compressed tells whether the data shall be serialized in compressed form.
	Compressed bool

	// BlockProvider is the keeper of original block data.
	// This provider will be referred to if no other data was explicitly set.
	BlockProvider BlockProvider

	blockLimit int
	blocks     map[int][]byte
}

// BlockCount returns the number of available blocks in the resource.
// Simple resources will always have exactly one block.
func (res Resource) BlockCount() (count int) {
	count = res.providerBlockCount()
	if count < res.blockLimit {
		count = res.blockLimit
	}
	return
}

func (res Resource) providerBlockCount() (count int) {
	if res.BlockProvider != nil {
		count = res.BlockProvider.BlockCount()
	}
	return
}

// Block returns the reader for the identified block.
// Each call returns a new reader instance.
// Data provided by this reader is always uncompressed.
func (res Resource) Block(index int) (io.Reader, error) {
	if res.blocks != nil {
		data, set := res.blocks[index]
		if set {
			return bytes.NewReader(data), nil
		} else if res.providerBlockCount() <= index {
			return bytes.NewReader(nil), nil
		}
	}
	if res.BlockProvider == nil {
		return nil, fmt.Errorf("no blocks available")
	}
	return res.BlockProvider.Block(index)
}

// SetBlock registers new data for a block.
// For any block set this way, the block provider of this resource will no longer be queried.
func (res *Resource) SetBlock(index int, data []byte) {
	if index < 0 {
		panic(fmt.Errorf("index must be a non-negative value"))
	}
	res.ensureBlockMap()
	res.blocks[index] = data
	if res.blockLimit <= index {
		res.blockLimit = index + 1
	}
}

func (res *Resource) ensureBlockMap() {
	if res.blocks == nil {
		res.blocks = make(map[int][]byte)
	}
}

// ToView returns a view of this resource.
func (res Resource) ToView() View {
	return simpleView{&res}
}
