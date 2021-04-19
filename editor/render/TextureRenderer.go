package render

import (
	"unsafe"

	"github.com/inkyblackness/hacked/ui/opengl"
)

type TextureRenderer struct {
	gl opengl.OpenGL

	framebuffer uint32
	texture     uint32

	width  int32
	height int32
}

func NewTextureRenderer(gl opengl.OpenGL) *TextureRenderer {
	renderer := &TextureRenderer{
		gl: gl,

		framebuffer: gl.GenFramebuffers(1)[0],
		texture:     gl.GenTextures(1)[0],

		width:  400,
		height: 400,
	}

	gl.BindTexture(opengl.TEXTURE_2D, renderer.texture)
	gl.TexParameteri(opengl.TEXTURE_2D, opengl.TEXTURE_MAG_FILTER, opengl.NEAREST)
	gl.TexParameteri(opengl.TEXTURE_2D, opengl.TEXTURE_MIN_FILTER, opengl.NEAREST)
	gl.TexImage2D(opengl.TEXTURE_2D, 0, opengl.RGBA, renderer.width, renderer.height, 0, opengl.RGBA, opengl.UNSIGNED_BYTE, unsafe.Pointer(uintptr(0)))
	// gl.GenerateMipmap(opengl.TEXTURE_2D)
	gl.BindTexture(opengl.TEXTURE_2D, 0)

	renderer.onFramebuffer(func() {
		gl.BindTexture(opengl.TEXTURE_2D, renderer.texture)
		gl.FramebufferTexture(opengl.FRAMEBUFFER, opengl.COLOR_ATTACHMENT0, renderer.texture, 0)
		gl.DrawBuffers([]uint32{opengl.COLOR_ATTACHMENT0})
		gl.BindTexture(opengl.TEXTURE_2D, 0)

		// result := gl.CheckFramebufferStatus(opengl.FRAMEBUFFER)
		// fmt.Printf("status: 0x%04X\n", result)
		// result == FRAMEBUFFER_COMPLETE == 0x8CD5
	})
	return renderer
}

func (renderer *TextureRenderer) Dispose() {
	gl := renderer.gl
	gl.DeleteFramebuffers([]uint32{renderer.framebuffer})
	gl.DeleteTextures([]uint32{renderer.texture})
}

// Handle returns the texture handle.
func (renderer *TextureRenderer) Handle() uint32 {
	return renderer.texture
}

// Size returns the dimensions of the texture, in pixels.
func (renderer *TextureRenderer) Size() (width, height float32) {
	return float32(renderer.width), float32(renderer.height)
}

func (renderer *TextureRenderer) Render(nested func()) {
	gl := renderer.gl
	renderer.onFramebuffer(func() {
		gl.Viewport(0, 0, renderer.width, renderer.height)
		gl.ClearColor(0.0, 1.0, 0.0, 1.0)
		gl.Clear(opengl.COLOR_BUFFER_BIT | opengl.DEPTH_BUFFER_BIT)

		nested()
	})
}

func (renderer *TextureRenderer) onFramebuffer(nested func()) {
	gl := renderer.gl
	gl.BindFramebuffer(opengl.FRAMEBUFFER, renderer.framebuffer)
	nested()
	gl.BindFramebuffer(opengl.FRAMEBUFFER, 0)
}
