package media

import (
	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/content/movie"
	"github.com/inkyblackness/hacked/ss1/resource"
)

// MovieBlockGetter provides raw data of blocks.
type MovieBlockGetter interface {
	ModifiedBlock(lang resource.Language, id resource.ID, index int) []byte
}

// MovieViewerService provides read-only access to movie resources.
type MovieViewerService struct {
	movieCache *movie.Cache
	getter     MovieBlockGetter
}

// NewMovieViewerService returns a new instance.
func NewMovieViewerService(movieCache *movie.Cache, getter MovieBlockGetter) MovieViewerService {
	return MovieViewerService{
		movieCache: movieCache,
		getter:     getter,
	}
}

// Modified returns true if the identified text resource is marked as modified.
func (service MovieViewerService) Modified(key resource.Key) bool {
	return len(service.getter.ModifiedBlock(key.Lang, key.ID, key.Index)) > 0
}

// Container returns the base container of the movie.
func (service MovieViewerService) Container(key resource.Key) (movie.Container, error) {
	return service.movieCache.Container(key)
}

// Audio returns the audio data associated with the given key.
func (service MovieViewerService) Audio(key resource.Key) audio.L8 {
	currentValue, cacheErr := service.movieCache.Audio(key)
	if cacheErr != nil {
		return audio.L8{}
	}
	return currentValue
}

// Subtitles returns the subtitles associated with the given key.
func (service MovieViewerService) Subtitles(key resource.Key, language resource.Language) movie.Subtitles {
	currentValue, cacheErr := service.movieCache.Subtitles(key, language)
	if cacheErr != nil {
		return movie.Subtitles{}
	}
	return currentValue
}
