package levels

import (
	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/hacked/editor/render"
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ui/opengl"
)

var mapTileGridVertexShaderSource = `
#version 150
precision mediump float;

in vec3 vertexPosition;
out float colorAlpha;

uniform mat4 modelMatrix;
uniform mat4 viewMatrix;
uniform mat4 projectionMatrix;

void main(void) {
	gl_Position = projectionMatrix * viewMatrix * modelMatrix * vec4(vertexPosition.xy, 0.0, 1.0);
	colorAlpha = vertexPosition.z;
}
`

var mapTileGridFragmentShaderSource = `
#version 150
precision mediump float;

uniform vec4 color;
in float colorAlpha;
out vec4 fragColor;

void main(void) {
	fragColor = vec4(color.rgb, color.a*colorAlpha);
}
`

// TileMapper returns basic information to draw a 2D map.
type TileMapper interface {
	MapGridInfo(x, y int) (level.TileType, level.TileSlopeControl, level.WallHeights)
}

// MapGrid renders the grid of the map, based on calculated wall heights.
type MapGrid struct {
	context *render.Context

	program                 uint32
	vao                     *opengl.VertexArrayObject
	vertexPositionBuffer    uint32
	vertexPositionAttrib    int32
	colorUniform            opengl.Vector4Uniform
	modelMatrixUniform      opengl.Matrix4Uniform
	viewMatrixUniform       opengl.Matrix4Uniform
	projectionMatrixUniform opengl.Matrix4Uniform

	tickVertices []float32
}

// NewMapGrid returns a new instance.
func NewMapGrid(context *render.Context) *MapGrid {
	gl := context.OpenGL
	program, programErr := opengl.LinkNewStandardProgram(gl, mapTileGridVertexShaderSource, mapTileGridFragmentShaderSource)
	if programErr != nil {
		panic(opengl.NamedShaderError{Name: "MapGrid", Nested: programErr})
	}
	grid := &MapGrid{
		context:                 context,
		program:                 program,
		vao:                     opengl.NewVertexArrayObject(gl, program),
		vertexPositionBuffer:    gl.GenBuffers(1)[0],
		vertexPositionAttrib:    gl.GetAttribLocation(program, "vertexPosition"),
		colorUniform:            opengl.Vector4Uniform(gl.GetUniformLocation(program, "color")),
		modelMatrixUniform:      opengl.Matrix4Uniform(gl.GetUniformLocation(program, "modelMatrix")),
		viewMatrixUniform:       opengl.Matrix4Uniform(gl.GetUniformLocation(program, "viewMatrix")),
		projectionMatrixUniform: opengl.Matrix4Uniform(gl.GetUniformLocation(program, "projectionMatrix")),
	}
	grid.tickVertices = grid.calculateTickVertices()

	grid.vao.WithSetter(func(gl opengl.OpenGL) {
		gl.EnableVertexAttribArray(uint32(grid.vertexPositionAttrib))
		gl.BindBuffer(opengl.ARRAY_BUFFER, grid.vertexPositionBuffer)
		gl.VertexAttribOffset(uint32(grid.vertexPositionAttrib), 3, opengl.FLOAT, false, 0, 0)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	})

	return grid
}

// Dispose releases any internal resources.
func (grid *MapGrid) Dispose() {
	gl := grid.context.OpenGL
	gl.DeleteProgram(grid.program)
	gl.DeleteBuffers([]uint32{grid.vertexPositionBuffer})
	grid.vao.Dispose()
}

func (grid *MapGrid) calculateTickVertices() []float32 {
	dotHalf := float32(0.05)
	dotBase := float32(0.5) - (dotHalf * 2.0)
	alpha := float32(0.9)

	top := dotBase + dotHalf
	topEnd := dotBase - dotHalf
	left := -dotBase - dotHalf
	leftEnd := -dotBase + dotHalf

	right := top
	rightEnd := topEnd
	bottom := left
	bottomEnd := leftEnd

	return []float32{
		left, top, alpha, leftEnd, top, alpha, left, topEnd, alpha,
		leftEnd, top, alpha, leftEnd, topEnd, alpha, left, topEnd, alpha,

		right, top, alpha, right, topEnd, alpha, rightEnd, top, alpha,
		rightEnd, top, alpha, right, topEnd, alpha, rightEnd, topEnd, alpha,

		right, bottom, alpha, rightEnd, bottom, alpha, right, bottomEnd, alpha,
		right, bottomEnd, alpha, rightEnd, bottom, alpha, rightEnd, bottomEnd, alpha,

		left, bottomEnd, alpha, leftEnd, bottom, alpha, left, bottom, alpha,
		left, bottomEnd, alpha, leftEnd, bottomEnd, alpha, leftEnd, bottom, alpha,
	}
}

