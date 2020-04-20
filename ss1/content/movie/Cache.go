package movie

import (
	"bytes"
	"errors"
	"io/ioutil"

	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/resource"
)

const sizeLimit = 0xFFFFFF

// Cache retrieves movie container from a localizer and keeps them decoded until they are invalidated.
type Cache struct {
	cp text.Codepage

	localizer resource.Localizer

	movies map[resource.Key]*cachedMovie
}

type cachedMovie struct {
	sizeWarning bool
	container   Container
	scenes      []Scene
}

func (cached *cachedMovie) video() []Scene {
	if len(cached.scenes) > 0 {
		return cached.scenes
	}
	cached.scenes, _ = cached.container.Video.Decompress()
	return cached.scenes
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
	container, err := Read(bytes.NewReader(data), cache.cp)
	if err != nil {
		return nil, err
	}
	cached := &cachedMovie{
		sizeWarning: len(data) > sizeLimit,
		container:   container,
	}
	cache.movies[key] = cached
	return cached, nil
}

// SizeWarning returns true if the given movie is cached and is larger than the underlying archive would support.
func (cache *Cache) SizeWarning(key resource.Key) bool {
	cached, existing := cache.movies[key]
	return existing && cached.sizeWarning
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
	return cached.container.Audio.Sound, nil
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
	return cached.container.Subtitles.PerLanguage[language], nil
}
