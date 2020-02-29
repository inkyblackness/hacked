package movie

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ss1/content/movie/internal/format"
	"github.com/inkyblackness/hacked/ss1/serial"
)

// Read tries to extract a MOVI container from the provided reader.
// On success the position of the reader is past the last data entry.
// On failure the position of the reader is undefined.
func Read(source io.ReadSeeker) (container Container, err error) {
	if source == nil {
		return nil, fmt.Errorf("source is nil")
	}

	var header format.Header
	startPos, _ := source.Seek(0, io.SeekCurrent)

	err = binary.Read(source, binary.LittleEndian, &header)
	if err != nil {
		return
	}
	builder := NewContainerBuilder()
	err = verifyAndExtractHeader(builder, &header)
	if err != nil {
		return
	}
	err = readPalette(source, builder)
	if err != nil {
		return
	}
	err = readIndexAndEntries(source, startPos, builder, &header)
	if err != nil {
		return
	}

	return builder.Build(), nil
}

func verifyAndExtractHeader(builder *ContainerBuilder, header *format.Header) error {
	if !bytes.Equal(header.Tag[:], bytes.NewBufferString(format.Tag).Bytes()) {
		return errors.New("not a MOVI format")
	}

	builder.EndTimestamp(Timestamp{Second: header.DurationSeconds, Fraction: header.DurationFraction})
	builder.VideoHeight(header.VideoHeight)
	builder.VideoWidth(header.VideoWidth)
	builder.AudioSampleRate(header.SampleRate)
	return nil
}

func readPalette(source io.Reader, builder *ContainerBuilder) error {
	var pal bitmap.Palette
	decoder := serial.NewDecoder(source)
	decoder.Code(&pal)

	if decoder.FirstError() != nil {
		return decoder.FirstError()
	}

	builder.StartPalette(&pal)
	return nil
}

func readIndexAndEntries(source io.ReadSeeker, startPos int64, builder *ContainerBuilder, header *format.Header) error {
	indexEntries := make([]format.IndexTableEntry, header.IndexEntryCount)
	err := binary.Read(source, binary.LittleEndian, indexEntries)
	if err != nil {
		return err
	}
	for index, indexEntry := range indexEntries {
		entryType := DataType(indexEntry.Type)

		if entryType != endOfMedia {
			timestamp := Timestamp{
				Second:   indexEntry.TimestampSecond,
				Fraction: indexEntry.TimestampFraction,
			}
			length := int(indexEntries[index+1].DataOffset - indexEntry.DataOffset)
			data := make([]byte, length)

			_, err = source.Seek(startPos+int64(indexEntry.DataOffset), 0)
			if err != nil {
				return err
			}
			_, err = source.Read(data)
			if err != nil {
				return err
			}

			builder.AddEntry(NewMemoryEntry(timestamp, entryType, data))
		}
	}
	return nil
}
