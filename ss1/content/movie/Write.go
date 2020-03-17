package movie

import (
	"bytes"
	"encoding/binary"
	"io"
	"sort"

	"github.com/inkyblackness/hacked/ss1/content/movie/internal/format"
	"github.com/inkyblackness/hacked/ss1/content/text"
)

const indexHeaderSizeIncrement = 0x0400

// Write encodes the provided container into the given writer.
func Write(dest io.Writer, container Container, cp text.Codepage) error {
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

	var buckets []EntryBucket
	buckets = append(buckets, container.Audio.Encode()...)
	buckets = append(buckets, container.Video.Encode()...)
	subtitleBucketsList := container.Subtitles.Encode(cp)
	for _, subtitleBuckets := range subtitleBucketsList {
		buckets = append(buckets, subtitleBuckets...)
	}
	sort.Slice(buckets, func(a, b int) bool {
		bucketA := buckets[a]
		bucketB := buckets[b]
		if bucketA.Timestamp.IsBefore(bucketB.Timestamp) {
			return true
		}
		if bucketB.Timestamp.IsBefore(bucketA.Timestamp) {
			return false
		}
		return bucketA.Priority < bucketB.Priority
	})

	// create index
	var entries []Entry
	for _, bucket := range buckets {
		entries = append(entries, bucket.Entries...)
	}
	moveSubtitlesAfterSceneChange(entries)

	for _, dataEntry := range entries {
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
	for _, dataEntry := range entries {
		_, err = dest.Write(dataEntry.Data.Bytes())
		if err != nil {
			return err
		}
	}
	return nil
}

func moveSubtitlesAfterSceneChange(entries []Entry) {
	nextSceneChangeIndex := -1
	var nextSceneChangeTimestamp Timestamp
	delta := TimestampFromSeconds(0.5)
	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]

		switch entry.Data.Type() {
		case DataTypePalette:
			nextSceneChangeIndex = i
			nextSceneChangeTimestamp = entry.Timestamp
		case DataTypeSubtitle:
			if (nextSceneChangeIndex > 0) && nextSceneChangeTimestamp.IsBefore(entry.Timestamp.Plus(delta)) {
				copy(entries[i:nextSceneChangeIndex], entries[i+1:nextSceneChangeIndex+1])
				entries[nextSceneChangeIndex] = entry
				nextSceneChangeIndex--
			}
		}
	}

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
