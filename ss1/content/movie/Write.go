package movie

import (
	"bytes"
	"encoding/binary"
	"io"
	"sort"

	"github.com/inkyblackness/hacked/ss1/content/movie/internal/format"
	"github.com/inkyblackness/hacked/ss1/content/text"
)

// Write encodes the provided container into the given writer.
func Write(dest io.Writer, container Container, cp text.Codepage) error {
	var indexEntries []format.IndexTableEntry
	var header format.Header
	palette := paletteDataFromContainer(container)

	// setup header
	copy(header.Tag[:], bytes.NewBufferString(format.Tag).Bytes())
	endTimestamp := container.duration()
	header.Duration.Number = int16(endTimestamp.Second)
	header.Duration.Fraction = endTimestamp.Fraction
	header.VideoWidth = container.Video.Width
	header.VideoHeight = container.Video.Height

	if header.VideoWidth != 0 {
		header.VideoBitsPerPixel = 8
		header.VideoPalettePresent = 1
		header.VideoFrameRate = format.FixFromFloat(12.5) // take a good guess on what a typical framerate would be.
	}
	if !container.Audio.Sound.Empty() {
		header.AudioChannelCount = 1
		header.AudioBytesPerSample = 1
		header.AudioSampleRate = format.FixFromFloat(container.Audio.Sound.SampleRate)
	}

	var buckets []format.EntryBucket
	buckets = append(buckets, container.Audio.encode()...)
	buckets = append(buckets, container.Video.encode()...)
	subtitleBucketsList := container.Subtitles.encode(cp)
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
	var entries []format.Entry
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
	header.IndexSize = int32(format.IndexTableSizeFor(len(indexEntries) + 1))
	dataStartOffset := format.HeaderSize + int32(len(palette)) + header.IndexSize
	for i := range indexEntries {
		indexEntries[i].DataOffset += dataStartOffset
	}

	// determine end
	lastEntry := format.IndexTableEntry{
		TimestampSecond:   byte(header.Duration.Number),
		TimestampFraction: header.Duration.Fraction,
		Type:              byte(format.DataTypeEndOfMedia),
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

func moveSubtitlesAfterSceneChange(entries []format.Entry) {
	nextSceneChangeIndex := -1
	var nextSceneChangeTimestamp format.Timestamp
	delta := format.TimestampFromSeconds(0.5)
	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]

		switch entry.Data.Type() {
		case format.DataTypePalette:
			nextSceneChangeIndex = i
			nextSceneChangeTimestamp = entry.Timestamp
		case format.DataTypeSubtitle:
			if (nextSceneChangeIndex > 0) && nextSceneChangeTimestamp.IsBefore(entry.Timestamp.Plus(delta)) {
				copy(entries[i:nextSceneChangeIndex], entries[i+1:nextSceneChangeIndex+1])
				entries[nextSceneChangeIndex] = entry
				nextSceneChangeIndex--
			}
		}
	}
}

func paletteDataFromContainer(container Container) []byte {
	palette := container.Video.StartPalette()
	buf := bytes.NewBuffer(nil)
	_ = binary.Write(buf, binary.LittleEndian, &palette)
	return buf.Bytes()
}
