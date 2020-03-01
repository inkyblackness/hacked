package movie

import (
	"bytes"
	"encoding/binary"

	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ss1/content/movie/internal/compression"
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/serial"
	"github.com/inkyblackness/hacked/ss1/serial/rle"
)

// MediaDispatcher iterates through the entries of a container and provides resulting media
// to a handler. The dispatcher handles intermediate container entries to create consumable media.
type MediaDispatcher struct {
	handler MediaHandler

	container Container
	nextIndex int

	codepage text.Codepage

	palette        bitmap.Palette
	decoderBuilder *compression.FrameDecoderBuilder
	frameBuffer    []byte
}

// NewMediaDispatcher returns a new instance of a dispatcher reading the provided container.
func NewMediaDispatcher(container Container, handler MediaHandler) *MediaDispatcher {
	width := int(container.VideoWidth)
	height := int(container.VideoHeight)
	dispatcher := &MediaDispatcher{
		handler:        handler,
		container:      container,
		codepage:       text.DefaultCodepage(),
		frameBuffer:    make([]byte, width*height),
		decoderBuilder: compression.NewFrameDecoderBuilder(width, height)}

	startPalette := container.StartPalette
	dispatcher.setPalette(&startPalette)
	dispatcher.decoderBuilder.ForStandardFrame(dispatcher.frameBuffer, width)

	return dispatcher
}

// DispatchNext processes the next entries from the container to call the handler.
// Returns false if the dispatcher reached the end of the container.
func (dispatcher *MediaDispatcher) DispatchNext() (result bool, err error) {
	for !result && (dispatcher.nextIndex < len(dispatcher.container.Entries)) {
		entry := dispatcher.container.Entries[dispatcher.nextIndex]
		result, err = dispatcher.process(entry)
		dispatcher.nextIndex++
	}

	return
}

func (dispatcher *MediaDispatcher) process(entry Entry) (dispatched bool, err error) {
	switch entry.Type() {
	case Audio:
		{
			dispatcher.handler.OnAudio(entry.Timestamp(), entry.Data())
			dispatched = true
		}
	case Subtitle:
		{
			var subtitleHeader SubtitleHeader

			err = binary.Read(bytes.NewReader(entry.Data()), binary.LittleEndian, &subtitleHeader)
			if err != nil {
				return
			}
			subtitle := dispatcher.codepage.Decode(entry.Data()[SubtitleHeaderSize:])
			dispatcher.handler.OnSubtitle(entry.Timestamp(), subtitleHeader.Control, subtitle)
			dispatched = true
		}

	case Palette:
		{
			var pal bitmap.Palette
			decoder := serial.NewDecoder(bytes.NewReader(entry.Data()))
			decoder.Code(&pal)
			if decoder.FirstError() == nil {
				dispatcher.setPalette(&pal)
				dispatcher.clearFrameBuffer()
			} else {
				err = decoder.FirstError()
			}
		}
	case ControlDictionary:
		{
			words, wordsErr := compression.UnpackControlWords(entry.Data())

			if wordsErr == nil {
				dispatcher.decoderBuilder.WithControlWords(words)
			} else {
				err = wordsErr
			}
		}
	case PaletteLookupList:
		{
			dispatcher.decoderBuilder.WithPaletteLookupList(entry.Data())
		}

	case LowResVideo:
		{
			var videoHeader LowResVideoHeader
			reader := bytes.NewReader(entry.Data())

			err = binary.Read(reader, binary.LittleEndian, &videoHeader)
			if err != nil {
				return
			}
			frameErr := rle.Decompress(reader, dispatcher.frameBuffer)
			if frameErr == nil {
				dispatcher.notifyVideoFrame(entry.Timestamp())
				dispatched = true
			} else {
				err = frameErr
			}
		}
	case HighResVideo:
		{
			var videoHeader HighResVideoHeader
			reader := bytes.NewReader(entry.Data())

			err = binary.Read(reader, binary.LittleEndian, &videoHeader)
			if err != nil {
				return
			}
			bitstreamData := entry.Data()[HighResVideoHeaderSize:videoHeader.PixelDataOffset]
			maskstreamData := entry.Data()[videoHeader.PixelDataOffset:]
			decoder := dispatcher.decoderBuilder.Build()

			err = decoder.Decode(bitstreamData, maskstreamData)
			if err != nil {
				return
			}
			dispatcher.notifyVideoFrame(entry.Timestamp())
			dispatched = true
		}
	}

	return
}

func (dispatcher *MediaDispatcher) setPalette(newPalette *bitmap.Palette) {
	dispatcher.palette = *newPalette
}

func (dispatcher *MediaDispatcher) clearFrameBuffer() {
	for pixel := 0; pixel < len(dispatcher.frameBuffer); pixel++ {
		dispatcher.frameBuffer[pixel] = 0x00
	}
}

func (dispatcher *MediaDispatcher) notifyVideoFrame(timestamp Timestamp) {
	bmp := bitmap.Bitmap{
		Header: bitmap.Header{
			Type:   bitmap.TypeFlat8Bit,
			Width:  int16(dispatcher.container.VideoWidth),
			Height: int16(dispatcher.container.VideoHeight),
			Stride: dispatcher.container.VideoWidth,
		},
		Palette: &dispatcher.palette,
		Pixels:  dispatcher.frameBuffer,
	}
	dispatcher.handler.OnVideo(timestamp, bmp)
}
