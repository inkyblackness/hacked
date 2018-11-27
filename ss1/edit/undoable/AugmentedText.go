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

func (service AugmentedTextService) GetText(key resource.Key) string {
	return service.wrapped.GetText(key)
}

func (service AugmentedTextService) GetSound(key resource.Key) (sound audio.L8) {
	return service.wrapped.GetSound(key)
}

func (service AugmentedTextService) RequestSetText(key resource.Key, value string, restoreFunc func()) {
	service.requestCommand(
		func(setter edit.AugmentedTextBlockSetter) {
			service.wrapped.SetText(setter, key, value)
		},
		service.wrapped.RestoreTextFunc(key), // TODO move RestoreTextFunc to here
		restoreFunc)
}

func (service AugmentedTextService) requestCommand(
	forward func(trans edit.AugmentedTextBlockSetter),
	backward func(trans edit.AugmentedTextBlockSetter),
	restore func()) {
	c := command{
		forward:  func(trans cmd.Transaction) { forward(trans) },
		backward: func(trans cmd.Transaction) { backward(trans) },
		restore:  restore,
	}
	service.commander.Queue(c)
}
