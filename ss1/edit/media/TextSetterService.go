package media

import (
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

// TextBlockSetter modifies storage of raw resource data.
type TextBlockSetter interface {
	SetResourceBlock(lang resource.Language, id resource.ID, index int, data []byte)
	SetResourceBlocks(lang resource.Language, id resource.ID, data [][]byte)
	DelResource(lang resource.Language, id resource.ID)
}

// TextSetterService provides methods to change text resources.
type TextSetterService struct {
	cp text.Codepage
}

// NewTextSetterService returns a new instance.
func NewTextSetterService(cp text.Codepage) TextSetterService {
	return TextSetterService{
		cp: cp,
	}
}

// Remove deletes any text resource for given key.
func (service TextSetterService) Remove(setter TextBlockSetter, key resource.Key) {
	info, _ := ids.Info(key.ID)
	if info.List {
		setter.SetResourceBlock(key.Lang, key.ID, key.Index, nil)
	} else {
		id := key.ID.Plus(key.Index)
		setter.DelResource(key.Lang, id)
	}
}

// Clear resets the identified text resource to an empty string.
func (service TextSetterService) Clear(setter TextBlockSetter, key resource.Key) {
	service.Set(setter, key, "")
}

// Set stores the given text as the identified resource.
func (service TextSetterService) Set(setter TextBlockSetter, key resource.Key, value string) {
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
