package compression

import (
	"errors"
	"sort"
)

// TileColorOp describes one operation how a tile should be colored.
type TileColorOp struct {
	Type   ControlType
	Offset uint32
}

// ControlWordSequencer receives a list of requested tile coloring operations
// and produces a sequence of low-level control words to reproduce the requested list.
type ControlWordSequencer struct {
	ops map[TileColorOp]uint32
}

// Add extends the list of requested coloring operations with the given entry.
// The offset of the entry must use up to 20 bits, otherwise the function returns an error.
func (seq *ControlWordSequencer) Add(op TileColorOp) error {
	if op.Offset >= 0x100000 {
		return errors.New("too high operation offset")
	}
	if seq.ops == nil {
		seq.ops = make(map[TileColorOp]uint32)
	}
	seq.ops[op]++
	return nil
}

// Sequence packs the requested list of entries into a sequence of control words.
// An error is returned if the operations would exceed the possible storage space.
func (seq ControlWordSequencer) Sequence() (ControlWordSequence, error) {
	var result ControlWordSequence
	var sortedOps []TileColorOp
	for op := range seq.ops {
		sortedOps = append(sortedOps, op)
	}
	sort.Slice(sortedOps, func(a, b int) bool {
		return seq.ops[sortedOps[a]] > seq.ops[sortedOps[b]]
	})
	for _, op := range sortedOps {
		result.words = append(result.words, ControlWordOf(12, op.Type, op.Offset))
	}

	return result, nil
}

// ControlWordSequence is a finalized set of control words to reproduce a list of
// tile coloring operations. Based on this sequence, a bitstream can be created based
// on a selection of such coloring operations (i.e., per frame).
type ControlWordSequence struct {
	words []ControlWord
}

// ControlWords returns the list of low-level control words of the sequence.
func (seq ControlWordSequence) ControlWords() []ControlWord {
	return seq.words
}

// BitstreamFor returns the bitstream to reproduce the provided list of coloring operations
// from this sequence.
func (seq ControlWordSequence) BitstreamFor(ops []TileColorOp) ([]byte, error) {
	return nil, errors.New("not implemented")
}
