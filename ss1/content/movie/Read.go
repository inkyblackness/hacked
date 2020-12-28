package movie

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/inkyblackness/hacked/ss1"
	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ss1/content/movie/internal/compression"
	"github.com/inkyblackness/hacked/ss1/content/movie/internal/format"
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/serial"
)

const (
	errSourceIsNil             ss1.StringError = "source is nil"
	errInvalidFormat           ss1.StringError = "not a MOVI format"
	errLowResVideoNotSupported ss1.StringError = "low-res video not supported"
)

// Read tries to extract a MOVI container from the provided reader.
// On success the position of the reader is past the last data entry.
// On failure the position of the reader is undefined.
func Read(source io.ReadSeeker, cp text.Codepage) (Container, error) {
	if source == nil {
		return Container{}, errSourceIsNil
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
		cp, startPalette, format.Timestamp{Second: byte(header.Duration.Number), Fraction: header.Duration.Fraction})
	if err != nil {
		return Container{}, err
	}

	return container, nil
}

func verifyAndExtractHeader(header *format.Header, container *Container) error {
	if !bytes.Equal(header.Tag[:], bytes.NewBufferString(format.Tag).Bytes()) {
		return errInvalidFormat
	}

	container.Video.Height = header.VideoHeight
	container.Video.Width = header.VideoWidth
	container.Audio.Sound.SampleRate = header.AudioSampleRate.ToFloat()
	return nil
}

func readPalette(source io.Reader, pal *bitmap.Palette) error {
	decoder := serial.NewDecoder(source)
	decoder.Code(pal)
	return decoder.FirstError()
}

func readIndexAndEntries(source io.ReadSeeker, startPos int64, header *format.Header) ([]format.Entry, error) {
	entries := make([]format.Entry, 0, header.IndexEntryCount)
	indexEntries := make([]format.IndexTableEntry, header.IndexEntryCount)
	err := binary.Read(source, binary.LittleEndian, indexEntries)
	if err != nil {
		return nil, err
	}
	for index, indexEntry := range indexEntries {
		entryType := format.DataType(indexEntry.Type)

		if entryType != format.DataTypeEndOfMedia {
			timestamp := format.Timestamp{
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
			entries = append(entries, format.Entry{
				Timestamp: timestamp,
				Data:      entryData,
			})
		}
	}
	return entries, nil
}

func readEntry(entryType format.DataType, source io.Reader, dataSize int) (format.EntryData, error) {
	limited := io.LimitReader(source, int64(dataSize))
	switch entryType {
	case format.DataTypeLowResVideo:
		return format.LowResVideoEntryFrom(limited, dataSize)
	case format.DataTypeHighResVideo:
		return format.HighResVideoEntryFrom(limited, dataSize)
	case format.DataTypeAudio:
		return format.AudioEntryFrom(limited, dataSize)
	case format.DataTypeSubtitle:
		return format.SubtitleEntryFrom(limited, dataSize)
	case format.DataTypePalette:
		return format.PaletteEntryFrom(limited)
	case format.DataTypePaletteReset:
		return format.PaletteResetEntryFrom()
	case format.DataTypePaletteLookupList:
		return format.PaletteLookupEntryFrom(limited, dataSize)
	case format.DataTypeControlDictionary:
		return format.ControlDictionaryEntryFrom(limited, dataSize)
	default:
		return format.UnknownEntryFrom(entryType, limited, dataSize)
	}
}

func parseEntries(entries []format.Entry, container *Container,
	cp text.Codepage, startPalette bitmap.Palette, endTimestamp format.Timestamp) error {
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
	finishFrame := func(timestamp format.Timestamp) {
		if frame != nil {
			frame.displayTime = timestamp.Minus(frame.displayTime)
			highResScene.frames = append(highResScene.frames, *frame)
		}
		frame = nil
	}
	for _, entry := range entries {
		switch data := entry.Data.(type) {
		case format.AudioEntryData:
			container.Audio.Sound.Samples = append(container.Audio.Sound.Samples, data.Samples...)
		case format.SubtitleEntryData:
			subText := cp.Decode(data.Text)
			switch data.Control {
			case format.SubtitleTextStd:
				container.Subtitles.add(resource.LangDefault, entry.Timestamp.ToDuration(), subText)
			case format.SubtitleTextGer:
				container.Subtitles.add(resource.LangGerman, entry.Timestamp.ToDuration(), subText)
			case format.SubtitleTextFrn:
				container.Subtitles.add(resource.LangFrench, entry.Timestamp.ToDuration(), subText)
			case format.SubtitleArea:
			default:
			}
		case format.LowResVideoEntryData:
			// ignored for now
			return errLowResVideoNotSupported
		case format.PaletteLookupEntryData:
			sceneChanging = true
			paletteLookup = data.List
		case format.ControlDictionaryEntryData:
			sceneChanging = true
			controlDictionary = data.Words
		case format.PaletteResetEntryData:
			sceneChanging = true
			palette = bitmap.Palette{}
		case format.PaletteEntryData:
			sceneChanging = true
			palette = data.Colors
		case format.HighResVideoEntryData:
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
