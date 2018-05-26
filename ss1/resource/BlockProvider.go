package resource

import "io"

// BlockProvider are capable of returning block data of a resource.
type BlockProvider interface {
	// BlockCount returns the number of available blocks in the resource.
	// Simple resources will always have exactly one block.
	BlockCount() int

	// Block returns the reader for the identified block.
	// Each call returns a new reader instance.
	// Data provided by this reader is always uncompressed.
	Block(index int) (io.Reader, error)
}
