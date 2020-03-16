package movie

import (
	"bytes"
	"errors"
	"io/ioutil"

	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ss1/content/movie/internal/compression"
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/serial/rle"
)

// Cache retrieves movie container from a localizer and keeps them decoded until they are invalidated.
type Cache struct {
	cp text.Codepage

	localizer resource.Localizer

	movies map[resource.Key]*cachedMovie
}

type cachedMovie struct {
	cp text.Codepage

	container Container

	sound           *audio.L8
	scenes          []Scene
	subtitlesByLang map[resource.Language]*SubtitleList
}

func (cached *cachedMovie) audio() audio.L8 {
	if cached.sound != nil {
		return *cached.sound
	}
	var samples []byte
	for _, entry := range cached.container.Entries {
		if audioData, isAudio := entry.Data.(AudioEntryData); isAudio {
			samples = append(samples, audioData.Samples...)
		}
	}
	cached.sound = &audio.L8{
		Samples:    samples,
		SampleRate: cached.container.Audio.Sound.SampleRate,
	}
	return *cached.sound
}

func (cached *cachedMovie) video() []Scene {
	if len(cached.scenes) > 0 {
		return cached.scenes
	}

	var scenes []Scene
	currentPalette := cached.container.StartPalette
	width := int(cached.container.Video.Width)
	height := int(cached.container.Video.Height)
	frameBuffer := make([]byte, width*height)
	decoderBuilder := compression.NewFrameDecoderBuilder(width, height)
	decoderBuilder.ForStandardFrame(frameBuffer, width)

	clonePalette := func() *bitmap.Palette {
		paletteCopy := currentPalette
		return &paletteCopy
	}
	cloneFramebuffer := func() []byte {
		bufferCopy := make([]byte, len(frameBuffer))
		copy(bufferCopy, frameBuffer)
		return bufferCopy
	}

	var currentScene *Scene

	setPreviousFrameEndTime := func(ts Timestamp) {
		if currentScene != nil && len(currentScene.Frames) > 0 {
			previousIndex := len(currentScene.Frames) - 1
			previousFrame := currentScene.Frames[previousIndex]
			if ts.IsAfter(previousFrame.DisplayTime) {
				previousFrame.DisplayTime = previousFrame.DisplayTime.DeltaTo(ts)
			} else {
				previousFrame.DisplayTime = Timestamp{}
			}
			currentScene.Frames[previousIndex] = previousFrame
		}
	}
	finishScene := func(now Timestamp) {
		if currentScene != nil {
			setPreviousFrameEndTime(now)
			scenes = append(scenes, *currentScene)
		}
		currentScene = nil
	}
	for _, entry := range cached.container.Entries {
		switch entryData := entry.Data.(type) {
		case PaletteEntryData:
			finishScene(entry.Timestamp)
			currentPalette = entryData.Colors
		case ControlDictionaryEntryData:
			decoderBuilder.WithControlWords(entryData.Words)
		case PaletteLookupEntryData:
			finishScene(entry.Timestamp)
			decoderBuilder.WithPaletteLookupList(entryData.List)
		case LowResVideoEntryData:
			err := rle.Decompress(bytes.NewReader(entryData.Packed), frameBuffer)
			if err != nil {
				break
			}
		case HighResVideoEntryData:
			decoder := decoderBuilder.Build()

			err := decoder.Decode(entryData.Bitstream, entryData.Maskstream)
			if err != nil {
				break
			}
			if currentScene == nil {
				currentScene = &Scene{}
			}

			bmp := bitmap.Bitmap{
				Header: bitmap.Header{
					Type:   bitmap.TypeFlat8Bit,
					Width:  int16(cached.container.Video.Width),
					Height: int16(cached.container.Video.Height),
					Stride: cached.container.Video.Width,
				},
				Palette: clonePalette(),
				Pixels:  cloneFramebuffer(),
			}
			setPreviousFrameEndTime(entry.Timestamp)
			currentScene.Frames = append(currentScene.Frames, Frame{
				Bitmap:      bmp,
				DisplayTime: entry.Timestamp,
			})
		}
	}
	finishScene(cached.container.EndTimestamp)

	cached.scenes = scenes

	return cached.scenes
}

func (cached *cachedMovie) subtitles(language resource.Language) SubtitleList {
	sub := cached.subtitlesByLang[language]
	if sub != nil {
		return *sub
	}

	sub = &SubtitleList{}
	expectedControl := SubtitleControlForLanguage(language)

	for _, entry := range cached.container.Entries {
		subtitleData, isSubtitle := entry.Data.(SubtitleEntryData)
		if !isSubtitle {
			continue
		}
		if subtitleData.Control == expectedControl {
			sub.add(entry.Timestamp, cached.cp.Decode(subtitleData.Text))
		}
	}
	if (len(sub.Entries) > 0) && (len(sub.Entries[len(sub.Entries)-1].Text) > 0) {
		sub.add(cached.container.EndTimestamp, "")
	}
	if cached.subtitlesByLang == nil {
		cached.subtitlesByLang = make(map[resource.Language]*SubtitleList)
	}
	cached.subtitlesByLang[language] = sub
	return *sub
}

// NewCache returns a new instance.
func NewCache(cp text.Codepage, localizer resource.Localizer) *Cache {
	cache := &Cache{
		cp:        cp,
		localizer: localizer,

		movies: make(map[resource.Key]*cachedMovie),
	}
	return cache
}

// InvalidateResources lets the cache remove any movies from resources that are specified in the given slice.
func (cache *Cache) InvalidateResources(ids []resource.ID) {
	for _, id := range ids {
		for key := range cache.movies {
			if key.ID == id {
				delete(cache.movies, key)
			}
		}
	}
}

func (cache *Cache) cached(key resource.Key) (*cachedMovie, error) {
	value, existing := cache.movies[key]
	if existing {
		return value, nil
	}
	selector := cache.localizer.LocalizedResources(key.Lang)
	view, err := selector.Select(key.ID.Plus(key.Index))
	if err != nil {
		return nil, errors.New("no movie found")
	}
	if (view.ContentType() != resource.Movie) || view.Compound() || (view.BlockCount() != 1) {
		return nil, errors.New("invalid resource type")
	}
	reader, err := view.Block(0)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	container, err := Read(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	cached := &cachedMovie{
		cp:        cache.cp,
		container: container,
	}
	cache.movies[key] = cached
	return cached, nil
}

// Container retrieves and caches the underlying movie, and returns the complete container.
func (cache *Cache) Container(key resource.Key) (Container, error) {
	cached, err := cache.cached(key)
	if err != nil {
		return Container{}, err
	}
	return cached.container, nil
}

// Audio retrieves and caches the underlying movie, and returns the audio only.
func (cache *Cache) Audio(key resource.Key) (sound audio.L8, err error) {
	cached, err := cache.cached(key)
	if err != nil {
		return
	}
	return cached.audio(), nil
}

// Video retrieves and caches the underlying movie, and returns the complete list of decoded scenes.
func (cache *Cache) Video(key resource.Key) ([]Scene, error) {
	cached, err := cache.cached(key)
	if err != nil {
		return nil, err
	}
	return cached.video(), nil
}

// Subtitles retrieves and caches the underlying movie, and returns the subtitles for given language.
func (cache *Cache) Subtitles(key resource.Key, language resource.Language) (SubtitleList, error) {
	cached, err := cache.cached(key)
	if err != nil {
		return SubtitleList{}, err
	}
	return cached.subtitles(language), nil
}
