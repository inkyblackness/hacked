package lgres

import (
	"io"

	"github.com/inkyblackness/hacked/ss1/resource/lgres/compression"
	"github.com/inkyblackness/hacked/ss1/serial"
)

// CompoundResourceWriter writes a resource with zero, one, or more blocks.
// Multiple blocks can be created and then written concurrently. Only when the
// resource is finished, the blocks are finalized.
type CompoundResourceWriter struct {
	target *serial.PositioningEncoder

	compressed      bool
	dataPaddingSize int
	blockStores     []*serial.ByteStore
	blockWriter     []*BlockWriter
}

// CreateBlock provides a new, dedicated writer for a new block.
func (writer *CompoundResourceWriter) CreateBlock() *BlockWriter {
	store := serial.NewByteStore()
	blockWriter := &BlockWriter{target: serial.NewEncoder(store), finisher: func() {}}
	writer.blockStores = append(writer.blockStores, store)
	writer.blockWriter = append(writer.blockWriter, blockWriter)
	return blockWriter
}

func (writer *CompoundResourceWriter) finish() (length uint32) {
	var unpackedSize uint32
	blockCount := len(writer.blockStores)
	writer.target.Code(uint16(blockCount))
	offset := 2 + (blockCount+1)*4 + writer.dataPaddingSize
	for index, store := range writer.blockStores {
		unpackedSize += writer.blockWriter[index].finish()
		writer.target.Code(uint32(offset))
		offset += len(store.Data())
	}
	writer.target.Code(uint32(offset))
	unpackedSize += writer.target.CurPos() + uint32(writer.dataPaddingSize)

	writer.writeBlocks()

	return unpackedSize
}

func (writer *CompoundResourceWriter) writeBlocks() {
	var targetWriter io.Writer = writer.target
	targetFinisher := func() {}
	if writer.compressed {
		compressor := compression.NewCompressor(targetWriter)
		targetWriter = compressor
		targetFinisher = func() { compressor.Close() } // nolint: errcheck
	}
	for i := 0; i < writer.dataPaddingSize; i++ {
		targetWriter.Write([]byte{0x00}) // nolint: errcheck
	}
	for _, store := range writer.blockStores {
		targetWriter.Write(store.Data()) // nolint: errcheck
	}
	targetFinisher()
}
