package graphics

import "github.com/inkyblackness/hacked/ui/opengl"

// BitmapTexture wraps an OpenGL handle for a downloaded image.
type BitmapTexture struct {
	gl     opengl.OpenGL
	handle uint32

	width, height float32
	u, v          float32
}

func powerOfTwo(value int) int {
	result := 2

	for (result < value) && (result < 0x1000) {
		result *= 2
	}

	return result
}

// NewBitmapTexture downloads the provided raw data to OpenGL and returns a BitmapTexture instance.
func NewBitmapTexture(gl opengl.OpenGL, width, height int, pixelData []byte) *BitmapTexture {
	textureWidth := powerOfTwo(width)
	textureHeight := powerOfTwo(height)
	tex := &BitmapTexture{
		gl:     gl,
		width:  float32(width),
		height: float32(height),
		handle: gl.GenTextures(1)[0]}
	tex.u = tex.width / float32(textureWidth)
	tex.v = tex.height / float32(textureHeight)

	gl.BindTexture(opengl.TEXTURE_2D, tex.handle)
	gl.TexImage2D(opengl.TEXTURE_2D, 0, opengl.RED, int32(textureWidth), int32(textureHeight),
		0, opengl.RED, opengl.UNSIGNED_BYTE, pixelData)
	gl.TexParameteri(opengl.TEXTURE_2D, opengl.TEXTURE_MAG_FILTER, opengl.NEAREST)
	gl.TexParameteri(opengl.TEXTURE_2D, opengl.TEXTURE_MIN_FILTER, opengl.NEAREST)
	gl.GenerateMipmap(opengl.TEXTURE_2D)
	gl.BindTexture(opengl.TEXTURE_2D, 0)

	return tex
}

// Dispose releases the OpenGL texture.
func (tex *BitmapTexture) Dispose() {
	if tex.handle != 0 {
		tex.gl.DeleteTextures([]uint32{tex.handle})
		tex.handle = 0
	}
}

// Handle returns the texture handle.
func (tex *BitmapTexture) Handle() uint32 {
	return tex.handle
}

// Size returns the dimensions of the bitmap, in pixels.
func (tex *BitmapTexture) Size() (width, height float32) {
	return tex.width, tex.height
}

// UV returns the maximum U and V values for the bitmap. The bitmap will be
// stored in a power-of-two texture, which may be larger than the bitmap.
func (tex *BitmapTexture) UV() (u, v float32) {
	return tex.u, tex.v
}
