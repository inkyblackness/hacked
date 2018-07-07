package graphics

import (
	"encoding/binary"
	"errors"

	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/hacked/ui/opengl"
)

// PaletteCache loads palettes and provides OpenGL textures.
type PaletteCache struct {
	gl        opengl.OpenGL
	localizer resource.Localizer

	palettes map[resource.Key]*PaletteTexture
}

// NewPaletteCache returns a new instance.
func NewPaletteCache(gl opengl.OpenGL, localizer resource.Localizer) *PaletteCache {
	cache := &PaletteCache{
		gl:        gl,
		localizer: localizer,
		palettes:  make(map[resource.Key]*PaletteTexture),
	}
	return cache
}

// InvalidateResources lets the cache remove any palette from resources that are specified in the given slice.
func (cache *PaletteCache) InvalidateResources(ids []resource.ID) {
	for _, id := range ids {
		for key, texture := range cache.palettes {
			if key.ID == id {
				texture.Dispose()
				delete(cache.palettes, key)
			}
		}
	}
}

// Palette returns the palette with given index - if available.
func (cache *PaletteCache) Palette(index int) (*PaletteTexture, error) {
	key := resource.KeyOf(ids.GamePalettesStart.Plus(index), resource.LangAny, 0)
	pal, existing := cache.palettes[key]
	if existing {
		return pal, nil
	}
	selector := cache.localizer.LocalizedResources(key.Lang)
	view, err := selector.Select(key.ID)
	if err != nil {
		return nil, err
	}
	if view.ContentType() != resource.Palette {
		return nil, errors.New("resource not a palette")
	}
	reader, err := view.Block(key.Index)
	if err != nil {
		return nil, err
	}
	var palette bitmap.Palette
	err = binary.Read(reader, binary.LittleEndian, &palette)
	if err != nil {
		return nil, err
	}

	pal = NewPaletteTexture(cache.gl, palette)
	cache.palettes[key] = pal

	return pal, nil
}
