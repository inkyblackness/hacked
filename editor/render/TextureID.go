package render

import (
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ui/gui"
	"github.com/inkyblackness/imgui-go"
)

// TextureIDForBitmapTexture returns a texture ID that identifies a bitmap texture.
func TextureIDForBitmapTexture(key resource.Key) imgui.TextureID {
	id := imgui.TextureID(gui.ImageTypeBitmapTexture) << 56
	id |= imgui.TextureID(key.Lang) << 32
	id |= imgui.TextureID(key.ID) << 16
	id |= imgui.TextureID(key.Index & 0xFFFF)
	return id
}
