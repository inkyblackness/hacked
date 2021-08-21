package gui

import "github.com/inkyblackness/imgui-go/v3"

// ImageType is an identifier how an image is to be rendered.
type ImageType byte

const (
	// ImageTypeSimpleTexture identifies textures that use the image key as OpenGL texture handle.
	ImageTypeSimpleTexture ImageType = 0
	// ImageTypeBitmapTexture identifies bitmap textures.
	ImageTypeBitmapTexture ImageType = 1
	// ImageTypeColorTexture identifies textures with own color information.
	ImageTypeColorTexture ImageType = 2
)

// TextureIDForSimpleTexture returns a TextureID with ImageTypeSimpleTexture.
func TextureIDForSimpleTexture(handle uint32) imgui.TextureID {
	return imgui.TextureID(ImageTypeSimpleTexture)<<56 | imgui.TextureID(handle)
}

// TextureIDForColorTexture returns a TextureID with ImageTypeSimpleTexture.
func TextureIDForColorTexture(handle uint32) imgui.TextureID {
	return imgui.TextureID(ImageTypeColorTexture)<<56 | imgui.TextureID(handle)
}

// ImageTypeFromID returns the image type the given texture ID specifies.
func ImageTypeFromID(id imgui.TextureID) ImageType {
	return ImageType(id >> 56)
}
