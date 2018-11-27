package cyber

import (
	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/cyber/media"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

type AugmentedTextBlockSetter interface {
	SetResourceBlock(lang resource.Language, id resource.ID, index int, data []byte)
	SetResourceBlocks(lang resource.Language, id resource.ID, data [][]byte)
	DelResource(lang resource.Language, id resource.ID)
}

func TrapMessageSoundKeyFor(key resource.Key) resource.Key {
	soundKey := key
	soundKey.ID = ids.TrapMessagesAudioStart.Plus(key.Index)
	soundKey.Index = 0
	return soundKey
}

type AugmentedTextService struct {
	getText  media.GetTextService
	setText  media.SetTextService
	getAudio media.GetAudioService
	setAudio media.SetAudioService
}

func NewAugmentedTextService(
	getText media.GetTextService, setText media.SetTextService,
	getAudio media.GetAudioService, setAudio media.SetAudioService) AugmentedTextService {
	return AugmentedTextService{
		getText:  getText,
		setText:  setText,
		getAudio: getAudio,
		setAudio: setAudio,
	}
}

func (service AugmentedTextService) IsTrapMessage(key resource.Key) bool {
	return key.ID == ids.TrapMessageTexts
}

func (service AugmentedTextService) GetText(key resource.Key) string {
	return service.getText.Current(key)
}

func (service AugmentedTextService) SetText(setter AugmentedTextBlockSetter, key resource.Key, value string) {
	service.setText.Set(setter, key, value)
}

func (service AugmentedTextService) RestoreTextFunc(key resource.Key) func(setter AugmentedTextBlockSetter) {
	oldText := service.getText.Current(key)
	isModified := service.getText.Modified(key)

	return func(setter AugmentedTextBlockSetter) {
		if isModified {
			service.setText.Set(setter, key, oldText)
		} else {
			service.setText.Remove(setter, key)
		}
	}
}

func (service AugmentedTextService) GetSound(key resource.Key) (sound audio.L8) {
	if !service.IsTrapMessage(key) {
		return
	}
	return service.getAudio.Get(TrapMessageSoundKeyFor(key))
}

func (service AugmentedTextService) SetSound(setter AugmentedTextBlockSetter, key resource.Key, sound audio.L8) {
	if service.IsTrapMessage(key) {
		service.setAudio.Set(setter, TrapMessageSoundKeyFor(key), sound)
	}
}

func (service AugmentedTextService) RestoreSoundFunc(key resource.Key) func(setter AugmentedTextBlockSetter) {
	if !service.IsTrapMessage(key) {
		return func(setter AugmentedTextBlockSetter) {}
	}

	soundKey := TrapMessageSoundKeyFor(key)
	isSoundModified := service.getAudio.Modified(soundKey)
	oldSound := service.getAudio.Get(soundKey)

	return func(setter AugmentedTextBlockSetter) {
		if isSoundModified {
			service.setAudio.Set(setter, soundKey, oldSound)
		} else {
			service.setAudio.Remove(setter, soundKey)
		}
	}
}

func (service AugmentedTextService) Clear(setter AugmentedTextBlockSetter, key resource.Key) {
	service.setText.Clear(setter, key)
	if service.IsTrapMessage(key) {
		service.setAudio.Clear(setter, TrapMessageSoundKeyFor(key))
	}
}

func (service AugmentedTextService) Remove(setter AugmentedTextBlockSetter, key resource.Key) {
	service.setText.Remove(setter, key)
	if service.IsTrapMessage(key) {
		service.setAudio.Remove(setter, TrapMessageSoundKeyFor(key))
	}
}

func (service AugmentedTextService) RestoreFunc(key resource.Key) func(setter AugmentedTextBlockSetter) {
	restoreText := service.RestoreTextFunc(key)
	restoreSound := service.RestoreSoundFunc(key)
	return func(setter AugmentedTextBlockSetter) {
		restoreText(setter)
		restoreSound(setter)
	}
}
