package text

import (
	"errors"
	"io/ioutil"

	"github.com/inkyblackness/hacked/ss1/resource"
)

type textReader func(resource.Selector, resource.Key, Codepage) (string, error)

// Cache retrieves texts from a localizer and keeps them decoded until they are invalidated.
type Cache struct {
	cp        Codepage
	localizer resource.Localizer
	reader    textReader

	texts map[resource.Key]string
}

func newCache(cp Codepage, localizer resource.Localizer, reader textReader) *Cache {
	cache := &Cache{
		cp:        cp,
		localizer: localizer,
		reader:    reader,

		texts: make(map[resource.Key]string),
	}
	return cache
}

func NewLineCache(cp Codepage, localizer resource.Localizer) *Cache {
	return newCache(cp, localizer, readLine)
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
	value, existing := cache.texts[key]
	if existing {
		return value, nil
	}
	selector := cache.localizer.LocalizedResources(key.Lang)
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
	value = cache.cp.Decode(raw)
	cache.texts[key] = value
	return value, nil
}
