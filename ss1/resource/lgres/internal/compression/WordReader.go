package compression

import "github.com/inkyblackness/hacked/ss1/serial"

// WordReader provides word instances from a serialized stream.
type WordReader struct {
	coder serial.Coder

	buffer              byte
	bufferBitsAvailable uint
}

// NewWordReader returns a new instance.
func NewWordReader(coder serial.Coder) *WordReader {
	reader := &WordReader{coder: coder, buffer: 0x00, bufferBitsAvailable: 0}

	return reader
}

// Read returns the next word from the stream.
func (reader *WordReader) Read() (value Word) {
	remaining := bitsPerWord

	for remaining > reader.bufferBitsAvailable {
		value = Word((uint(value) << reader.bufferBitsAvailable) | uint(reader.buffer))
		remaining -= reader.bufferBitsAvailable

		reader.bufferByte()
	}

	value = Word((uint(value) << remaining) | uint(reader.buffer)>>(8-remaining))
	reader.buffer &= 1<<(8-remaining) - 1
	reader.bufferBitsAvailable -= remaining

	return
}

func (reader *WordReader) bufferByte() {
	reader.coder.Code(&reader.buffer)
	reader.bufferBitsAvailable = 8
}
