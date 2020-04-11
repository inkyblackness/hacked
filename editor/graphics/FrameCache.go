package graphics

import (
	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ui/opengl"
)

// FrameCacheKey is used to identify a cached frame.
type FrameCacheKey uint16

// FrameCache loads bitmaps with their palettes and provides OpenGL textures.
type FrameCache struct {
	gl opengl.OpenGL

	textures map[FrameCacheKey]*BitmapTexture
	palettes map[FrameCacheKey]*PaletteTexture

	keyCounter uint16
}

// NewFrameCache returns a new instance.
func NewFrameCache(gl opengl.OpenGL) *FrameCache {
	cache := &FrameCache{
		gl:       gl,
		textures: make(map[FrameCacheKey]*BitmapTexture),
		palettes: make(map[FrameCacheKey]*PaletteTexture),
	}
	return cache
}

// AllocateKey returns a new key to be used with uploading a new frame.
func (cache *FrameCache) AllocateKey() FrameCacheKey {
	cache.keyCounter++
	return FrameCacheKey(cache.keyCounter)
}

// SetTexture registers a texture based on given bitmap under given key.
// The bitmap should contain a palette, otherwise it will most likely not be displayed.
func (cache *FrameCache) SetTexture(key FrameCacheKey, width, height uint16, pixels []byte, palette *bitmap.Palette) {
	cache.DropTextureForKey(key)
	tex := NewBitmapTexture(cache.gl, int(width), int(height), pixels)
	cache.textures[key] = tex
	if palette == nil {
		return
	}
	pal := NewPaletteTexture(cache.gl, palette)
	cache.palettes[key] = pal
}

// DropTextureForKey removes the currently cached frame.
func (cache *FrameCache) DropTextureForKey(key FrameCacheKey) {
	if tex, existing := cache.textures[key]; existing {
		tex.Dispose()
		delete(cache.textures, key)
	}
	if pal, existing := cache.palettes[key]; existing {
		pal.Dispose()
		delete(cache.palettes, key)
	}
}

// HandlesForKey returns the OpenGL handles for both palette and texture for given key.
func (cache FrameCache) HandlesForKey(key FrameCacheKey) (palette uint32, texture uint32) {
	if tex, existing := cache.textures[key]; existing {
		texture = tex.Handle()
	}
	if pal, existing := cache.palettes[key]; existing {
		palette = pal.Handle()
	}
	return
}

// Texture returns the buffered texture for given key.
func (cache *FrameCache) Texture(key FrameCacheKey) *BitmapTexture {
	return cache.textures[key]
}
