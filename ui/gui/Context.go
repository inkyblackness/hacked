package gui

import (
	"github.com/inkyblackness/hacked/ui/opengl"
	"github.com/inkyblackness/imgui-go"
	"time"
)

// Context describes a scope for a graphical user interface.
// It is based on ImGui.
type Context struct {
	imguiContext *imgui.Context
	window       opengl.Window

	lastRenderTime time.Time
	mouseDown      [3]bool

	fontTexture            uint32
	shaderHandle           uint32
	attribLocationTex      int32
	attribLocationProjMtx  int32
	attribLocationPosition int32
	attribLocationUV       int32
	attribLocationColor    int32
	vboHandle              uint32
	elementsHandle         uint32
}

// NewContext initializes a new UI context based on the provided OpenGL window.
func NewContext(window opengl.Window) (context *Context, err error) {
	context = &Context{
		imguiContext: imgui.CreateContext(nil),
		window:       window,
	}

	err = context.createDeviceObjects()
	if err != nil {
		context.Destroy()
		context = nil
	}

	return
}

// Destroy cleans up the resources of the graphical user interface.
func (context *Context) Destroy() {
	context.destroyDeviceObjects(context.window.OpenGl())
	context.imguiContext.Destroy()
}

// NewFrame must be called at the start of rendering.
func (context *Context) NewFrame() {
	io := imgui.CurrentIO()

	windowWidth, windowHeight := context.window.Size()
	io.SetDisplaySize(imgui.Vec2{X: float32(windowWidth), Y: float32(windowHeight)})

	now := time.Now()
	if !context.lastRenderTime.IsZero() {
		elapsed := now.Sub(context.lastRenderTime)
		io.SetDeltaTime(float32(elapsed.Seconds()))
	}
	context.lastRenderTime = now

	imgui.NewFrame()
}

// Render must be called at the end of rendering.
func (context *Context) Render() {
	imgui.Render()

}

func (context *Context) createDeviceObjects() (err error) {
	gl := context.window.OpenGl()
	glslVersion := "#version 150"

	vertexShaderSource := glslVersion + `
uniform mat4 ProjMtx;
in vec2 Position;
in vec2 UV;
in vec4 Color;
out vec2 Frag_UV;
out vec4 Frag_Color;
void main()
{
	Frag_UV = UV;
	Frag_Color = Color;
	gl_Position = ProjMtx * vec4(Position.xy,0,1);
}
`
	fragmentShaderSource := glslVersion + `
uniform sampler2D Texture;
in vec2 Frag_UV;
in vec4 Frag_Color;
out vec4 Out_Color;
void main()
{
	Out_Color = Frag_Color * texture( Texture, Frag_UV.st).r;
}
`
	context.shaderHandle, err = opengl.LinkNewStandardProgram(gl, vertexShaderSource, fragmentShaderSource)
	if err != nil {
		return
	}

	context.attribLocationTex = gl.GetUniformLocation(context.shaderHandle, "Texture")
	context.attribLocationProjMtx = gl.GetUniformLocation(context.shaderHandle, "ProjMtx")
	context.attribLocationPosition = gl.GetAttribLocation(context.shaderHandle, "Position")
	context.attribLocationUV = gl.GetAttribLocation(context.shaderHandle, "UV")
	context.attribLocationColor = gl.GetAttribLocation(context.shaderHandle, "Color")

	buffers := gl.GenBuffers(2)
	context.vboHandle = buffers[0]
	context.elementsHandle = buffers[1]

	context.createFontsTexture(gl)

	return
}

func (context *Context) createFontsTexture(gl opengl.OpenGl) {
	io := imgui.CurrentIO()
	image := io.Fonts().TextureDataAlpha8()

	context.fontTexture = gl.GenTextures(1)[0]
	gl.BindTexture(opengl.TEXTURE_2D, context.fontTexture)
	gl.TexParameteri(opengl.TEXTURE_2D, opengl.TEXTURE_MIN_FILTER, opengl.LINEAR)
	gl.TexParameteri(opengl.TEXTURE_2D, opengl.TEXTURE_MAG_FILTER, opengl.LINEAR)
	gl.PixelStorei(opengl.UNPACK_ROW_LENGTH, 0)
	gl.TexImage2D(opengl.TEXTURE_2D, 0, opengl.RED, int32(image.Width), int32(image.Height),
		0, opengl.RED, opengl.UNSIGNED_BYTE, image.Pixels)

	io.Fonts().SetTextureID(imgui.TextureID(context.fontTexture))

	gl.BindTexture(opengl.TEXTURE_2D, 0)
}

func (context *Context) destroyDeviceObjects(gl opengl.OpenGl) {
	if context.vboHandle != 0 {
		gl.DeleteBuffers([]uint32{context.vboHandle})
	}
	context.vboHandle = 0
	if context.elementsHandle != 0 {
		gl.DeleteBuffers([]uint32{context.elementsHandle})
	}
	context.elementsHandle = 0

	if context.shaderHandle != 0 {
		gl.DeleteProgram(context.shaderHandle)
	}
	context.shaderHandle = 0

	if context.fontTexture != 0 {
		gl.DeleteTextures([]uint32{context.fontTexture})
		imgui.CurrentIO().Fonts().SetTextureID(0)
		context.fontTexture = 0
	}
}
