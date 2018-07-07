package graphics

import (
	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ui/opengl"
)

// PaletteTexture contains a palette stored as OpenGL texture.
type PaletteTexture struct {
	gl opengl.OpenGL

	handle uint32
}

// NewPaletteTexture creates a new PaletteTexture instance.
func NewPaletteTexture(gl opengl.OpenGL, palette bitmap.Palette) *PaletteTexture {
	tex := &PaletteTexture{
		gl:     gl,
		handle: gl.GenTextures(1)[0]}

	tex.Update(palette)

	return tex
}

// Dispose releases the OpenGL texture.
func (tex *PaletteTexture) Dispose() {
	if tex.handle != 0 {
		tex.gl.DeleteTextures([]uint32{tex.handle})
		tex.handle = 0
	}
}

// Handle returns the texture handle.
func (tex *PaletteTexture) Handle() uint32 {
	return tex.handle
}

// Update reloads the palette.
func (tex *PaletteTexture) Update(palette bitmap.Palette) {
	gl := tex.gl
	const bytesPerRGBA = 4
	const colors = 256
	var data [colors * bytesPerRGBA]byte

	for i := 0; i < colors; i++ {
		entry := palette[i]
		data[i*bytesPerRGBA+0] = entry.Red
		data[i*bytesPerRGBA+1] = entry.Green
		data[i*bytesPerRGBA+2] = entry.Blue
		data[i*bytesPerRGBA+3] = 0xFF
	}
	data[3] = 0x00

	gl.BindTexture(opengl.TEXTURE_2D, tex.handle)
	gl.TexImage2D(opengl.TEXTURE_2D, 0, opengl.RGBA, int32(colors), 1, 0, opengl.RGBA, opengl.UNSIGNED_BYTE, palette[:])
	gl.TexParameteri(opengl.TEXTURE_2D, opengl.TEXTURE_MAG_FILTER, opengl.NEAREST)
	gl.TexParameteri(opengl.TEXTURE_2D, opengl.TEXTURE_MIN_FILTER, opengl.NEAREST)
	gl.GenerateMipmap(opengl.TEXTURE_2D)
	gl.BindTexture(opengl.TEXTURE_2D, 0)
}
