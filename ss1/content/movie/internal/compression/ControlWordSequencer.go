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

	bitstreamIndexLimit := seq.BitstreamIndexLimit
	if bitstreamIndexLimit == 0 {
		bitstreamIndexLimit = 0xFFF
	}
	maximumDirect := int(bitstreamIndexLimit / 2)
	requiredExtensions := 0

	if len(sortedOps) > maximumDirect {
		requiredExtensions = (len(sortedOps) - maximumDirect + 15) / 16
		if requiredExtensions > int(bitstreamIndexLimit)-maximumDirect {
			return result, errors.New("too many words")
		}
	}
	extensionStart := maximumDirect + requiredExtensions

	result.opPaths = make(map[TileColorOp]nestedTileColorOp)
	for opIndex, op := range sortedOps {
		if opIndex == maximumDirect {
			// write required amount of long offset entries
			for i := 0; i < requiredExtensions; i++ {
				result.words = append(result.words, LongOffsetOf(uint32(extensionStart+i*16)))
			}
		}
		if opIndex >= maximumDirect {
			offset := uint32(opIndex - maximumDirect)
			result.opPaths[op] = nestedTileColorOp{
				parent:        &nestedTileColorOp{relOffsetBits: 12, relOffset: uint32(maximumDirect) + offset/16},
				relOffsetBits: 4,
				relOffset:     offset % 16}
			result.words = append(result.words, ControlWordOf(4, op.Type, op.Offset))
		} else {
			result.opPaths[op] = nestedTileColorOp{parent: nil, relOffsetBits: 12, relOffset: uint32(opIndex)}
			result.words = append(result.words, ControlWordOf(12, op.Type, op.Offset))
		}
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
	// HTiles specifies the amount of horizontal tiles (= operations) a frame has.
	// This number is relevant for skip operations. If 0, then no skip operation-compression is done.
	HTiles uint32

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
	writeOp := func(op TileColorOp) {
		// TODO: if previous op is to be written, write repeat?
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
	pendingSkips := uint32(0)
	writePendingSkips := func() {
		for pendingSkips != 0 {
			toSkip := pendingSkips
			if toSkip >= 0x1F {
				toSkip = 0x1E
			}
			writeOp(TileColorOp{Type: CtrlSkip})
			writer.Write(5, toSkip-1)
			pendingSkips -= toSkip
		}
	}
	writeLineSkip := func() {
		if pendingSkips > 0 {
			writeOp(TileColorOp{Type: CtrlSkip})
			writer.Write(5, 0x1F)
			pendingSkips = 0
		}
	}
	for opIndex, op := range ops {
		if (pendingSkips > 0) && (uint32(opIndex)%seq.HTiles) == 0 {
			writeLineSkip()
		}
		switch {
		case (op.Type == CtrlSkip) && (seq.HTiles == 0):
			pendingSkips = 1
			writePendingSkips()
		case op.Type == CtrlSkip:
			pendingSkips++
		default:
			writePendingSkips()
			writeOp(op)
		}
	}
	if pendingSkips > 0 {
		writeLineSkip()
	}
	return writer.Buffer(), nil
}
