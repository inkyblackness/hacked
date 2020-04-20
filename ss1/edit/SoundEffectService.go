package edit

import (
	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/edit/media"
	"github.com/inkyblackness/hacked/ss1/resource"
)

// SoundEffectService provides read/write functionality.
type SoundEffectService struct {
	viewer media.SoundEffectViewerService
	setter media.SoundEffectSetterService
}

// NewSoundEffectService returns a new instance based on given accessor.
func NewSoundEffectService(
	viewer media.SoundEffectViewerService, setter media.SoundEffectSetterService) SoundEffectService {
	return SoundEffectService{
		viewer: viewer,
		setter: setter,
	}
}

// RestoreFunc creates a snapshot of the current sound and returns a function to restore it.
func (service SoundEffectService) RestoreFunc(key resource.Key) func(setter media.SoundEffectBlockSetter) {
	oldAudio := service.viewer.Audio(key)
	isModified := service.viewer.Modified(key)

	return func(setter media.SoundEffectBlockSetter) {
		if isModified {
			service.setter.Set(setter, key, oldAudio)
		} else {
			service.setter.Remove(setter, key)
		}
	}
}

// Modified returns true if the identified sound resource is marked as modified.
func (service SoundEffectService) Modified(key resource.Key) bool {
	return service.viewer.Modified(key)
}

// Remove erases the sound from the resources.
func (service SoundEffectService) Remove(setter media.SoundEffectBlockSetter, key resource.Key) {
	service.setter.Remove(setter, key)
}

// Clear resets the identified audio resource to a silent one-sample audio.
func (service SoundEffectService) Clear(setter media.SoundEffectBlockSetter, key resource.Key) {
	service.setter.Clear(setter, key)
}

// Audio returns the audio component of identified sound effect.
func (service SoundEffectService) Audio(key resource.Key) audio.L8 {
	return service.viewer.Audio(key)
}

// SetAudio sets the audio component of identified sound effect.
func (service SoundEffectService) SetAudio(setter media.SoundEffectBlockSetter, key resource.Key, data audio.L8) {
	service.setter.Set(setter, key, data)
}
