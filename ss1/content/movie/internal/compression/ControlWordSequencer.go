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
	// BitstreamIndexLimit specifies the highest value the sequencer may use to store offset values in the bitstream.
	// If this value is 0, then the default of 0xFFF (12 bits) is used.
	BitstreamIndexLimit uint32

	ops map[TileColorOp]uint32
}

// Add extends the list of requested coloring operations with the given entry.
// The offset of the entry must be a value less than or equal ControlWordParamLimit, otherwise the function returns an error.
func (seq *ControlWordSequencer) Add(op TileColorOp) error {
	if op.Offset > ControlWordParamLimit {
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
		opA := sortedOps[a]
		opB := sortedOps[b]
		countA := seq.ops[opA]
		countB := seq.ops[opB]
		if countA != countB {
			return countA > countB
		}
		if opA.Offset != opB.Offset {
			return opA.Offset < opB.Offset
		}
		if opA.Type != opB.Type {
			return opA.Type < opB.Type
		}
		return false
	})

	bitstreamOffset := 12
	bitstreamIndexLimit := seq.BitstreamIndexLimit
	if bitstreamIndexLimit == 0 {
		bitstreamIndexLimit = 0xFFF
	}
	result.opPaths = make(map[TileColorOp]nestedTileColorOp)
	var nestedParent *nestedTileColorOp
	relOffset := uint32(0)
	for _, op := range sortedOps {
		wordCount := uint32(len(result.words))
		if wordCount == bitstreamIndexLimit {
			bitstreamOffset = 4
			result.words = append(result.words, LongOffsetOf(wordCount+1))
			nestedParent = &nestedTileColorOp{relOffsetBits: 12, relOffset: wordCount}
			relOffset = 0
		}
		if (wordCount > bitstreamIndexLimit) && (relOffset == 15) {
			result.words = append(result.words, LongOffsetOf(wordCount+1))
			nestedParent = &nestedTileColorOp{parent: nestedParent, relOffsetBits: 4, relOffset: 15}
			relOffset = 0
		}

		result.opPaths[op] = nestedTileColorOp{parent: nestedParent, relOffsetBits: uint(bitstreamOffset), relOffset: relOffset}
		result.words = append(result.words, ControlWordOf(bitstreamOffset, op.Type, op.Offset))
		relOffset++
	}

	return result, nil
}

type nestedTileColorOp struct {
	parent        *nestedTileColorOp
	relOffsetBits uint
	relOffset     uint32
}

// ControlWordSequence is a finalized set of control words to reproduce a list of
// tile coloring operations. Based on this sequence, a bitstream can be created based
// on a selection of such coloring operations (i.e., per frame).
type ControlWordSequence struct {
	words   []ControlWord
	opPaths map[TileColorOp]nestedTileColorOp
}

// ControlWords returns the list of low-level control words of the sequence.
func (seq ControlWordSequence) ControlWords() []ControlWord {
	return seq.words
}

// BitstreamFor returns the bitstream to reproduce the provided list of coloring operations
// from this sequence.
func (seq ControlWordSequence) BitstreamFor(ops []TileColorOp) ([]byte, error) {
	var writer BitstreamWriter
	for _, op := range ops {
		nested := seq.opPaths[op]
		path := []nestedTileColorOp{nested}
		for nested.parent != nil {
			nested = *nested.parent
			path = append(path, nested)
		}
		for index := len(path) - 1; index >= 0; index-- {
			nested := path[index]
			writer.Write(nested.relOffsetBits, nested.relOffset)
		}
	}
	return writer.Buffer(), nil
}
