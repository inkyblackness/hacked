package media

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

// TextViewerService provides read-only access to text resources.
type TextViewerService struct {
	lineCache *text.Cache
	pageCache *text.Cache
	getter    TextBlockGetter
}

// NewTextViewerService returns a new instance.
func NewTextViewerService(lineCache, pageCache *text.Cache, getter TextBlockGetter) TextViewerService {
	return TextViewerService{
		lineCache: lineCache,
		pageCache: pageCache,
		getter:    getter,
	}
}

// Text returns the text data associated with the given key.
func (service TextViewerService) Text(key resource.Key) string {
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

// Modified returns true if the identified text resource is marked as modified.
func (service TextViewerService) Modified(key resource.Key) bool {
	info, _ := ids.Info(key.ID)
	if info.List {
		return len(service.getter.ModifiedBlock(key.Lang, key.ID, key.Index)) > 0
	}
	return len(service.getter.ModifiedBlocks(key.Lang, key.ID.Plus(key.Index))) > 0
}
