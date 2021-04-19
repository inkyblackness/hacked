package render

import (
	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/hacked/ui/opengl"
)

var orientationViewVertexShaderSource = `
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

var orientationViewFragmentShaderSource = `
#version 150
precision mediump float;

out vec4 fragColor;

void main(void) {
	fragColor = vec4(1.0, 0.0, 0.0, 1.0);
}
`

type OrientationView struct {
	context Context

	program                 uint32
	vao                     *opengl.VertexArrayObject
	vertexPositionBuffer    uint32
	vertexPositionAttrib    int32
	modelMatrixUniform      opengl.Matrix4Uniform
	viewMatrixUniform       opengl.Matrix4Uniform
	projectionMatrixUniform opengl.Matrix4Uniform
}

func NewOrientationView(context Context, zero mgl.Vec3, rotation mgl.Vec3) *OrientationView {
	gl := context.OpenGL
	program, programErr := opengl.LinkNewStandardProgram(gl, orientationViewVertexShaderSource, orientationViewFragmentShaderSource)

	if programErr != nil {
		panic(opengl.NamedShaderError{Name: "OrientationViewShader", Nested: programErr})
	}

	view := &OrientationView{
		context: context,
		program: program,

		vao:                     opengl.NewVertexArrayObject(gl, program),
		vertexPositionBuffer:    gl.GenBuffers(1)[0],
		vertexPositionAttrib:    gl.GetAttribLocation(program, "vertexPosition"),
		modelMatrixUniform:      opengl.Matrix4Uniform(gl.GetUniformLocation(program, "modelMatrix")),
		viewMatrixUniform:       opengl.Matrix4Uniform(gl.GetUniformLocation(program, "viewMatrix")),
		projectionMatrixUniform: opengl.Matrix4Uniform(gl.GetUniformLocation(program, "projectionMatrix")),
	}
	return view
}

func (view *OrientationView) Dispose() {
	gl := view.context.OpenGL
	view.vao.Dispose()
	gl.DeleteBuffers([]uint32{view.vertexPositionBuffer})
	gl.DeleteProgram(view.program)
}

func (view *OrientationView) Render(orientation mgl.Vec3) {

}
