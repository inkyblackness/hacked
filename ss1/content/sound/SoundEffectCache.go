package sound

import (
	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/content/audio/voc"
	"github.com/inkyblackness/hacked/ss1/resource"
)

// EffectCache retrieves audio samples stored as VOC files from a localizer and
// keeps them cached until they are invalidated.
type EffectCache struct {
	localizer resource.Localizer

	sounds map[resource.Key]audio.L8
}

// NewEffectCache returns a new instance.
func NewEffectCache(localizer resource.Localizer) *EffectCache {
	cache := &EffectCache{
		localizer: localizer,

		sounds: make(map[resource.Key]audio.L8),
	}
	return cache
}

// InvalidateResources lets the cache remove any sounds from resources that are specified in the given slice.
func (cache *EffectCache) InvalidateResources(ids []resource.ID) {
	for _, id := range ids {
		for key := range cache.sounds {
			if key.ID == id {
				delete(cache.sounds, key)
			}
		}
	}
}

func (cache *EffectCache) cached(key resource.Key) (audio.L8, error) {
	value, existing := cache.sounds[key]
	if existing {
		return value, nil
	}
	selector := cache.localizer.LocalizedResources(key.Lang)
	view, err := selector.Select(key.ID.Plus(key.Index))
	if err != nil {
		return audio.L8{}, err
	}
	if view.ContentType() != resource.Sound {
		return audio.L8{}, resource.ErrWrongType(key, resource.Sound)
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
func (cache *EffectCache) Audio(key resource.Key) (audio.L8, error) {
	cached, err := cache.cached(key)
	if err != nil {
		return audio.L8{}, err
	}
	return cached, nil
}
