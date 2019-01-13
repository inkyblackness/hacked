package resource

import (
	"bytes"
	"fmt"
	"io"
)

// Block is one set of bytes stored in a resource.
// The interpretation of the bytes is dependent on the content type of the resource.
type Block []byte

// Blocks is a list of blocks in memory.
type Blocks struct {
	data [][]byte
}

// BlocksFrom returns a blocks instance from given data.
func BlocksFrom(data [][]byte) Blocks {
	return Blocks{data: data}
}

// BlockCount returns the number of available blocks.
func (blocks Blocks) BlockCount() int {
	return len(blocks.data)
}

// Block returns the reader for the identified block.
// Each call returns a new reader instance.
func (blocks Blocks) Block(index int) (io.Reader, error) {
	available := len(blocks.data)
	if (index < 0) || (index >= available) {
		return nil, fmt.Errorf("block index wrong: %v/%v", index, available)
	}
	return bytes.NewBuffer(blocks.data[index]), nil
}
