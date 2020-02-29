package movie

import (
	"bytes"
	"errors"
	"io/ioutil"

	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/resource"
)

// Cache retrieves movie container from a localizer and keeps them decoded until they are invalidated.
type Cache struct {
	localizer resource.Localizer

	movies map[resource.Key]*cachedMovie
}

type cachedMovie struct {
	container Container

	sound *audio.L8
}

func (cached *cachedMovie) audio() (sound audio.L8) {
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

// NewCache returns a new instance.
func NewCache(localizer resource.Localizer) *Cache {
	cache := &Cache{
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
	cached := &cachedMovie{container: container}
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
