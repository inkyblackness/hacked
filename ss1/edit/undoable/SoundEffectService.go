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
