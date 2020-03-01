package movie

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io/ioutil"

	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/resource"
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
	sceneFrames     [][]bitmap.Bitmap
	subtitlesByLang map[resource.Language]*Subtitles
}

func (cached *cachedMovie) audio() audio.L8 {
	if cached.sound != nil {
		return *cached.sound
	}
	var samples []byte
	for index := 0; index < cached.container.EntryCount(); index++ {
		entry := cached.container.Entry(index)
		if entry.Type() == Audio {
			samples = append(samples, entry.Data()...)
		}
	}
	cached.sound = &audio.L8{
		Samples:    samples,
		SampleRate: float32(cached.container.AudioSampleRate()),
	}
	return *cached.sound
}

func (cached *cachedMovie) video() [][]bitmap.Bitmap {
	if len(cached.sceneFrames) > 0 {
		return cached.sceneFrames
	}
	// TODO
	return nil
}

func (cached *cachedMovie) subtitles(language resource.Language) Subtitles {
	sub := cached.subtitlesByLang[language]
	if sub != nil {
		return *sub
	}

	sub = &Subtitles{}
	expectedControl := SubtitleControlForLanguage(language)

	for index := 0; index < cached.container.EntryCount(); index++ {
		entry := cached.container.Entry(index)
		if entry.Type() != Subtitle {
			continue
		}
		var subtitleHeader SubtitleHeader
		err := binary.Read(bytes.NewReader(entry.Data()), binary.LittleEndian, &subtitleHeader)
		if err != nil {
			continue
		}
		if subtitleHeader.Control == expectedControl {
			sub.add(entry.Timestamp(), cached.cp.Decode(entry.Data()[SubtitleHeaderSize:]))
		}
	}
	if (len(sub.entries) > 0) && (len(sub.entries[len(sub.entries)-1].Text) > 0) {
		sub.add(cached.container.EndTimestamp(), "")
	}
	if cached.subtitlesByLang == nil {
		cached.subtitlesByLang = make(map[resource.Language]*Subtitles)
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

// Audio retrieves and caches the underlying movie, and returns the audio only.
func (cache *Cache) Audio(key resource.Key) (sound audio.L8, err error) {
	cached, err := cache.cached(key)
	if err != nil {
		return
	}
	return cached.audio(), nil
}

func (cache *Cache) Video(key resource.Key) ([][]bitmap.Bitmap, error) {
	cached, err := cache.cached(key)
	if err != nil {
		return nil, err
	}
	return cached.video(), nil
}

func (cache *Cache) Subtitles(key resource.Key, language resource.Language) (Subtitles, error) {
	cached, err := cache.cached(key)
	if err != nil {
		return Subtitles{}, err
	}
	return cached.subtitles(language), nil
}
