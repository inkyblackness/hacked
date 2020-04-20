package sound

import (
	"errors"

	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/content/audio/voc"
	"github.com/inkyblackness/hacked/ss1/resource"
)

// SoundEffectCache retrieves audio samples stored as VOC files from a localizer and
// keeps them cached until they are invalidated.
type SoundEffectCache struct {
	localizer resource.Localizer

	sounds map[resource.Key]audio.L8
}

// NewSoundCache returns a new instance.
func NewSoundCache(localizer resource.Localizer) *SoundEffectCache {
	cache := &SoundEffectCache{
		localizer: localizer,

		sounds: make(map[resource.Key]audio.L8),
	}
	return cache
}

// InvalidateResources lets the cache remove any sounds from resources that are specified in the given slice.
func (cache *SoundEffectCache) InvalidateResources(ids []resource.ID) {
	for _, id := range ids {
		for key := range cache.sounds {
			if key.ID == id {
				delete(cache.sounds, key)
			}
		}
	}
}

func (cache *SoundEffectCache) cached(key resource.Key) (audio.L8, error) {
	value, existing := cache.sounds[key]
	if existing {
		return value, nil
	}
	selector := cache.localizer.LocalizedResources(key.Lang)
	view, err := selector.Select(key.ID.Plus(key.Index))
	if err != nil {
		return audio.L8{}, errors.New("no sound found")
	}
	if view.ContentType() != resource.Sound {
		return audio.L8{}, errors.New("invalid resource type")
	}
	reader, err := view.Block(0)
	if err != nil {
		return audio.L8{}, err
	}
	sound, err := voc.Load(reader)
	if err != nil {
		return audio.L8{}, err
	}
	cache.sounds[key] = sound
	return sound, nil
}

// Audio retrieves and caches the underlying sound.
func (cache *SoundEffectCache) Audio(key resource.Key) (audio.L8, error) {
	cached, err := cache.cached(key)
	if err != nil {
		return audio.L8{}, err
	}
	return cached, nil
}
