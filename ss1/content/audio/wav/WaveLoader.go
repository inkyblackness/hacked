package wav

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

func l8FromL8(input []byte) []byte {
	return input
}

func l8FromL16(input []byte) []byte {
	samples := len(input) / 2
	output := make([]byte, samples)

	for i := 0; i < samples; i++ {
		output[i] = byte(int(input[i*2+1]) + 0x80)
	}

	return output
}

type waveLoader struct {
	dataRead   bool
	formatRead bool
	samples    []byte
	sampleRate float32

	reader io.Reader
	err    error

	dataConverter func([]byte) []byte
}

func (loader *waveLoader) load(reader io.Reader) {
	loader.reader = reader
	loader.dataConverter = l8FromL8
	loader.loadRiff()
}

func (loader *waveLoader) read(data interface{}) bool {
	if loader.err == nil {
		loader.err = binary.Read(loader.reader, binary.LittleEndian, data)
	}
	return loader.err == nil
}

func (loader *waveLoader) readBytes(size uint32) (data []byte) {
	data = make([]byte, int(size))
	_, loader.err = loader.reader.Read(data)
	return
}

func (loader *waveLoader) loadChunk(handler func(riffChunkType, uint32)) {
	var tag riffChunkTag

	if loader.read(&tag) {
		handler(tag.ChunkType, tag.Size)
	}
}

func (loader *waveLoader) loadRiff() {
	loader.loadChunk(loader.handleRiff)
}

func (loader *waveLoader) handleRiff(chunkType riffChunkType, size uint32) {
	if chunkType == riffChunkTypeRiff {
		var contentType riffContentType
		if loader.read(&contentType) && (contentType == riffContentTypeWave) {
			loader.loadWave()
		} else if loader.err == nil {
			loader.err = errNotASupportedWave
		}
	} else {
		loader.err = errNotASupportedWave
	}
}

func (loader *waveLoader) loadWave() {
	for !loader.isDone() {
		loader.loadFormatOrData()
	}
	loader.samples = loader.dataConverter(loader.samples)
}

func (loader *waveLoader) loadFormatOrData() {
	loader.loadChunk(loader.handleFormatOrData)
}

func (loader *waveLoader) handleFormatOrData(chunkType riffChunkType, size uint32) {
	if chunkType == riffChunkTypeFmt {
		loader.loadFormat(size)
	} else if chunkType == riffChunkTypeData {
		loader.loadData(size)
	}
}

func (loader *waveLoader) loadFormat(size uint32) {
	headerData := loader.readBytes(size)
	headerReader := bytes.NewReader(headerData)
	var header formatHeader

	loader.formatRead = true
	binary.Read(headerReader, binary.LittleEndian, &header.base)
	binary.Read(headerReader, binary.LittleEndian, &header.extension.BitsPerSample)
	loader.sampleRate = float32(header.base.SamplesPerSec)

	if header.extension.BitsPerSample == 16 {
		loader.dataConverter = l8FromL16
	}

	if (header.base.FormatType != waveFormatTypePcm) ||
		(header.base.Channels != 1) ||
		((header.extension.BitsPerSample != 8) && (header.extension.BitsPerSample != 16)) {
		loader.err = fmt.Errorf("unsupported WAVE format")
	}
}

func (loader *waveLoader) loadData(size uint32) {
	loader.dataRead = true
	loader.samples = loader.readBytes(size)
}

func (loader *waveLoader) isDone() bool {
	return (loader.err != nil) || (loader.dataRead && loader.formatRead)
}
