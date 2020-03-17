package movie

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ss1/content/movie/internal/compression"
	"github.com/inkyblackness/hacked/ss1/content/movie/internal/format"
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/serial"
)

// Read tries to extract a MOVI container from the provided reader.
// On success the position of the reader is past the last data entry.
// On failure the position of the reader is undefined.
func Read(source io.ReadSeeker, cp text.Codepage) (Container, error) {
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
	err = verifyAndExtractHeader(&header, &container)
	if err != nil {
		return Container{}, err
	}
	var startPalette bitmap.Palette
	err = readPalette(source, &startPalette)
	if err != nil {
		return Container{}, err
	}
	entries, err := readIndexAndEntries(source, startPos, &header)
	if err != nil {
		return Container{}, err
	}
	err = parseEntries(entries, &container,
		cp, startPalette, Timestamp{Second: header.DurationSeconds, Fraction: header.DurationFraction})
	if err != nil {
		return Container{}, err
	}

	return container, nil
}

func verifyAndExtractHeader(header *format.Header, container *Container) error {
	if !bytes.Equal(header.Tag[:], bytes.NewBufferString(format.Tag).Bytes()) {
		return errors.New("not a MOVI format")
	}

	container.Video.Height = header.VideoHeight
	container.Video.Width = header.VideoWidth
	container.Audio.Sound.SampleRate = float32(header.SampleRate)
	return nil
}

func readPalette(source io.Reader, pal *bitmap.Palette) error {
	decoder := serial.NewDecoder(source)
	decoder.Code(pal)
	return decoder.FirstError()
}

func readIndexAndEntries(source io.ReadSeeker, startPos int64, header *format.Header) ([]Entry, error) {
	entries := make([]Entry, 0, header.IndexEntryCount)
	indexEntries := make([]format.IndexTableEntry, header.IndexEntryCount)
	err := binary.Read(source, binary.LittleEndian, indexEntries)
	if err != nil {
		return nil, err
	}
	for index, indexEntry := range indexEntries {
		entryType := format.DataType(indexEntry.Type)

		if entryType != format.DataTypeEndOfMedia {
			timestamp := Timestamp{
				Second:   indexEntry.TimestampSecond,
				Fraction: indexEntry.TimestampFraction,
			}
			length := int(indexEntries[index+1].DataOffset - indexEntry.DataOffset)

			_, err = source.Seek(startPos+int64(indexEntry.DataOffset), io.SeekStart)
			if err != nil {
				return nil, err
			}

			entryData, err := readEntry(entryType, source, length)
			if err != nil {
				return nil, err
			}
			entries = append(entries, Entry{
				Timestamp: timestamp,
				Data:      entryData,
			})
		}
	}
	return entries, nil
}

func readEntry(entryType format.DataType, source io.Reader, dataSize int) (EntryData, error) {
	limited := io.LimitReader(source, int64(dataSize))
	switch entryType {
	case format.DataTypeLowResVideo:
		return LowResVideoEntryFrom(limited, dataSize)
	case format.DataTypeHighResVideo:
		return HighResVideoEntryFrom(limited, dataSize)
	case format.DataTypeAudio:
		return AudioEntryFrom(limited, dataSize)
	case format.DataTypeSubtitle:
		return SubtitleEntryFrom(limited, dataSize)
	case format.DataTypePalette:
		return PaletteEntryFrom(limited)
	case format.DataTypePaletteReset:
		return PaletteResetEntryFrom()
	case format.DataTypePaletteLookupList:
		return PaletteLookupEntryFrom(limited, dataSize)
	case format.DataTypeControlDictionary:
		return ControlDictionaryEntryFrom(limited, dataSize)
	default:
		return UnknownEntryFrom(entryType, limited, dataSize)
	}
}

func parseEntries(entries []Entry, container *Container,
	cp text.Codepage, startPalette bitmap.Palette, endTimestamp Timestamp) error {
	palette := startPalette
	var paletteLookup []byte
	var controlDictionary []compression.ControlWord
	var highResScene *HighResScene
	sceneChanging := true
	finishScene := func() {
		if highResScene != nil {
			container.Video.Scenes = append(container.Video.Scenes, *highResScene)
		}
		highResScene = nil
	}
	var frame *HighResFrame
	finishFrame := func(timestamp Timestamp) {
		if frame != nil {
			frame.displayTime = timestamp.Minus(frame.displayTime)
			highResScene.frames = append(highResScene.frames, *frame)
		}
		frame = nil
	}
	for _, entry := range entries {
		switch data := entry.Data.(type) {
		case AudioEntryData:
			container.Audio.Sound.Samples = append(container.Audio.Sound.Samples, data.Samples...)
		case SubtitleEntryData:
			subText := cp.Decode(data.Text)
			switch data.Control {
			case format.SubtitleTextStd:
				container.Subtitles.Add(resource.LangDefault, entry.Timestamp, subText)
			case format.SubtitleTextGer:
				container.Subtitles.Add(resource.LangGerman, entry.Timestamp, subText)
			case format.SubtitleTextFrn:
				container.Subtitles.Add(resource.LangFrench, entry.Timestamp, subText)
			case format.SubtitleArea:
			default:
			}
		case LowResVideoEntryData:
			// ignored for now
			return fmt.Errorf("low-res video not supported")
		case PaletteLookupEntryData:
			sceneChanging = true
			paletteLookup = data.List
		case ControlDictionaryEntryData:
			sceneChanging = true
			controlDictionary = data.Words
		case PaletteResetEntryData:
			sceneChanging = true
			palette = bitmap.Palette{}
		case PaletteEntryData:
			sceneChanging = true
			palette = data.Colors
		case HighResVideoEntryData:
			finishFrame(entry.Timestamp)
			if sceneChanging {
				finishScene()
				highResScene = &HighResScene{
					palette:       palette,
					paletteLookup: paletteLookup,
					controlWords:  controlDictionary,
				}
				sceneChanging = false
			}
			frame = &HighResFrame{
				bitstream:   data.Bitstream,
				maskstream:  data.Maskstream,
				displayTime: entry.Timestamp,
			}
		default:
		}
	}
	finishFrame(endTimestamp)
	finishScene()

	return nil
}
