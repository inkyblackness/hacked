package texts

import (
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

type BlockKeeper interface {
	ModifiedBlock(lang resource.Language, id resource.ID, index int) []byte
	ModifiedBlocks(lang resource.Language, id resource.ID) [][]byte
}

type BlockSetter interface {
	SetResourceBlock(lang resource.Language, id resource.ID, index int, data []byte)
	SetResourceBlocks(lang resource.Language, id resource.ID, data [][]byte)
	DelResource(lang resource.Language, id resource.ID)
}

type Adapter struct {
	lineCache *text.Cache
	pageCache *text.Cache
	mod       BlockKeeper
	cp        text.Codepage
}

func NewAdapter(lineCache, pageCache *text.Cache, keeper BlockKeeper, cp text.Codepage) Adapter {
	return Adapter{
		lineCache: lineCache,
		pageCache: pageCache,
		mod:       keeper,
		cp:        cp,
	}
}

func (adapter Adapter) CurrentText(key resource.Key) string {
	var cache *text.Cache
	resourceInfo, existing := ids.Info(key.ID)
	if !existing || resourceInfo.List {
		cache = adapter.lineCache
	} else {
		cache = adapter.pageCache
	}
	currentValue, cacheErr := cache.Text(key)
	if cacheErr != nil {
		currentValue = ""
	}
	return currentValue
}

func (adapter Adapter) Modification(key resource.Key) (data [][]byte, isList bool) {
	info, _ := ids.Info(key.ID)
	if info.List {
		data = [][]byte{adapter.mod.ModifiedBlock(key.Lang, key.ID, key.Index)}
	} else {
		data = adapter.mod.ModifiedBlocks(key.Lang, key.ID.Plus(key.Index))
	}
	return data, info.List
}

func (adapter Adapter) RemoveText(modifier BlockSetter, key resource.Key) {
	info, _ := ids.Info(key.ID)
	if info.List {
		modifier.SetResourceBlock(key.Lang, key.ID, key.Index, nil)
	} else {
		id := key.ID.Plus(key.Index)
		modifier.DelResource(key.Lang, id)
	}
}

func (adapter Adapter) SetText(modifier BlockSetter, key resource.Key, value string) {
	blockedValue := text.Blocked(value)
	info, _ := ids.Info(key.ID)
	if info.List {
		newData := adapter.cp.Encode(blockedValue[0])
		modifier.SetResourceBlock(key.Lang, key.ID, key.Index, newData)
	} else {
		newData := make([][]byte, len(blockedValue))
		for index, blockLine := range blockedValue {
			newData[index] = adapter.cp.Encode(blockLine)
		}
		id := key.ID.Plus(key.Index)
		modifier.SetResourceBlocks(key.Lang, id, newData)
	}
}
