package levels

import (
	"fmt"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/hacked/editor/render"
	"github.com/inkyblackness/hacked/ui/opengl"
)

var HighlighterVertexShaderSource = `
#version 150
precision mediump float;

in vec3 vertexPosition;

uniform mat4 modelMatrix;
uniform mat4 viewMatrix;
uniform mat4 projectionMatrix;

void main(void) {
	gl_Position = projectionMatrix * viewMatrix * modelMatrix * vec4(vertexPosition, 1.0);
}
`

var HighlighterFragmentShaderSource = `
#version 150
precision mediump float;

uniform vec4 inColor;
out vec4 fragColor;

void main(void) {
	fragColor = inColor;
}
`

// Highlighter draws a simple highlighting of a rectangular area.
type Highlighter struct {
	context *render.Context

	program                 uint32
	vao                     *opengl.VertexArrayObject
	vertexPositionBuffer    uint32
	vertexPositionAttrib    int32
	modelMatrixUniform      opengl.Matrix4Uniform
	viewMatrixUniform       opengl.Matrix4Uniform
	projectionMatrixUniform opengl.Matrix4Uniform
	inColorUniform          opengl.Vector4Uniform
}

// NewHighlighter returns a new instance of Highlighter.
func NewHighlighter(context *render.Context) *Highlighter {
	gl := context.OpenGL
	program, programErr := opengl.LinkNewStandardProgram(gl, HighlighterVertexShaderSource, HighlighterFragmentShaderSource)

	if programErr != nil {
		panic(fmt.Errorf("highlighter shader failed: %v", programErr))
	}
	highlighter := &Highlighter{
		context: context,
		program: program,

		vao:                     opengl.NewVertexArrayObject(gl, program),
		vertexPositionBuffer:    gl.GenBuffers(1)[0],
		vertexPositionAttrib:    gl.GetAttribLocation(program, "vertexPosition"),
		modelMatrixUniform:      opengl.Matrix4Uniform(gl.GetUniformLocation(program, "modelMatrix")),
		viewMatrixUniform:       opengl.Matrix4Uniform(gl.GetUniformLocation(program, "viewMatrix")),
		projectionMatrixUniform: opengl.Matrix4Uniform(gl.GetUniformLocation(program, "projectionMatrix")),
		inColorUniform:          opengl.Vector4Uniform(gl.GetUniformLocation(program, "inColor"))}

	{
		gl.BindBuffer(opengl.ARRAY_BUFFER, highlighter.vertexPositionBuffer)
		half := float32(0.5)
		var vertices = []float32{
			-half, -half, 0.0,
			half, -half, 0.0,
			half, half, 0.0,

			half, half, 0.0,
			-half, half, 0.0,
			-half, -half, 0.0}
		gl.BufferData(opengl.ARRAY_BUFFER, len(vertices)*4, vertices, opengl.STATIC_DRAW)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	}

	highlighter.vao.OnShader(func() {
		gl.EnableVertexAttribArray(uint32(highlighter.vertexPositionAttrib))
		gl.BindBuffer(opengl.ARRAY_BUFFER, highlighter.vertexPositionBuffer)
		gl.VertexAttribOffset(uint32(highlighter.vertexPositionAttrib), 3, opengl.FLOAT, false, 0, 0)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	})

	return highlighter
}

// Dispose releases all resources.
func (highlighter *Highlighter) Dispose() {
	gl := highlighter.context.OpenGL

	highlighter.vao.Dispose()
	gl.DeleteBuffers([]uint32{highlighter.vertexPositionBuffer})
	gl.DeleteShader(highlighter.program)
}

// Render renders the highlights.
func (highlighter *Highlighter) Render(positions []MapPosition, sideLength float32, color [4]float32) {
	gl := highlighter.context.OpenGL

	highlighter.vao.OnShader(func() {
		highlighter.viewMatrixUniform.Set(gl, highlighter.context.ViewMatrix)
		highlighter.projectionMatrixUniform.Set(gl, &highlighter.context.ProjectionMatrix)
		highlighter.inColorUniform.Set(gl, &color)

		for _, pos := range positions {
			modelMatrix := mgl.Ident4().
				Mul4(mgl.Translate3D(float32(pos.X), float32(pos.Y), 0.0)).
				Mul4(mgl.Scale3D(sideLength, sideLength, 1.0))

			highlighter.modelMatrixUniform.Set(gl, &modelMatrix)

			gl.DrawArrays(opengl.TRIANGLES, 0, 6)
		}
	})
}
