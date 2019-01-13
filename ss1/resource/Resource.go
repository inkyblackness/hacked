package resource

import (
	"fmt"
	"io"
)

// Properties describe the meta information about a resource.
type Properties struct {
	// Compound tells whether the resource should be serialized with a directory.
	// Compound resources can have zero, one, or more blocks.
	// Simple resources always have exactly one block.
	Compound bool

	// ContentType describes how the block data shall be interpreted.
	ContentType ContentType

	// Compressed tells whether the data shall be serialized in compressed form.
	Compressed bool
}

// Resource provides meta information as well as access to its contained blocks.
type Resource struct {
	Properties

	// Blocks is the keeper of original block data.
	// This provider will be referred to if no other data was explicitly set.
	Blocks BlockProvider
}

// BlockCount returns the number of available blocks in the resource.
// Simple resources will always have exactly one block.
func (res Resource) BlockCount() (count int) {
	if res.Blocks != nil {
		count = res.Blocks.BlockCount()
	}
	return
}

// Block returns the reader for the identified block.
// Each call returns a new reader instance.
// Data provided by this reader is always uncompressed.
func (res Resource) Block(index int) (io.Reader, error) {
	if res.Blocks == nil {
		return nil, fmt.Errorf("no blocks available")
	}
	return res.Blocks.Block(index)
}

// ToView returns a view of this resource.
func (res Resource) ToView() View {
	return simpleView{res: &res}
}
