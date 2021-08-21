package render

import (
	"math"

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

out vec3 position;
out float zCenter;

void main(void) {
	gl_Position = projectionMatrix * viewMatrix * modelMatrix * vec4(vertexPosition, 1.0);
	zCenter = (projectionMatrix * viewMatrix * modelMatrix * vec4(0.0, 0.0, 0.0, 1.0)).z;
	position = gl_Position.xyz;
}
`

var orientationViewFragmentShaderSource = `
#version 150
precision mediump float;

uniform vec4 foregroundColor;
uniform vec4 backgroundColor;

in vec3 position;
in float zCenter;

out vec4 fragColor;

void main(void) {
	if (position.z <= zCenter)
	{
    	fragColor = foregroundColor;
	}
  	else
	{
    	fragColor = backgroundColor;
	}
}
`

// OrientationView is a control that displays how an object is oriented.
// This is still not working properly. The orientation arrow is not properly rotated,
// and I suspect the typical issue of three-angle rotation versus quaternion rotation.
// Though I doubt the original engine used quaternions, it remains too confusing for me.
type OrientationView struct {
	context Context

	program                 uint32
	vao                     *opengl.VertexArrayObject
	vertexPositionBuffer    uint32
	vertexPositionAttrib    int32
	modelMatrixUniform      opengl.Matrix4Uniform
	viewMatrixUniform       opengl.Matrix4Uniform
	projectionMatrixUniform opengl.Matrix4Uniform
	foregroundColorUniform  opengl.Vector4Uniform
	backgroundColorUniform  opengl.Vector4Uniform

	baseOrientation    mgl.Mat4
	rotation           mgl.Vec3
	xRingVerticesStart int
	yRingVerticesStart int
	zRingVerticesStart int
	vertices           int
}

// NewOrientationView returns a new instance.
func NewOrientationView(context Context, baseOrientation mgl.Mat4, rotation mgl.Vec3) *OrientationView {
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
		foregroundColorUniform:  opengl.Vector4Uniform(gl.GetUniformLocation(program, "foregroundColor")),
		backgroundColorUniform:  opengl.Vector4Uniform(gl.GetUniformLocation(program, "backgroundColor")),

		baseOrientation: baseOrientation,
		rotation:        rotation,
	}
	{
		gl.BindBuffer(opengl.ARRAY_BUFFER, view.vertexPositionBuffer)
		var vertices []float32
		radius := 0.5
		vertices = append(vertices, 0.0, 0.0, 0.0)
		vertices = append(vertices, 0.4, 0.0, 0.0)
		vertices = append(vertices, 0.0, 0.0, 0.0)
		vertices = append(vertices, 0.0, 0.4, 0.0)
		vertices = append(vertices, 0.0, 0.0, 0.0)
		vertices = append(vertices, 0.0, 0.0, 0.4)

		view.xRingVerticesStart = len(vertices)
		for angle := 0.0; angle < float64(toRad(360.0)); angle += float64(toRad(4.0)) {
			vertices = append(vertices, 0.0, float32(radius*math.Cos(angle)), float32(radius*math.Sin(angle)))
		}
		vertices = append(vertices, 0.0, float32(radius*math.Cos(0)), float32(radius*math.Sin(0)))

		view.yRingVerticesStart = len(vertices)
		for angle := 0.0; angle < float64(toRad(360.0)); angle += float64(toRad(4.0)) {
			vertices = append(vertices, float32(radius*math.Cos(angle)), 0.0, float32(radius*math.Sin(angle)))
		}
		vertices = append(vertices, float32(radius*math.Cos(0)), 0.0, float32(radius*math.Sin(0)))

		view.zRingVerticesStart = len(vertices)
		for angle := 0.0; angle < float64(toRad(360.0)); angle += float64(toRad(4.0)) {
			vertices = append(vertices, float32(radius*math.Cos(angle)), float32(radius*math.Sin(angle)), 0.0)
		}
		vertices = append(vertices, float32(radius*math.Cos(0)), float32(radius*math.Sin(0)), 0.0)

		view.vertices = len(vertices)
		gl.BufferData(opengl.ARRAY_BUFFER, len(vertices)*4, vertices, opengl.STATIC_DRAW)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	}
	view.vao.WithSetter(func(gl opengl.OpenGL) {
		gl.EnableVertexAttribArray(uint32(view.vertexPositionAttrib))
		gl.BindBuffer(opengl.ARRAY_BUFFER, view.vertexPositionBuffer)
		gl.VertexAttribOffset(uint32(view.vertexPositionAttrib), 3, opengl.FLOAT, false, 0, 0)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	})

	return view
}

// Dispose cleans up resources of the view.
func (view *OrientationView) Dispose() {
	gl := view.context.OpenGL
	view.vao.Dispose()
	gl.DeleteBuffers([]uint32{view.vertexPositionBuffer})
	gl.DeleteProgram(view.program)
}

func toRad(degree float32) float32 {
	return (degree * math.Pi * 2.0) / 360.0
}

// Render renders the orientation view for given orientation.
func (view *OrientationView) Render(orientation mgl.Vec3) {
	gl := view.context.OpenGL

	view.vao.OnShader(func() {
		view.projectionMatrixUniform.Set(gl, &view.context.ProjectionMatrix)
		view.viewMatrixUniform.Set(gl, view.context.ViewMatrix)

		view.renderRings(false)
		view.renderArrow(
			mgl.Ident4().
				Mul4(mgl.HomogRotate3D(toRad(orientation[0]), mgl.Vec3{view.rotation[0], 0.0, 0.0})).
				Mul4(mgl.HomogRotate3D(toRad(orientation[1]), mgl.Vec3{0.0, view.rotation[1], 0.0})).
				Mul4(mgl.HomogRotate3D(toRad(orientation[2]), mgl.Vec3{0.0, 0.0, view.rotation[2]})))
		view.renderRings(true)
	})
}

func (view *OrientationView) renderRings(front bool) {
	foreground := func(color [4]float32) [4]float32 {
		if !front {
			return [4]float32{0.0, 0.0, 0.0, 0.0}
		}
		return color
	}
	background := func(color [4]float32) [4]float32 {
		if front {
			return [4]float32{0.0, 0.0, 0.0, 0.0}
		}
		return color
	}

	// draw Z-rotation ring
	view.renderRing(view.zRingVerticesStart, view.vertices,
		foreground([4]float32{0.0, 0.0, 1.0, 1.0}), background([4]float32{0.0, 0.0, 0.8, 0.7}))

	// draw Y-rotation ring
	view.renderRing(view.yRingVerticesStart, view.zRingVerticesStart,
		foreground([4]float32{0.0, 1.0, 0.0, 1.0}), background([4]float32{0.0, 0.8, 0.0, 0.7}))

	// draw X-rotation ring
	view.renderRing(view.xRingVerticesStart, view.yRingVerticesStart,
		foreground([4]float32{1.0, 0.0, 0.0, 1.0}), background([4]float32{0.8, 0.0, 0.0, 0.7}))
}

func (view *OrientationView) renderRing(start, end int, foregroundColor, backgroundColor [4]float32) {
	gl := view.context.OpenGL

	modelMatrix := mgl.Ident4().Mul4(view.baseOrientation)
	view.foregroundColorUniform.Set(gl, &foregroundColor)
	view.backgroundColorUniform.Set(gl, &backgroundColor)
	view.modelMatrixUniform.Set(gl, &modelMatrix)
	gl.DrawArrays(opengl.LINES, int32(start)/3, int32(end-start)/3)
}

func (view *OrientationView) renderArrow(rotation mgl.Mat4) {
	gl := view.context.OpenGL

	modelMatrix := mgl.Ident4().Mul4(view.baseOrientation).Mul4(rotation)
	view.foregroundColorUniform.Set(gl, &[4]float32{0.0, 0.0, 1.0, 1.0})
	view.backgroundColorUniform.Set(gl, &[4]float32{0.0, 0.0, 0.6, 1.0})
	view.modelMatrixUniform.Set(gl, &modelMatrix)
	gl.DrawArrays(opengl.LINES, 4, 2)

	modelMatrix = mgl.Ident4().Mul4(view.baseOrientation).Mul4(rotation)
	view.foregroundColorUniform.Set(gl, &[4]float32{0.0, 1.0, 0.0, 1.0})
	view.backgroundColorUniform.Set(gl, &[4]float32{0.0, 0.6, 0.0, 1.0})
	view.modelMatrixUniform.Set(gl, &modelMatrix)
	gl.DrawArrays(opengl.LINES, 2, 2)

	modelMatrix = mgl.Ident4().Mul4(view.baseOrientation).Mul4(rotation)
	view.foregroundColorUniform.Set(gl, &[4]float32{1.0, 0.0, 0.0, 1.0})
	view.backgroundColorUniform.Set(gl, &[4]float32{0.6, 0.0, 0.0, 1.0})
	view.modelMatrixUniform.Set(gl, &modelMatrix)
	gl.DrawArrays(opengl.LINES, 0, 2)
}
