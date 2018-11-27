package cyber

import (
	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/cyber/media"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

type TrapMessage struct {
	Text  string
	Sound audio.L8
}

type GetTrapMessageService struct {
	getText  media.GetTextService
	getAudio media.GetAudioService
}

func NewGetTrapMessageService(getText media.GetTextService, getAudio media.GetAudioService) GetTrapMessageService {
	return GetTrapMessageService{
		getText:  getText,
		getAudio: getAudio,
	}
}

func (service GetTrapMessageService) Current(key resource.Key) TrapMessage {
	text := service.getText.Current(key)
	sound := service.getAudio.Get(trapMessageSoundKeyFor(key))
	return TrapMessage{text, sound}
}

func trapMessageSoundKeyFor(key resource.Key) resource.Key {
	soundKey := key
	soundKey.ID = ids.TrapMessagesAudioStart.Plus(key.Index)
	soundKey.Index = 0
	return soundKey
}

func (service GetTrapMessageService) Modified(key resource.Key) bool {
	return service.getText.Modified(key) ||
		service.getAudio.Modified(trapMessageSoundKeyFor(key))
}
