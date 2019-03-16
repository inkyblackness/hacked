package internal

import (
	"errors"

	"github.com/inkyblackness/hacked/ss1/content/movie/compression"
)

// TileColorOp describes one operation how a tile should be colored.
type TileColorOp struct {
	Type   compression.ControlType
	Offset uint32
}

// ControlWordSequencer receives a list of requested tile coloring operations
// and produces a sequence of low-level control words to reproduce the requested list.
type ControlWordSequencer struct {
}

// Add extends the list of requested coloring operations with the given entry.
// The offset of the entry must use up to 20 bits, otherwise the function returns an error.
func (seq *ControlWordSequencer) Add(op TileColorOp) error {
	return errors.New("not implemented")
}

// Sequence packs the requested list of entries into a sequence of control words.
// An error is returned if the operations would exceed the possible storage space.
func (seq *ControlWordSequencer) Sequence() (ControlWordSequence, error) {
	return ControlWordSequence{}, errors.New("not implemented")
}

// ControlWordSequence is a finalized set of control words to reproduce a list of
// tile coloring operations. Based on this sequence, a bitstream can be created based
// on a selection of such coloring operations (i.e., per frame).
type ControlWordSequence struct {
}

// ControlWords returns the list of low-level control words of the sequence.
func (seq ControlWordSequence) ControlWords() []compression.ControlWord {
	return nil
}

// BitstreamFor returns the bitstream to reproduce the provided list of coloring operations
// from this sequence.
func (seq ControlWordSequence) BitstreamFor(ops []TileColorOp) ([]byte, error) {
	return nil, errors.New("not implemented")
}
