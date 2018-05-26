package compression

import "github.com/inkyblackness/hacked/ss1/serial"

type wordWriter struct {
	coder serial.Coder

	bufferedBits uint
	scratch      uint32

	outBuffer [1024]byte
	outUsed   int
}

func newWordWriter(coder serial.Coder) *wordWriter {
	writer := &wordWriter{coder: coder, bufferedBits: 0, scratch: 0}

	return writer
}

func (writer *wordWriter) close() {
	writer.write(endOfStream)
	if writer.bufferedBits > 0 {
		writer.writeByte(byte(writer.scratch >> 16))
	}
	writer.writeByte(byte(0x00))
	writer.flushBuffer()
}

func (writer *wordWriter) write(value word) {
	writer.scratch |= uint32(value) << ((16 - bitsPerWord) + (8 - writer.bufferedBits))
	writer.bufferedBits += bitsPerWord
	writer.writeByte(byte(writer.scratch >> 16))
	writer.scratch <<= 8
	writer.bufferedBits -= 8

	if writer.bufferedBits >= 8 {
		writer.writeByte(byte(writer.scratch >> 16))
		writer.scratch <<= 8
		writer.bufferedBits -= 8
	}
}

func (writer *wordWriter) writeByte(value byte) {
	writer.outBuffer[writer.outUsed] = value
	writer.outUsed++
	if writer.outUsed >= len(writer.outBuffer) {
		writer.flushBuffer()
	}
}

func (writer *wordWriter) flushBuffer() {
	writer.coder.Code(writer.outBuffer[:writer.outUsed])
	writer.outUsed = 0
}
