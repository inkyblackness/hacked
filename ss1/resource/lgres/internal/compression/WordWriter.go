package compression

import "github.com/inkyblackness/hacked/ss1/serial"

// WordWriter serializes a stream of word entries.
type WordWriter struct {
	coder serial.Coder

	bufferedBits uint
	scratch      uint32

	outBuffer [1024]byte
	outUsed   int
}

// NewWordWriter returns a new instance.
func NewWordWriter(coder serial.Coder) *WordWriter {
	writer := &WordWriter{coder: coder, bufferedBits: 0, scratch: 0}

	return writer
}

// Close finishes the stream of words.
func (writer *WordWriter) Close() {
	writer.Write(EndOfStream)
	if writer.bufferedBits > 0 {
		writer.writeByte(byte(writer.scratch >> 16))
	}
	writer.writeByte(byte(0x00))
	writer.flushBuffer()
}

// Write adds the given word to the stream.
func (writer *WordWriter) Write(value Word) {
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

func (writer *WordWriter) writeByte(value byte) {
	writer.outBuffer[writer.outUsed] = value
	writer.outUsed++
	if writer.outUsed >= len(writer.outBuffer) {
		writer.flushBuffer()
	}
}

func (writer *WordWriter) flushBuffer() {
	writer.coder.Code(writer.outBuffer[:writer.outUsed])
	writer.outUsed = 0
}
