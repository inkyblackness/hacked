package graphics

import (
	"errors"

	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ui/opengl"
)

// TextureCache loads bitmaps and provides OpenGL textures.
type TextureCache struct {
	gl        opengl.OpenGL
	localizer resource.Localizer

	textures map[resource.Key]*BitmapTexture
}

// NewTextureCache returns a new instance.
func NewTextureCache(gl opengl.OpenGL, localizer resource.Localizer) *TextureCache {
	cache := &TextureCache{
		gl:        gl,
		localizer: localizer,
		textures:  make(map[resource.Key]*BitmapTexture),
	}
	return cache
}

// InvalidateResources lets the cache remove any textures from resources that are specified in the given slice.
func (cache *TextureCache) InvalidateResources(ids []resource.ID) {
	for _, id := range ids {
		for key, texture := range cache.textures {
			if key.ID == id {
				texture.Dispose()
				delete(cache.textures, key)
			}
		}
	}
}

// Texture returns the texture with given key - if available.
func (cache *TextureCache) Texture(key resource.Key) (*BitmapTexture, error) {
	tex, existing := cache.textures[key]
	if existing {
		return tex, nil
	}
	selector := cache.localizer.LocalizedResources(key.Lang)
	view, err := selector.Select(key.ID)
	if err != nil {
		return nil, err
	}
	if view.ContentType() != resource.Bitmap {
		return nil, errors.New("resource not a bitmap")
	}
	reader, err := view.Block(key.Index)
	if err != nil {
		return nil, err
	}
	bmp, err := bitmap.Decode(reader)
	if err != nil {
		return nil, err
	}

	tex = NewBitmapTexture(cache.gl, int(bmp.Header.Width), int(bmp.Header.Height), bmp.Pixels)
	cache.textures[key] = tex

	return tex, nil
}
