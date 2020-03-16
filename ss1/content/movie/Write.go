package movie

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/inkyblackness/hacked/ss1/content/movie/internal/format"
)

const indexHeaderSizeIncrement = 0x0400

// Write encodes the provided container into the given writer.
func Write(dest io.Writer, container Container) error {
	var indexEntries []format.IndexTableEntry
	var header format.Header
	palette := paletteDataFromContainer(container)

	// setup header
	copy(header.Tag[:], bytes.NewBufferString(format.Tag).Bytes())
	header.DurationSeconds = container.EndTimestamp.Second
	header.DurationFraction = container.EndTimestamp.Fraction
	header.VideoWidth = container.Video.Width
	header.VideoHeight = container.Video.Height
	header.SampleRate = uint16(container.Audio.Sound.SampleRate)

	if header.VideoWidth != 0 {
		header.Unknown001C = 0x0008
		header.Unknown001E = 0x0001
	}
	header.Unknown0020 = 0x0001
	header.Unknown0022 = 0x00000001

	// create index
	for _, dataEntry := range container.Entries {
		indexEntry := format.IndexTableEntry{
			Type:              byte(dataEntry.Data.Type()),
			DataOffset:        header.ContentSize,
			TimestampSecond:   dataEntry.Timestamp.Second,
			TimestampFraction: dataEntry.Timestamp.Fraction,
		}

		header.ContentSize += int32(len(dataEntry.Data.Bytes()))
		indexEntries = append(indexEntries, indexEntry)
	}
	// calculate size fields
	header.IndexSize = int32(indexTableSizeFor(len(indexEntries) + 1))
	dataStartOffset := format.HeaderSize + int32(len(palette)) + header.IndexSize
	for i := range indexEntries {
		indexEntries[i].DataOffset += dataStartOffset
	}

	// determine end
	lastEntry := format.IndexTableEntry{
		TimestampSecond:   header.DurationSeconds,
		TimestampFraction: header.DurationFraction,
		Type:              byte(dataTypeEndOfMedia),
		DataOffset:        dataStartOffset + header.ContentSize}
	indexEntries = append(indexEntries, lastEntry)
	header.IndexEntryCount = int32(len(indexEntries))

	// write data
	err := binary.Write(dest, binary.LittleEndian, &header)
	if err != nil {
		return err
	}
	_, err = dest.Write(palette)
	if err != nil {
		return err
	}
	err = binary.Write(dest, binary.LittleEndian, indexEntries)
	if err != nil {
		return err
	}
	_, err = dest.Write(make([]byte, int(header.IndexSize)-int(header.IndexEntryCount)*format.IndexTableEntrySize))
	if err != nil {
		return err
	}
	for _, dataEntry := range container.Entries {
		_, err = dest.Write(dataEntry.Data.Bytes())
		if err != nil {
			return err
		}
	}
	return nil
}

func paletteDataFromContainer(container Container) []byte {
	palette := container.StartPalette
	buf := bytes.NewBuffer(nil)
	_ = binary.Write(buf, binary.LittleEndian, &palette)
	return buf.Bytes()
}

func indexTableSizeFor(entryCount int) int {
	size := indexHeaderSizeIncrement
	requiredSize := entryCount * format.IndexTableEntrySize

	if requiredSize > size {
		size *= requiredSize/size + 2
	}

	return size
}
