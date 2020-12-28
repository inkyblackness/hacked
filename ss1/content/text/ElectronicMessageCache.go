package text

import (
	"github.com/inkyblackness/hacked/ss1/resource"
)

// ElectronicMessageCache retrieves messages from a localizer and keeps them decoded until they are invalidated.
type ElectronicMessageCache struct {
	cp        Codepage
	localizer resource.Localizer

	messages map[resource.Key]ElectronicMessage
}

// NewElectronicMessageCache returns a new instance.
func NewElectronicMessageCache(cp Codepage, localizer resource.Localizer) *ElectronicMessageCache {
	cache := &ElectronicMessageCache{
		cp:        cp,
		localizer: localizer,

		messages: make(map[resource.Key]ElectronicMessage),
	}
	return cache
}

// InvalidateResources lets the cache remove any texts from resources that are specified in the given slice.
func (cache *ElectronicMessageCache) InvalidateResources(ids []resource.ID) {
	for _, id := range ids {
		for key := range cache.messages {
			if key.ID == id {
				delete(cache.messages, key)
			}
		}
	}
}

// Message retrieves and caches the message of given key.
func (cache *ElectronicMessageCache) Message(key resource.Key) (ElectronicMessage, error) {
	cacheKey := resource.KeyOf(key.ID.Plus(key.Index), key.Lang, 0)
	value, existing := cache.messages[cacheKey]
	if existing {
		return value, nil
	}
	selector := cache.localizer.LocalizedResources(key.Lang)
	view, err := selector.Select(cacheKey.ID)
	if err != nil {
		return EmptyElectronicMessage(), err
	}
	if (view.ContentType() != resource.Text) || !view.Compound() {
		return EmptyElectronicMessage(), resource.ErrWrongType(key, resource.Text)
	}
	value, err = DecodeElectronicMessage(cache.cp, view)
	if err != nil {
		return EmptyElectronicMessage(), err
	}
	cache.messages[cacheKey] = value
	return value, nil
}
