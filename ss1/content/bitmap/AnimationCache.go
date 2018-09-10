package bitmap

import (
	"errors"

	"github.com/inkyblackness/hacked/ss1/resource"
)

// AnimationCache retrieves animations from a localizer and keeps them decoded until they are invalidated.
type AnimationCache struct {
	localizer resource.Localizer

	animations map[resource.Key]Animation
}

// NewAnimationCache returns a new instance.
func NewAnimationCache(localizer resource.Localizer) *AnimationCache {
	cache := &AnimationCache{
		localizer:  localizer,
		animations: make(map[resource.Key]Animation),
	}
	return cache
}

// InvalidateResources lets the cache remove any animations from resources that are specified in the given slice.
func (cache *AnimationCache) InvalidateResources(ids []resource.ID) {
	for _, id := range ids {
		for key := range cache.animations {
			if key.ID == id {
				delete(cache.animations, key)
			}
		}
	}
}

// Animation tries to look up given animation.
func (cache *AnimationCache) Animation(key resource.Key) (anim Animation, err error) {
	anim, existing := cache.animations[key]
	if existing {
		return anim, nil
	}
	selector := cache.localizer.LocalizedResources(key.Lang)
	view, err := selector.Select(key.ID)
	if err != nil {
		return
	}
	if (view.ContentType() != resource.Animation) || (view.BlockCount() != 1) {
		return anim, errors.New("resource is not an animation")
	}
	reader, err := view.Block(0)
	if err != nil {
		return
	}
	anim, err = ReadAnimation(reader)
	if err != nil {
		return
	}
	cache.animations[key] = anim
	return anim, nil
}
