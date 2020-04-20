package media

import (
	"bytes"

	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/content/audio/voc"
	"github.com/inkyblackness/hacked/ss1/resource"
)

// SoundEffectBlockSetter modifies storage of raw resource data.
type SoundEffectBlockSetter interface {
	SetResourceBlocks(lang resource.Language, id resource.ID, data [][]byte)
	DelResource(lang resource.Language, id resource.ID)
}

// SoundEffectSetterService provides methods to change sound resources.
type SoundEffectSetterService struct{}

// NewSoundSetterService returns a new instance.
func NewSoundSetterService() SoundEffectSetterService {
	return SoundEffectSetterService{}
}

// Remove deletes any audio resource for given key.
func (service SoundEffectSetterService) Remove(setter SoundEffectBlockSetter, key resource.Key) {
	setter.DelResource(key.Lang, key.ID.Plus(key.Index))
}

// Clear resets the identified audio resource to a silent one-sample audio.
func (service SoundEffectSetterService) Clear(setter SoundEffectBlockSetter, key resource.Key) {
	silence := audio.L8{SampleRate: 22050, Samples: []byte{0x80}}
	service.Set(setter, key, silence)
}

// Set stores the given sound as the identified resource.
func (service SoundEffectSetterService) Set(setter SoundEffectBlockSetter, key resource.Key, data audio.L8) {
	buf := bytes.NewBuffer(nil)
	_ = voc.Save(buf, data.SampleRate, data.Samples)
	blockData := [][]byte{buf.Bytes()}
	setter.SetResourceBlocks(key.Lang, key.ID.Plus(key.Index), blockData)
}
