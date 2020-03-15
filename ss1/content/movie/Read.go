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
func Read(source io.ReadSeeker) (Container, error) {
	if source == nil {
		return Container{}, fmt.Errorf("source is nil")
	}

	var header format.Header
	startPos, _ := source.Seek(0, io.SeekCurrent)
	var container Container
	err := binary.Read(source, binary.LittleEndian, &header)
	if err != nil {
		return Container{}, err
	}
	err = verifyAndExtractHeader(&container, &header)
	if err != nil {
		return Container{}, err
	}
	err = readPalette(source, &container)
	if err != nil {
		return Container{}, err
	}
	err = readIndexAndEntries(source, startPos, &container, &header)
	if err != nil {
		return Container{}, err
	}

	return container, nil
}

func verifyAndExtractHeader(container *Container, header *format.Header) error {
	if !bytes.Equal(header.Tag[:], bytes.NewBufferString(format.Tag).Bytes()) {
		return errors.New("not a MOVI format")
	}

	container.EndTimestamp = Timestamp{Second: header.DurationSeconds, Fraction: header.DurationFraction}
	container.VideoHeight = header.VideoHeight
	container.VideoWidth = header.VideoWidth
	container.AudioSampleRate = header.SampleRate
	return nil
}

func readPalette(source io.Reader, container *Container) error {
	var pal bitmap.Palette
	decoder := serial.NewDecoder(source)
	decoder.Code(&pal)

	if decoder.FirstError() != nil {
		return decoder.FirstError()
	}

	container.StartPalette = pal
	return nil
}

func readIndexAndEntries(source io.ReadSeeker, startPos int64, container *Container, header *format.Header) error {
	indexEntries := make([]format.IndexTableEntry, header.IndexEntryCount)
	err := binary.Read(source, binary.LittleEndian, indexEntries)
	if err != nil {
		return err
	}
	for index, indexEntry := range indexEntries {
		entryType := DataType(indexEntry.Type)

		if entryType != dataTypeEndOfMedia {
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

			entry, _ := UnknownEntryFrom(timestamp, entryType, data)
			container.AddEntry(entry)
		}
	}
	return nil
}
