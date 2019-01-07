package undoable

import (
	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/edit"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
)

// AugmentedTextService provides read/write functionality with undo capability.
type AugmentedTextService struct {
	wrapped   edit.AugmentedTextService
	commander cmd.Commander
}

// NewAugmentedTextService returns a new instance of a service.
func NewAugmentedTextService(wrapped edit.AugmentedTextService, commander cmd.Commander) AugmentedTextService {
	return AugmentedTextService{
		wrapped:   wrapped,
		commander: commander,
	}
}

// IsTrapMessage returns true if the provided resource identifies a trap resource.
func (service AugmentedTextService) IsTrapMessage(key resource.Key) bool {
	return service.wrapped.IsTrapMessage(key)
}

// Text returns the textual value of the identified text resource.
func (service AugmentedTextService) Text(key resource.Key) string {
	return service.wrapped.Text(key)
}

// RequestSetText queues the change to update the text.
func (service AugmentedTextService) RequestSetText(key resource.Key, value string, restoreFunc func()) {
	service.requestCommand(
		func(setter edit.AugmentedTextBlockSetter) {
			service.wrapped.SetText(setter, key, value)
		},
		service.wrapped.RestoreTextFunc(key),
		restoreFunc)
}

// Sound returns the audio value of the identified text resource.
// In case the text resource has no audio, an empty sound will be returned.
func (service AugmentedTextService) Sound(key resource.Key) audio.L8 {
	return service.wrapped.Sound(key)
}

// RequestSetSound queues the change to update the sound.
func (service AugmentedTextService) RequestSetSound(key resource.Key, sound audio.L8, restoreFunc func()) {
	service.requestCommand(
		func(setter edit.AugmentedTextBlockSetter) {
			service.wrapped.SetSound(setter, key, sound)
		},
		service.wrapped.RestoreSoundFunc(key),
		restoreFunc)
}

// RequestClear queues the change to set both the text and the sound empty.
func (service AugmentedTextService) RequestClear(key resource.Key, restoreFunc func()) {
	service.requestCommand(
		func(setter edit.AugmentedTextBlockSetter) {
			service.wrapped.Clear(setter, key)
		},
		service.wrapped.RestoreFunc(key),
		restoreFunc)
}

// RequestRemove queues the change to remove both the text and the sound from the storage.
func (service AugmentedTextService) RequestRemove(key resource.Key, restoreFunc func()) {
	service.requestCommand(
		func(setter edit.AugmentedTextBlockSetter) {
			service.wrapped.Remove(setter, key)
		},
		service.wrapped.RestoreFunc(key),
		restoreFunc)
}

func (service AugmentedTextService) requestCommand(
	forward func(trans edit.AugmentedTextBlockSetter),
	reverse func(trans edit.AugmentedTextBlockSetter),
	restore func()) {
	c := command{
		forward: func(trans cmd.Transaction) { forward(trans) },
		reverse: func(trans cmd.Transaction) { reverse(trans) },
		restore: restore,
	}
	service.commander.Queue(c)
}
