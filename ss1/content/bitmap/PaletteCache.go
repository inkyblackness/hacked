package bitmap

import (
	"errors"

	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/serial"
)

// PaletteCache retrieves palettes from a localizer and keeps them decoded until they are invalidated.
type PaletteCache struct {
	localizer resource.Localizer

	defaultPalette Palette
	palettes       map[resource.Key]Palette
}

// NewPaletteCache returns a new instance.
func NewPaletteCache(localizer resource.Localizer) *PaletteCache {
	cache := &PaletteCache{
		localizer: localizer,
		palettes:  make(map[resource.Key]Palette),
	}
	for i := 0; i < len(cache.defaultPalette); i++ {
		col := &cache.defaultPalette[i]
		col.Red = uint8(i)
		col.Green = uint8(i)
		col.Blue = uint8(i)
	}
	return cache
}

// InvalidateResources lets the cache remove any palettes from resources that are specified in the given slice.
func (cache *PaletteCache) InvalidateResources(ids []resource.ID) {
	for _, id := range ids {
		for key := range cache.palettes {
			if key.ID == id {
				delete(cache.palettes, key)
			}
		}
	}
}

// Palette tries to look up given palette.
func (cache *PaletteCache) Palette(key resource.Key) (pal Palette, err error) {
	pal, existing := cache.palettes[key]
	if existing {
		return pal, nil
	}
	selector := cache.localizer.LocalizedResources(key.Lang)
	view, err := selector.Select(key.ID)
	if err != nil {
		return cache.defaultPalette, err
	}
	if (view.ContentType() != resource.Palette) || (view.BlockCount() != 1) {
		return cache.defaultPalette, errors.New("resource is not a palette")
	}
	reader, err := view.Block(0)
	if err != nil {
		return cache.defaultPalette, err
	}
	decoder := serial.NewDecoder(reader)
	decoder.Code(&pal)
	if decoder.FirstError() != nil {
		return cache.defaultPalette, decoder.FirstError()
	}
	cache.palettes[key] = pal
	return pal, nil
}
