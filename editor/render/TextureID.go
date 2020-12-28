package render

import (
	"github.com/inkyblackness/imgui-go/v3"

	"github.com/inkyblackness/hacked/editor/graphics"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ui/gui"
)

const (
	// BitmapTextureTypeResource specifies a bitmap texture directly from resources.
	BitmapTextureTypeResource = 0x00
	// BitmapTextureTypeFrame specifies a bitmap texture from the frame cache.
	BitmapTextureTypeFrame = 0x01
)

// TextureIDForBitmapTexture returns a texture ID that identifies a bitmap texture from resources.
func TextureIDForBitmapTexture(key resource.Key) imgui.TextureID {
	id := imgui.TextureID(gui.ImageTypeBitmapTexture) << 56
	id |= imgui.TextureID(BitmapTextureTypeResource) << 48
	id |= imgui.TextureID(key.Lang) << 32
	id |= imgui.TextureID(key.ID) << 16
	id |= imgui.TextureID(key.Index & 0xFFFF)
	return id
}

// TextureIDForBitmapFrame returns a texture ID that identifies a bitmap texture from frame cache.
func TextureIDForBitmapFrame(key graphics.FrameCacheKey) imgui.TextureID {
	id := imgui.TextureID(gui.ImageTypeBitmapTexture) << 56
	id |= imgui.TextureID(BitmapTextureTypeFrame) << 48
	id |= imgui.TextureID(key)
	return id
}
