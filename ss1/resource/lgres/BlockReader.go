package lgres

import "io"

type blockFunc func(index int) (io.Reader, error)

type blockReader struct {
	blockCount int
	blockFunc  blockFunc
}

// BlockCount returns the number of available blocks.
func (reader *blockReader) BlockCount() int {
	return reader.blockCount
}

// Block returns the reader for the identified block.
// Each call returns a new reader instance.
// Data provided by this reader is always uncompressed.
func (reader *blockReader) Block(index int) (io.Reader, error) {
	return reader.blockFunc(index)
}
