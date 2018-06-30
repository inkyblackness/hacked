package levels

import (
	"fmt"
	"github.com/inkyblackness/hacked/editor/render"
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ui/opengl"
)

var mapTileGridVertexShaderSource = `
#version 150
precision mediump float;

in vec3 vertexPosition;

uniform mat4 viewMatrix;
uniform mat4 projectionMatrix;

out float height;

void main(void) {
	gl_Position = projectionMatrix * viewMatrix * vec4(vertexPosition.xy, 0.0, 1.0);
	height = vertexPosition.z;
}
`

var mapTileGridFragmentShaderSource = `
#version 150
precision mediump float;

in float height;
out vec4 fragColor;

void main(void) {
	fragColor = vec4(0.0, 0.8, 0.0, height);
}
`

// WallMapper returns basic information to draw a 2D map.
type WallMapper interface {
	MapGridInfo(x, y int) (level.TileType, level.WallHeights)
}

// MapGrid renders the grid of the map, based on calculated wall heights.
type MapGrid struct {
	context *render.Context

	program                 uint32
	vao                     *opengl.VertexArrayObject
	vertexPositionBuffer    uint32
	vertexPositionAttrib    int32
	viewMatrixUniform       opengl.Matrix4Uniform
	projectionMatrixUniform opengl.Matrix4Uniform
}

// NewMapGrid returns a new instance.
func NewMapGrid(context *render.Context) *MapGrid {
	gl := context.OpenGL
	program, programErr := opengl.LinkNewStandardProgram(gl, mapTileGridVertexShaderSource, mapTileGridFragmentShaderSource)

	if programErr != nil {
		panic(fmt.Errorf("MapGrid shader failed: %v", programErr))
	}
	grid := &MapGrid{
		context:                 context,
		program:                 program,
		vao:                     opengl.NewVertexArrayObject(gl, program),
		vertexPositionBuffer:    gl.GenBuffers(1)[0],
		vertexPositionAttrib:    gl.GetAttribLocation(program, "vertexPosition"),
		viewMatrixUniform:       opengl.Matrix4Uniform(gl.GetUniformLocation(program, "viewMatrix")),
		projectionMatrixUniform: opengl.Matrix4Uniform(gl.GetUniformLocation(program, "projectionMatrix")),
	}

	grid.vao.WithSetter(func(gl opengl.OpenGL) {
		gl.EnableVertexAttribArray(uint32(grid.vertexPositionAttrib))
		gl.BindBuffer(opengl.ARRAY_BUFFER, grid.vertexPositionBuffer)
		gl.VertexAttribOffset(uint32(grid.vertexPositionAttrib), 3, opengl.FLOAT, false, 0, 0)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	})

	return grid
}

// Dispose releases any internal resources
func (grid *MapGrid) Dispose() {
	gl := grid.context.OpenGL
	gl.DeleteProgram(grid.program)
	gl.DeleteBuffers([]uint32{grid.vertexPositionBuffer})
	grid.vao.Dispose()
}

// Render renders
func (grid *MapGrid) Render(mapper WallMapper) {
	gl := grid.context.OpenGL

	grid.vao.OnShader(func() {
		grid.viewMatrixUniform.Set(gl, grid.context.ViewMatrix)
		grid.projectionMatrixUniform.Set(gl, &grid.context.ProjectionMatrix)

		gl.BindBuffer(opengl.ARRAY_BUFFER, grid.vertexPositionBuffer)

		vertices := make([]float32, 0, ((4*3)+2)*2*3)
		heightFactor := [32*2 + 1]float32{
			0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
			0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
			0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
			0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
			0.0,
			0.1, 0.2, 0.3, 0.3, 0.4, 0.4, 0.5, 0.7,
			1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0,
			1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0,
			1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0,
		}
		for y := 0; y < 64; y++ {
			for x := 0; x < 64; x++ {
				vertices = vertices[0:0]
				tileType, wallHeights := mapper.MapGridInfo(x, y)

				finePerFraction := fineCoordinatesPerTileSide / 3
				left := float32(x) * fineCoordinatesPerTileSide
				right := left + fineCoordinatesPerTileSide
				bottom := float32(y) * fineCoordinatesPerTileSide
				top := bottom + fineCoordinatesPerTileSide

				for i, height := range wallHeights.North {
					vertices = append(vertices,
						left+finePerFraction*float32(i), top, heightFactor[int(height)+32],
						left+finePerFraction*float32(i+1), top, heightFactor[int(height)+32])
				}
				for i, height := range wallHeights.South {
					vertices = append(vertices,
						right-finePerFraction*float32(i), bottom, heightFactor[int(height)+32],
						right-finePerFraction*float32(i+1), bottom, heightFactor[int(height)+32])
				}
				for i, height := range wallHeights.East {
					vertices = append(vertices,
						right, top-finePerFraction*float32(i), heightFactor[int(height)+32],
						right, top-finePerFraction*float32(i+1), heightFactor[int(height)+32])
				}
				for i, height := range wallHeights.West {
					vertices = append(vertices,
						left, bottom+finePerFraction*float32(i), heightFactor[int(height)+32],
						left, bottom+finePerFraction*float32(i+1), heightFactor[int(height)+32])
				}
				if tileType == level.TileTypeDiagonalOpenNorthEast || tileType == level.TileTypeDiagonalOpenSouthWest {
					vertices = append(vertices, left, top, 1.0, right, bottom, 1.0)
				}
				if tileType == level.TileTypeDiagonalOpenNorthWest || tileType == level.TileTypeDiagonalOpenSouthEast {
					vertices = append(vertices, left, bottom, 1.0, right, top, 1.0)
				}

				if len(vertices) > 0 {
					gl.BufferData(opengl.ARRAY_BUFFER, len(vertices)*4, vertices, opengl.STATIC_DRAW)
					gl.DrawArrays(opengl.LINES, 0, int32(len(vertices)/3))
				}
			}
		}
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	})
}
