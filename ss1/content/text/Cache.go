package text

import (
	"errors"
	"io/ioutil"

	"github.com/inkyblackness/hacked/ss1/resource"
)

type keyResolver func(resource.Key) resource.Key
type textReader func(resource.Selector, resource.Key, Codepage) (string, error)

// Cache retrieves texts from a localizer and keeps them decoded until they are invalidated.
type Cache struct {
	cp        Codepage
	localizer resource.Localizer
	reader    textReader

	keyResolver keyResolver
	texts       map[resource.Key]string
}

func newCache(cp Codepage, localizer resource.Localizer, keyResolver keyResolver, reader textReader) *Cache {
	cache := &Cache{
		cp:        cp,
		localizer: localizer,
		reader:    reader,

		keyResolver: keyResolver,
		texts:       make(map[resource.Key]string),
	}
	return cache
}

// NewLineCache returns a cache for single-block texts.
func NewLineCache(cp Codepage, localizer resource.Localizer) *Cache {
	return newCache(cp, localizer, func(key resource.Key) resource.Key { return key }, readLine)
}

// NewPageCache returns a cache for resource-based texts.
func NewPageCache(cp Codepage, localizer resource.Localizer) *Cache {
	return newCache(cp, localizer, func(key resource.Key) resource.Key {
		return resource.KeyOf(key.ID.Plus(key.Index), key.Lang, 0)
	}, readPage)
}

func readLine(selector resource.Selector, key resource.Key, cp Codepage) (string, error) {
	view, err := selector.Select(key.ID)
	if err != nil {
		return "", err
	}
	if view.ContentType() != resource.Text {
		return "", errors.New("resource is not a text")
	}
	reader, err := view.Block(key.Index)
	if err != nil {
		return "", err
	}
	raw, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return cp.Decode(raw), nil
}

func readPage(selector resource.Selector, key resource.Key, cp Codepage) (string, error) {
	view, err := selector.Select(key.ID.Plus(key.Index))
	if err != nil {
		return "", err
	}
	if view.ContentType() != resource.Text {
		return "", errors.New("resource is not a text")
	}
	blockCount := view.BlockCount()
	value := ""
	for block := 0; block < blockCount; block++ {
		reader, err := view.Block(block)
		if err != nil {
			return "", err
		}
		raw, err := ioutil.ReadAll(reader)
		if err != nil {
			return "", err
		}
		value += cp.Decode(raw)
	}
	return value, nil
}

// InvalidateResources lets the cache remove any texts from resources that are specified in the given slice.
func (cache *Cache) InvalidateResources(ids []resource.ID) {
	for _, id := range ids {
		for key := range cache.texts {
			if key.ID == id {
				delete(cache.texts, key)
			}
		}
	}
}

// Text retrieves and caches the text of given key.
func (cache *Cache) Text(key resource.Key) (string, error) {
	cacheKey := cache.keyResolver(key)
	value, existing := cache.texts[cacheKey]
	if existing {
		return value, nil
	}
	selector := cache.localizer.LocalizedResources(key.Lang)
	value, err := cache.reader(selector, key, cache.cp)
	if err != nil {
		return "", err
	}
	cache.texts[cacheKey] = value
	return value, nil
}
