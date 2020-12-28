package media

import (
	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/content/sound"
	"github.com/inkyblackness/hacked/ss1/resource"
)

// SoundEffectBlockGetter provides raw data of blocks.
type SoundEffectBlockGetter interface {
	ModifiedBlock(lang resource.Language, id resource.ID, index int) []byte
}

// SoundEffectViewerService provides read-only access to audio resources.
type SoundEffectViewerService struct {
	soundCache *sound.EffectCache
	getter     SoundEffectBlockGetter
}

// NewSoundViewerService returns a new instance.
func NewSoundViewerService(soundCache *sound.EffectCache, getter SoundEffectBlockGetter) SoundEffectViewerService {
	return SoundEffectViewerService{
		soundCache: soundCache,
		getter:     getter,
	}
}

// Audio returns the audio data associated with the given key.
func (service SoundEffectViewerService) Audio(key resource.Key) audio.L8 {
	sound, _ := service.soundCache.Audio(key)
	return sound
}

// Modified returns true if the identified sound resource is marked as modified.
func (service SoundEffectViewerService) Modified(key resource.Key) bool {
	data := service.getter.ModifiedBlock(key.Lang, key.ID.Plus(key.Index), 0)
	return len(data) > 0
}
