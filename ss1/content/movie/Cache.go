package movie

import (
	"bytes"
	"errors"
	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/resource"
	"io/ioutil"
)

// Cache retrieves movie container from a localizer and keeps them decoded until they are invalidated.
type Cache struct {
	localizer resource.Localizer

	movies map[resource.Key]Container
}

// NewCache returns a new instance.
func NewCache(localizer resource.Localizer) *Cache {
	cache := &Cache{
		localizer: localizer,

		movies: make(map[resource.Key]Container),
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

// Movie retrieves and caches the message of given key.
func (cache *Cache) Movie(key resource.Key) (Container, error) {
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
	value, err = Read(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	cache.movies[key] = value
	return value, nil
}

// Audio retrieves and caches the underlying movie, and returns the audio only.
func (cache *Cache) Audio(key resource.Key) (sound audio.L8, err error) {
	container, err := cache.Movie(key)
	if err != nil {
		return
	}
	var samples []byte
	for index := 0; index < container.EntryCount(); index++ {
		entry := container.Entry(index)
		if entry.Type() == Audio {
			samples = append(samples, entry.Data()...)
		}
	}
	return audio.L8{
		Samples:    samples,
		SampleRate: float32(container.AudioSampleRate()),
	}, nil
}