// Render renders the grid.
func (grid *MapGrid) Render(columns, rows int, mapper TileMapper) {
	gl := grid.context.OpenGL

	var slopeTicksByType = [][]int{
		nil,
		nil,

		nil,
		nil,
		nil,
		nil,

		{0, 1},
		{1, 2},
		{2, 3},
		{3, 0},

		{3, 0, 1},
		{0, 1, 2},
		{1, 2, 3},
		{2, 3, 0},

		{2},
		{3},
		{0},
		{1},
	}

	floorTickStarts := []int32{0, 6, 12, 18}
	ceilingTickStarts := []int32{3, 9, 15, 21}

	grid.vao.OnShader(func() {
		modelMatrix := mgl.Ident4()
		lineColor := [4]float32{0.0, 0.8, 0.0, 1.0}
		floorTickColor := [4]float32{0.0, 0.8, 0.0, 0.9}
		ceilingTickColor := [4]float32{0.8, 0.0, 0.0, 0.9}

		grid.viewMatrixUniform.Set(gl, grid.context.ViewMatrix)
		grid.projectionMatrixUniform.Set(gl, &grid.context.ProjectionMatrix)
		gl.BindBuffer(opengl.ARRAY_BUFFER, grid.vertexPositionBuffer)

		// This table does not consider height shift of the level.
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
		var vertexBuffer [((4 * 3) + 2) * 2 * 3]float32
		for y := 0; y < rows; y++ {
			for x := 0; x < columns; x++ {
				modelMatrix = mgl.Ident4().
					Mul4(mgl.Translate3D((float32(x)+0.5)*fineCoordinatesPerTileSide, (float32(y)+0.5)*fineCoordinatesPerTileSide, 0.0)).
					Mul4(mgl.Scale3D(fineCoordinatesPerTileSide, fineCoordinatesPerTileSide, 1.0))
				grid.modelMatrixUniform.Set(gl, &modelMatrix)

				vertices := vertexBuffer[0:0]
				tileType, slopeControl, wallHeights := mapper.MapGridInfo(x, y)

				left := float32(-0.5)
				right := float32(0.5)
				top := float32(0.5)
				bottom := float32(-0.5)
				finePerFraction := float32(1.0 / 3.0)

				for i, height := range wallHeights.North {
					vertices = append(vertices,
						left+finePerFraction*float32(i), top, heightFactor[int(height)+32],
						left+finePerFraction*float32(i+1), top, heightFactor[int(height)+32])
				}
				for i, height := range wallHeights.East {
					vertices = append(vertices,
						right, top-finePerFraction*float32(i), heightFactor[int(height)+32],
						right, top-finePerFraction*float32(i+1), heightFactor[int(height)+32])
				}
				for i, height := range wallHeights.South {
					vertices = append(vertices,
						right-finePerFraction*float32(i), bottom, heightFactor[int(height)+32],
						right-finePerFraction*float32(i+1), bottom, heightFactor[int(height)+32])
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

				grid.colorUniform.Set(gl, &lineColor)
				gl.BufferData(opengl.ARRAY_BUFFER, len(vertices)*4, vertices, opengl.STATIC_DRAW)
				gl.DrawArrays(opengl.LINES, 0, int32(len(vertices)/3))

				var floorTicks []int
				var ceilingTicks []int

				if slopeControl != level.TileSlopeControlFloorFlat {
					floorTicks = slopeTicksByType[tileType]
				}
				if slopeControl == level.TileSlopeControlCeilingMirrored {
					ceilingTicks = slopeTicksByType[tileType]
				} else if slopeControl != level.TileSlopeControlCeilingFlat {
					ceilingTicks = slopeTicksByType[tileType.Info().SlopeInvertedType]
				}

				gl.BufferData(opengl.ARRAY_BUFFER, len(grid.tickVertices)*4, grid.tickVertices, opengl.STATIC_DRAW)
				grid.colorUniform.Set(gl, &floorTickColor)
				for _, index := range floorTicks {
					gl.DrawArrays(opengl.TRIANGLES, floorTickStarts[index], 3)
				}
				grid.colorUniform.Set(gl, &ceilingTickColor)
				for _, index := range ceilingTicks {
					gl.DrawArrays(opengl.TRIANGLES, ceilingTickStarts[index], 3)
				}
			}
		}
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	})
}
