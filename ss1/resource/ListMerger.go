package resource

import (
	"io"
)

type listMerger struct {
	list List
}

// Compound indicates whether the resource is a compound one, which can contain zero, one, or more blocks.
// Simple resources always contain exactly one block.
func (view listMerger) Compound() bool {
	return view.list[0].Compound
}

// ContentType specifies how the data in the resource shall be interpreted.
func (view listMerger) ContentType() ContentType {
	return view.list[0].ContentType
}

// Compressed indicates whether the data should be compressed when serialized.
func (view listMerger) Compressed() bool {
	return view.list[0].Compressed
}

// BlockCount returns the maximum count of all resources.
func (view listMerger) BlockCount() (max int) {
	for _, layer := range view.list {
		count := layer.BlockCount()
		if max < count {
			max = count
		}
	}
	return
}

// Block returns a reader for the identified block in the resource.
// The view returns the block from the first resource that has a non-empty block.
// The resources are checked from last-to-first.
func (view listMerger) Block(index int) (reader io.Reader, err error) {
	for layer := len(view.list) - 1; (layer >= 0) && (reader == nil); layer-- {
		tempReader, tempErr := view.list[layer].Block(index)
		if tempErr == nil {
			var buf [1]byte
			read, _ := tempReader.Read(buf[:]) // nolint
			if read == 1 {
				reader, err = view.list[layer].Block(index)
			}
		}
	}
	return
}
