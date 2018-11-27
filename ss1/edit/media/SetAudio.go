package media

import (
	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/content/movie"
	"github.com/inkyblackness/hacked/ss1/resource"
)

type AudioBlockSetter interface {
	SetResourceBlocks(lang resource.Language, id resource.ID, data [][]byte)
	DelResource(lang resource.Language, id resource.ID)
}

type SetAudioService struct{}

func NewSetAudioService() SetAudioService {
	return SetAudioService{}
}

func (service SetAudioService) Remove(setter AudioBlockSetter, key resource.Key) {
	setter.DelResource(key.Lang, key.ID)
}

func (service SetAudioService) Clear(setter AudioBlockSetter, key resource.Key) {
	silence := audio.L8{SampleRate: 22050, Samples: []byte{0x80}}
	service.Set(setter, key, silence)
}

func (service SetAudioService) Set(setter AudioBlockSetter, key resource.Key, sound audio.L8) {
	movieData := movie.ContainSoundData(sound)
	blockData := [][]byte{movieData}
	setter.SetResourceBlocks(key.Lang, key.ID, blockData)
}
