package voc

import (
	"bytes"
	"encoding/binary"
	"io"
)

// Save encodes the provided samples into the given writer.
func Save(writer io.Writer, sampleRate float32, samples []byte) error {
	err := writeHeader(writer)
	if err != nil {
		return err
	}
	err = writeBasicSoundData(writer, sampleRate, samples)
	if err != nil {
		return err
	}
	err = writeEndOfFile(writer)
	return err
}

func writeHeader(writer io.Writer) error {
	version := baseVersion

	_, err := writer.Write(bytes.NewBufferString(fileHeader).Bytes())
	if err != nil {
		return err
	}
	err = binary.Write(writer, binary.LittleEndian, standardHeaderSize)
	if err != nil {
		return err
	}
	err = binary.Write(writer, binary.LittleEndian, baseVersion)
	if err != nil {
		return err
	}
	return binary.Write(writer, binary.LittleEndian, (^version)+versionCheckValue)
}

func writeBlockHeader(writer io.Writer, block blockType, dataBytes int) error {
	_, err := writer.Write([]byte{byte(block), byte(dataBytes), byte(dataBytes >> 8), byte(dataBytes >> 16)})
	return err
}

func writeBasicSoundData(writer io.Writer, sampleRate float32, samples []byte) error {
	sampleType := byte(0)

	err := writeBlockHeader(writer, soundData, len(samples)+2)
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte{sampleRateToDivisor(sampleRate), sampleType})
	if err != nil {
		return err
	}
	_, err = writer.Write(samples)
	return err
}

func writeEndOfFile(writer io.Writer) error {
	_, err := writer.Write([]byte{byte(terminator)})
	return err
}
