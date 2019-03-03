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

// BlockRaw returns the raw byte slice stored in the identified block.
func (blocks Blocks) BlockRaw(index int) ([]byte, error) {
	available := len(blocks.data)
	if (index < 0) || (index >= available) {
		return nil, fmt.Errorf("block index wrong: %v/%v", index, available)
	}
	return blocks.data[index], nil
}

// Block returns the reader for the identified block.
// Each call returns a new reader instance.
func (blocks Blocks) Block(index int) (io.Reader, error) {
	raw, err := blocks.BlockRaw(index)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(raw), nil
}

// Set the data of all the blocks.
func (blocks *Blocks) Set(data [][]byte) {
	blocks.data = data
}

// SetBlock sets the data of the identified block.
func (blocks *Blocks) SetBlock(index int, data []byte) {
	if index < 0 {
		return
	}
	if index >= len(blocks.data) {
		newData := make([][]byte, index+1)
		copy(newData, blocks.data)
		blocks.data = newData
	}
	blocks.data[index] = data
}
