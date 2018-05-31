package resource

import (
	"io"
)

// ResourceView is a read-only view on a selected resource.
type ResourceView interface {

	// Compound tells whether the resource should be serialized with a directory.
	// Compound resources can have zero, one, or more blocks.
	// Simple resources always have exactly one block.
	Compound() bool

	// ContentType describes how the block data shall be interpreted.
	ContentType() ContentType

	// Compressed tells whether the data shall be serialized in compressed form.
	Compressed() bool

	// BlockCount returns the number of available blocks in the resource.
	// Simple resources will always have exactly one block.
	BlockCount() int

	// Block returns the reader for the identified block.
	// Each call returns a new reader instance.
	// Data provided by this reader is always uncompressed.
	Block(index int) (io.Reader, error)
}
