package media

import (
	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/content/movie"
	"github.com/inkyblackness/hacked/ss1/resource"
)

// AudioBlockSetter modifies storage of raw resource data.
type AudioBlockSetter interface {
	SetResourceBlocks(lang resource.Language, id resource.ID, data [][]byte)
	DelResource(lang resource.Language, id resource.ID)
}

// AudioSetterService provides methods to change audio resources.
type AudioSetterService struct{}

// NewAudioSetterService returns a new instance.
func NewAudioSetterService() AudioSetterService {
	return AudioSetterService{}
}

// Remove deletes any audio resource for given key.
func (service AudioSetterService) Remove(setter AudioBlockSetter, key resource.Key) {
	setter.DelResource(key.Lang, key.ID)
}

// Clear resets the identified audio resource to a silent one-sample audio.
func (service AudioSetterService) Clear(setter AudioBlockSetter, key resource.Key) {
	silence := audio.L8{SampleRate: 22050, Samples: []byte{0x80}}
	service.Set(setter, key, silence)
}

// Set stores the given sound as the identified resource.
func (service AudioSetterService) Set(setter AudioBlockSetter, key resource.Key, sound audio.L8) {
	movieData := movie.ContainSoundData(sound)
	blockData := [][]byte{movieData}
	setter.SetResourceBlocks(key.Lang, key.ID, blockData)
}
