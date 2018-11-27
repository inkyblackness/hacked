package cyber

import (
	"github.com/inkyblackness/hacked/ss1/cyber/media"
	"github.com/inkyblackness/hacked/ss1/resource"
)

type TrapMessageBlockSetter interface {
	SetResourceBlock(lang resource.Language, id resource.ID, index int, data []byte)
	SetResourceBlocks(lang resource.Language, id resource.ID, data [][]byte)
	DelResource(lang resource.Language, id resource.ID)
}

type SetTrapMessageService struct {
	setText  media.SetTextService
	setAudio media.SetAudioService
}

func NewSetTrapMessageService(setText media.SetTextService, setAudio media.SetAudioService) SetTrapMessageService {
	return SetTrapMessageService{
		setText:  setText,
		setAudio: setAudio,
	}
}

func (service SetTrapMessageService) Remove(setter TrapMessageBlockSetter, key resource.Key) {
	service.setText.Remove(setter, key)
	service.setAudio.Remove(setter, trapMessageSoundKeyFor(key))
}

func (service SetTrapMessageService) Clear(setter TrapMessageBlockSetter, key resource.Key) {
	service.setText.Clear(setter, key)
	service.setAudio.Clear(setter, trapMessageSoundKeyFor(key))
}

func (service SetTrapMessageService) Set(setter TrapMessageBlockSetter, key resource.Key, value TrapMessage) {
	service.setText.Set(setter, key, value.Text)
	service.setAudio.Set(setter, trapMessageSoundKeyFor(key), value.Sound)
}
