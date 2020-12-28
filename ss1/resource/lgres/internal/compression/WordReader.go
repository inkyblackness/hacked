package compression

import "github.com/inkyblackness/hacked/ss1/serial"

type wordReader struct {
	coder serial.Coder

	buffer              byte
	bufferBitsAvailable uint
}

func newWordReader(coder serial.Coder) *wordReader {
	reader := &wordReader{coder: coder, buffer: 0x00, bufferBitsAvailable: 0}

	return reader
}

func (reader *wordReader) read() (value word) {
	remaining := bitsPerWord

	for remaining > reader.bufferBitsAvailable {
		value = word((uint(value) << reader.bufferBitsAvailable) | uint(reader.buffer))
		remaining -= reader.bufferBitsAvailable

		reader.bufferByte()
	}

	value = word((uint(value) << remaining) | uint(reader.buffer)>>(8-remaining))
	reader.buffer &= 1<<(8-remaining) - 1
	reader.bufferBitsAvailable -= remaining

	return
}

func (reader *wordReader) bufferByte() {
	reader.coder.Code(&reader.buffer)
	reader.bufferBitsAvailable = 8
}
