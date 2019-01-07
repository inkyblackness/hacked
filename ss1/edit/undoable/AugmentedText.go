package undoable

import (
	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/edit"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
)

type AugmentedTextService struct {
	wrapped   edit.AugmentedTextService
	commander cmd.Commander
}

func NewAugmentedTextService(wrapped edit.AugmentedTextService, commander cmd.Commander) AugmentedTextService {
	return AugmentedTextService{
		wrapped:   wrapped,
		commander: commander,
	}
}

func (service AugmentedTextService) IsTrapMessage(key resource.Key) bool {
	return service.wrapped.IsTrapMessage(key)
}

func (service AugmentedTextService) GetText(key resource.Key) string {
	return service.wrapped.GetText(key)
}

func (service AugmentedTextService) RequestSetText(key resource.Key, value string, restoreFunc func()) {
	service.requestCommand(
		func(setter edit.AugmentedTextBlockSetter) {
			service.wrapped.SetText(setter, key, value)
		},
		service.wrapped.RestoreTextFunc(key),
		restoreFunc)
}

func (service AugmentedTextService) GetSound(key resource.Key) audio.L8 {
	return service.wrapped.GetSound(key)
}

func (service AugmentedTextService) RequestSetSound(key resource.Key, sound audio.L8, restoreFunc func()) {
	service.requestCommand(
		func(setter edit.AugmentedTextBlockSetter) {
			service.wrapped.SetSound(setter, key, sound)
		},
		service.wrapped.RestoreSoundFunc(key),
		restoreFunc)
}

func (service AugmentedTextService) RequestClear(key resource.Key, restoreFunc func()) {
	service.requestCommand(
		func(setter edit.AugmentedTextBlockSetter) {
			service.wrapped.Clear(setter, key)
		},
		service.wrapped.RestoreFunc(key),
		restoreFunc)
}

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
