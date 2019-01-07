package media

import (
	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/content/movie"
	"github.com/inkyblackness/hacked/ss1/resource"
)

// AudioBlockGetter provides raw data of blocks.
type AudioBlockGetter interface {
	ModifiedBlock(lang resource.Language, id resource.ID, index int) []byte
}

// AudioViewerService provides read-only access to audio resources.
type AudioViewerService struct {
	movieCache *movie.Cache
	getter     AudioBlockGetter
}

// NewAudioViewerService returns a new instance.
func NewAudioViewerService(movieCache *movie.Cache, getter AudioBlockGetter) AudioViewerService {
	return AudioViewerService{
		movieCache: movieCache,
		getter:     getter,
	}
}

// Audio returns the sound data associated with the given key.
func (service AudioViewerService) Audio(key resource.Key) audio.L8 {
	sound, _ := service.movieCache.Audio(key)
	return sound
}

// Modified returns true if the identified audio resource is marked as modified.
func (service AudioViewerService) Modified(key resource.Key) bool {
	data := service.getter.ModifiedBlock(key.Lang, key.ID, 0)
	return len(data) > 0
}
