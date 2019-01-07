package media

import (
	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/content/movie"
	"github.com/inkyblackness/hacked/ss1/resource"
)

type AudioBlockGetter interface {
	ModifiedBlock(lang resource.Language, id resource.ID, index int) []byte
}

type GetAudioService struct {
	movieCache *movie.Cache
	getter     AudioBlockGetter
}

func NewGetAudioService(movieCache *movie.Cache, getter AudioBlockGetter) GetAudioService {
	return GetAudioService{
		movieCache: movieCache,
		getter:     getter,
	}
}

func (service GetAudioService) Get(key resource.Key) audio.L8 {
	sound, _ := service.movieCache.Audio(key)
	return sound
}

func (service GetAudioService) Modified(key resource.Key) bool {
	data := service.getter.ModifiedBlock(key.Lang, key.ID, 0)
	return len(data) > 0
}
