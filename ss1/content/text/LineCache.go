package text

import (
	"errors"
	"io/ioutil"

	"github.com/inkyblackness/hacked/ss1/resource"
)

// LineCache is a cache for keeping single-block text lines.
type LineCache struct {
	// Codepage is used to decode the strings.
	Codepage  Codepage
	localizer resource.Localizer

	lines map[resource.Key]string
}

// NewLineCache returns a new instance.
func NewLineCache(cp Codepage, localizer resource.Localizer) *LineCache {
	cache := &LineCache{
		Codepage:  cp,
		localizer: localizer,

		lines: make(map[resource.Key]string),
	}
	return cache
}

// InvalidateResources lets the cache remove any lines from resources that are specified in the given slice.
func (cache *LineCache) InvalidateResources(ids []resource.ID) {
	for _, id := range ids {
		for key := range cache.lines {
			if key.ID == id {
				delete(cache.lines, key)
			}
		}
	}
}

// Line retrieves and caches the text line of given key.
func (cache *LineCache) Line(key resource.Key) (string, error) {
	line, existing := cache.lines[key]
	if existing {
		return line, nil
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
	line = cache.Codepage.Decode(raw)
	cache.lines[key] = line
	return line, nil
}
