package levels

import (
	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/hacked/editor/render"
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ui/opengl"
)

var mapColorsVertexShaderSource = `
#version 150
precision mediump float;

in vec4 vertexColor;
in vec2 vertexPosition;

uniform mat4 modelMatrix;
uniform mat4 viewMatrix;
uniform mat4 projectionMatrix;

out vec4 color;

void main(void) {
   gl_Position = projectionMatrix * viewMatrix * modelMatrix * vec4(vertexPosition.xy, 0.0, 1.0);
   color = vertexColor;
}
`

var mapColorsFragmentShaderSource = `
#version 150
precision mediump float;

in vec4 color;

out vec4 fragColor;

void main(void) {
   fragColor = color;
}
`

// ColorQuery returns the color for the specified tile.
type ColorQuery func(level.TilePosition) [4]float32

// MapColors is a renderable for the tile colorings.
type MapColors struct {
	context *render.Context

	program                 uint32
	vao                     *opengl.VertexArrayObject
	vertexPositionBuffer    uint32
	vertexPositionAttrib    int32
	vertexColorBuffer       uint32
	vertexColorAttrib       int32
	modelMatrixUniform      opengl.Matrix4Uniform
	viewMatrixUniform       opengl.Matrix4Uniform
	projectionMatrixUniform opengl.Matrix4Uniform

	lastColorGrid        [][][4]float32
	lastColorGridColumns int
	lastColorGridRows    int
}

// NewMapColors returns a new instance of a renderable for tile colorings.
func NewMapColors(context *render.Context) *MapColors {
	gl := context.OpenGL
	program, programErr := opengl.LinkNewStandardProgram(gl, mapColorsVertexShaderSource, mapColorsFragmentShaderSource)

	if programErr != nil {
		panic(opengl.NamedShaderError{Name: "MapColorsShader", Nested: programErr})
	}
	renderable := &MapColors{
		context:                 context,
		program:                 program,
		vao:                     opengl.NewVertexArrayObject(gl, program),
		vertexPositionBuffer:    gl.GenBuffers(1)[0],
		vertexPositionAttrib:    gl.GetAttribLocation(program, "vertexPosition"),
		vertexColorBuffer:       gl.GenBuffers(1)[0],
		vertexColorAttrib:       gl.GetAttribLocation(program, "vertexColor"),
		modelMatrixUniform:      opengl.Matrix4Uniform(gl.GetUniformLocation(program, "modelMatrix")),
		viewMatrixUniform:       opengl.Matrix4Uniform(gl.GetUniformLocation(program, "viewMatrix")),
		projectionMatrixUniform: opengl.Matrix4Uniform(gl.GetUniformLocation(program, "projectionMatrix")),
	}

	{
		top := float32(fineCoordinatesPerTileSide)
		left := float32(0.0)
		right := float32(fineCoordinatesPerTileSide)
		bottom := float32(0.0)

		vertices := []float32{
			left, bottom,
			left, top,
			right, top,

			right, top,
			right, bottom,
			left, bottom}

		gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexPositionBuffer)
		gl.BufferData(opengl.ARRAY_BUFFER, len(vertices)*4, vertices, opengl.STATIC_DRAW)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	}

	renderable.vao.WithSetter(func(gl opengl.OpenGL) {
		gl.EnableVertexAttribArray(uint32(renderable.vertexPositionAttrib))
		gl.EnableVertexAttribArray(uint32(renderable.vertexColorAttrib))
		gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexPositionBuffer)
		gl.VertexAttribOffset(uint32(renderable.vertexPositionAttrib), 2, opengl.FLOAT, false, 0, 0)
		gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexColorBuffer)
		gl.VertexAttribOffset(uint32(renderable.vertexColorAttrib), 4, opengl.FLOAT, false, 0, 0)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	})

	return renderable
}

// Dispose releases any internal resources.
func (renderable *MapColors) Dispose() {
	gl := renderable.context.OpenGL
	gl.DeleteProgram(renderable.program)
	gl.DeleteBuffers([]uint32{renderable.vertexPositionBuffer})
	renderable.vao.Dispose()
}

// Render renders the renderable.
func (renderable *MapColors) Render(columnCount, rowCount int, query ColorQuery) {
	gl := renderable.context.OpenGL

	if (renderable.lastColorGridColumns != columnCount) || (renderable.lastColorGridRows != rowCount) {
		renderable.lastColorGrid = make([][][4]float32, rowCount+1)
		for y := 0; y < rowCount; y++ {
			renderable.lastColorGrid[y] = make([][4]float32, columnCount+1)
		}
		renderable.lastColorGrid[rowCount] = make([][4]float32, columnCount+1)
		renderable.lastColorGridColumns = columnCount
		renderable.lastColorGridRows = rowCount
	}
	for y, row := range renderable.lastColorGrid {
		for x := 0; x < columnCount; x++ {
			row[x] = query(level.TilePosition{X: byte(x), Y: byte(y)})
		}
	}

	renderable.vao.OnShader(func() {
		var colors [24]float32

		renderable.viewMatrixUniform.Set(gl, renderable.context.ViewMatrix)
		renderable.projectionMatrixUniform.Set(gl, &renderable.context.ProjectionMatrix)

		gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexColorBuffer)

		for y := 0; y < rowCount; y++ {
			for x := 0; x < columnCount; x++ {
				colorBottomLeft := renderable.lastColorGrid[y][x]
				colorBottomRight := renderable.lastColorGrid[y][x+1]
				colorTopLeft := renderable.lastColorGrid[y+1][x]
				colorTopRight := renderable.lastColorGrid[y+1][x+1]

				modelMatrix := mgl.Ident4().
					Mul4(mgl.Translate3D((float32(x))*fineCoordinatesPerTileSide, (float32(y))*fineCoordinatesPerTileSide, 0.0))
				renderable.modelMatrixUniform.Set(gl, &modelMatrix)

				copy(colors[0:4], colorBottomLeft[:])
				copy(colors[4:8], colorTopLeft[:])
				copy(colors[8:12], colorTopRight[:])
				copy(colors[12:16], colorTopRight[:])
				copy(colors[16:20], colorBottomRight[:])
				copy(colors[20:24], colorBottomLeft[:])

				gl.BufferData(opengl.ARRAY_BUFFER, len(colors)*4, colors[:], opengl.STATIC_DRAW)
				gl.DrawArrays(opengl.TRIANGLES, 0, 6)
			}
		}

		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	})
}
