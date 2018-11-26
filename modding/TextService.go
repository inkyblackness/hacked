package modding

import (
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

// TODO: type for byte arrays? -> ss1/resource

// TextBlockGetter provides raw data of blocks.
type TextBlockGetter interface {
	ModifiedBlock(lang resource.Language, id resource.ID, index int) []byte
	ModifiedBlocks(lang resource.Language, id resource.ID) [][]byte
}

// TextBlockSetter modifies storage of raw text data
type TextBlockSetter interface {
	SetResourceBlock(lang resource.Language, id resource.ID, index int, data []byte)
	SetResourceBlocks(lang resource.Language, id resource.ID, data [][]byte)
	DelResource(lang resource.Language, id resource.ID)
}

// TextService is modding text.
type TextService struct {
	lineCache *text.Cache
	pageCache *text.Cache
	getter    TextBlockGetter
	cp        text.Codepage
}

func NewTextService(lineCache, pageCache *text.Cache, getter TextBlockGetter, cp text.Codepage) TextService {
	return TextService{
		lineCache: lineCache,
		pageCache: pageCache,
		getter:    getter,
		cp:        cp,
	}
}

func (service TextService) Current(key resource.Key) string {
	var cache *text.Cache
	resourceInfo, existing := ids.Info(key.ID)
	if !existing || resourceInfo.List {
		cache = service.lineCache
	} else {
		cache = service.pageCache
	}
	currentValue, cacheErr := cache.Text(key)
	if cacheErr != nil {
		currentValue = ""
	}
	return currentValue
}

func (service TextService) Remove(setter TextBlockSetter, key resource.Key) {
	info, _ := ids.Info(key.ID)
	if info.List {
		setter.SetResourceBlock(key.Lang, key.ID, key.Index, nil)
	} else {
		id := key.ID.Plus(key.Index)
		setter.DelResource(key.Lang, id)
	}
}

func (service TextService) Clear(setter TextBlockSetter, key resource.Key) {
	service.Set(setter, key, "")
}

func (service TextService) Set(setter TextBlockSetter, key resource.Key, value string) {
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
