package wav

import (
	"encoding/binary"
	"io"
)

// Save encodes the provided samples into the given writer
func Save(writer io.Writer, sampleRate float32, samples []byte) error {
	dataSize := uint32(len(samples))
	var fmt formatHeader
	fmtSize := fmt.size()
	contentType := riffContentTypeWave
	contentTypeSize := uint32(4)
	tagSize := uint32(8)
	tagSizes := uint32(tagSize * 2)

	fmt.base.FormatType = waveFormatTypePcm
	fmt.base.Channels = 1
	fmt.base.SamplesPerSec = uint32(sampleRate)
	fmt.base.AvgBytesPerSec = fmt.base.SamplesPerSec
	fmt.base.BlockAlign = 1
	fmt.extension.BitsPerSample = 8

	riffTag := riffChunkTag{ChunkType: riffChunkTypeRiff, Size: tagSizes + contentTypeSize + fmtSize + dataSize}

	_ = binary.Write(writer, binary.LittleEndian, &riffTag)
	_ = binary.Write(writer, binary.LittleEndian, &contentType)
	_ = binary.Write(writer, binary.LittleEndian, &riffChunkTag{ChunkType: riffChunkTypeFmt, Size: fmtSize})
	_ = binary.Write(writer, binary.LittleEndian, &fmt.base)
	_ = binary.Write(writer, binary.LittleEndian, &fmt.extension)
	_ = binary.Write(writer, binary.LittleEndian, &riffChunkTag{ChunkType: riffChunkTypeData, Size: dataSize})
	_, err := writer.Write(samples)
	return err
}
