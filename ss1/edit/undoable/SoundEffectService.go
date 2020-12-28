package undoable

import (
	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/edit"
	"github.com/inkyblackness/hacked/ss1/edit/media"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
)

// SoundEffectService provides read/write functionality with undo capability.
type SoundEffectService struct {
	wrapped   edit.SoundEffectService
	commander cmd.Commander
}

// NewSoundEffectService returns a new instance of a service.
func NewSoundEffectService(wrapped edit.SoundEffectService, commander cmd.Commander) SoundEffectService {
	return SoundEffectService{
		wrapped:   wrapped,
		commander: commander,
	}
}

// Modified returns true if the identified sound resource is marked as modified.
func (service SoundEffectService) Modified(key resource.Key) bool {
	return service.wrapped.Modified(key)
}

// RequestRemove queues to erase the sound from the resources.
func (service SoundEffectService) RequestRemove(key resource.Key, restoreFunc func()) {
	service.requestCommand(
		func(setter media.SoundEffectBlockSetter) {
			service.wrapped.Remove(setter, key)
		},
		service.wrapped.RestoreFunc(key),
		restoreFunc)
}

// RequestClear queues to reset the identified audio resource to a silent one-sample audio.
func (service SoundEffectService) RequestClear(key resource.Key, restoreFunc func()) {
	service.requestCommand(
		func(setter media.SoundEffectBlockSetter) {
			service.wrapped.Clear(setter, key)
		},
		service.wrapped.RestoreFunc(key),
		restoreFunc)
}

// Audio returns the audio component of identified sound effect.
func (service SoundEffectService) Audio(key resource.Key) audio.L8 {
	return service.wrapped.Audio(key)
}

// RequestSetAudio queues the change to update the audio.
func (service SoundEffectService) RequestSetAudio(key resource.Key, data audio.L8, restoreFunc func()) {
	service.requestCommand(
		func(setter media.SoundEffectBlockSetter) {
			service.wrapped.SetAudio(setter, key, data)
		},
		service.wrapped.RestoreFunc(key),
		restoreFunc)
}

func (service SoundEffectService) requestCommand(
	forward func(modder media.SoundEffectBlockSetter),
	reverse func(modder media.SoundEffectBlockSetter),
	restore func()) {
	c := command{
		forward: func(modder world.Modder) { forward(modder) },
		reverse: func(modder world.Modder) { reverse(modder) },
		restore: restore,
	}
	service.commander.Queue(c)
}
