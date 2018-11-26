package modding

import (
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

// TextBlockSetter modifies storage of raw text data
type TextBlockSetter interface {
	SetResourceBlock(lang resource.Language, id resource.ID, index int, data []byte)
	SetResourceBlocks(lang resource.Language, id resource.ID, data [][]byte)
	DelResource(lang resource.Language, id resource.ID)
}

// TextService is modding text.
type SetTextService struct {
	cp text.Codepage
}

func NewSetTextService(cp text.Codepage) SetTextService {
	return SetTextService{
		cp: cp,
	}
}

func (service SetTextService) Remove(setter TextBlockSetter, key resource.Key) {
	info, _ := ids.Info(key.ID)
	if info.List {
		setter.SetResourceBlock(key.Lang, key.ID, key.Index, nil)
	} else {
		id := key.ID.Plus(key.Index)
		setter.DelResource(key.Lang, id)
	}
}

func (service SetTextService) Clear(setter TextBlockSetter, key resource.Key) {
	service.Set(setter, key, "")
}

func (service SetTextService) Set(setter TextBlockSetter, key resource.Key, value string) {
	blockedValue := text.Blocked(value)
	info, _ := ids.Info(key.ID)
	if info.List {
		newData := service.cp.Encode(blockedValue[0])
		setter.SetResourceBlock(key.Lang, key.ID, key.Index, newData)
	} else {
		newData := make([][]byte, len(blockedValue))
		for index, blockLine := range blockedValue {
			newData[index] = service.cp.Encode(blockLine)
		}
		id := key.ID.Plus(key.Index)
		setter.SetResourceBlocks(key.Lang, id, newData)
	}
}
