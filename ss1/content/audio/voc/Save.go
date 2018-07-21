package voc

import (
	"bytes"
	"encoding/binary"
	"io"
)

// Save encodes the provided samples into the given writer
func Save(writer io.Writer, sampleRate float32, samples []byte) {
	writeHeader(writer)
	writeBasicSoundData(writer, sampleRate, samples)
	writeEndOfFile(writer)
}

func writeHeader(writer io.Writer) {
	version := baseVersion

	writer.Write(bytes.NewBufferString(fileHeader).Bytes())
	binary.Write(writer, binary.LittleEndian, standardHeaderSize)
	binary.Write(writer, binary.LittleEndian, baseVersion)
	binary.Write(writer, binary.LittleEndian, uint16(uint16(^version)+versionCheckValue))
}

func writeBlockHeader(writer io.Writer, block blockType, dataBytes int) {
	writer.Write([]byte{byte(block), byte(dataBytes), byte(dataBytes >> 8), byte(dataBytes >> 16)})
}

func writeBasicSoundData(writer io.Writer, sampleRate float32, samples []byte) {
	sampleType := byte(0)

	writeBlockHeader(writer, soundData, len(samples)+2)

	writer.Write([]byte{sampleRateToDivisor(sampleRate), sampleType})
	writer.Write(samples)
}

func writeEndOfFile(writer io.Writer) {
	writer.Write([]byte{byte(terminator)})
}
