package resource

import (
	"bytes"
	"fmt"
	"io"
)

// MemoryBlockProvider is a block provider backed by data in memory.
type MemoryBlockProvider [][]byte

// BlockCount returns the number of available blocks.
func (provider MemoryBlockProvider) BlockCount() int {
	return len(provider)
}

// Block returns the reader for the identified block.
// Each call returns a new reader instance.
func (provider MemoryBlockProvider) Block(index int) (io.Reader, error) {
	available := len(provider)
	if (index < 0) || (index >= available) {
		return nil, fmt.Errorf("block index wrong: %v/%v", index, available)
	}
	return bytes.NewBuffer(provider[index]), nil
}
