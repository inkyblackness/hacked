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

in vec4 vertexPosition;
out float colorAlpha;

uniform mat4 alphaMatrix;

uniform mat4 modelMatrix;
uniform mat4 viewMatrix;
uniform mat4 projectionMatrix;

void main(void) {
	gl_Position = projectionMatrix * viewMatrix * modelMatrix * vec4(vertexPosition.xy, 0.0, 1.0);
	colorAlpha = alphaMatrix[int(vertexPosition.z)][int(vertexPosition.w)];
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
	MapGridInfo(pos level.TilePosition) (level.TileType, level.TileSlopeControl, level.WallHeights)
}

// MapGrid renders the grid of the map, based on calculated wall heights.
type MapGrid struct {
	context *render.Context

	program                  uint32
	vao                      *opengl.VertexArrayObject
	vertexPositionBuffer     uint32
	tickVertexPositionBuffer uint32
	vertexPositionAttrib     int32
	colorUniform             opengl.Vector4Uniform

	/*
		The alpha matrix is set up to specify the alpha value used for the wall sections.
		There are four cardinal walls with three sections each, plus two diagonal walls (single section), as well
		as the slope-ticks that re-use the shader.
		The matrix is meant to be used like this:
			N N N D1
			E E E D2
			S S S -
			W W W T
		D1 and D2 being the two diagonals, and T specifying the cell for the tick-marks.
		The vertex buffers are an array of v4 entries, with z/w specifying the row and column index to take.
	*/
	alphaMatrixUniform      opengl.Matrix4Uniform
	modelMatrixUniform      opengl.Matrix4Uniform
	viewMatrixUniform       opengl.Matrix4Uniform
	projectionMatrixUniform opengl.Matrix4Uniform

	slopeTicksByType  [][]int
	floorTickStarts   []int32
	ceilingTickStarts []int32
}

// NewMapGrid returns a new instance.
func NewMapGrid(context *render.Context) *MapGrid {
	gl := context.OpenGL
	program, programErr := opengl.LinkNewStandardProgram(gl, mapTileGridVertexShaderSource, mapTileGridFragmentShaderSource)
	if programErr != nil {
		panic(opengl.NamedShaderError{Name: "MapGrid", Nested: programErr})
	}
	grid := &MapGrid{
		context:                  context,
		program:                  program,
		vao:                      opengl.NewVertexArrayObject(gl, program),
		vertexPositionBuffer:     gl.GenBuffers(1)[0],
		tickVertexPositionBuffer: gl.GenBuffers(1)[0],
		vertexPositionAttrib:     gl.GetAttribLocation(program, "vertexPosition"),
		colorUniform:             opengl.Vector4Uniform(gl.GetUniformLocation(program, "color")),
		alphaMatrixUniform:       opengl.Matrix4Uniform(gl.GetUniformLocation(program, "alphaMatrix")),
		modelMatrixUniform:       opengl.Matrix4Uniform(gl.GetUniformLocation(program, "modelMatrix")),
		viewMatrixUniform:        opengl.Matrix4Uniform(gl.GetUniformLocation(program, "viewMatrix")),
		projectionMatrixUniform:  opengl.Matrix4Uniform(gl.GetUniformLocation(program, "projectionMatrix")),
	}

	{
		tickVertices := grid.calculateTickVertices()
		gl.BindBuffer(opengl.ARRAY_BUFFER, grid.tickVertexPositionBuffer)
		gl.BufferData(opengl.ARRAY_BUFFER, len(tickVertices)*4, tickVertices, opengl.STATIC_DRAW)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	}
	{
		const parts = 3
		left := float32(-0.5)
		right := float32(0.5)
		top := float32(0.5)
		bottom := float32(-0.5)
		finePerPart := float32(1.0 / parts)

		var wallVertices []float32
		for i := 0; i < parts; i++ {
			wallVertices = append(wallVertices, left+(finePerPart*float32(i)), top, 0.0, float32(i))
			wallVertices = append(wallVertices, left+(finePerPart*float32(i+1)), top, 0.0, float32(i))
		}
		for i := 0; i < parts; i++ {
			wallVertices = append(wallVertices, right, top-(finePerPart*float32(i)), 1.0, float32(i))
			wallVertices = append(wallVertices, right, top-(finePerPart*float32(i+1)), 1.0, float32(i))
		}
		for i := 0; i < parts; i++ {
			wallVertices = append(wallVertices, right-(finePerPart*float32(i)), bottom, 2.0, float32(i))
			wallVertices = append(wallVertices, right-(finePerPart*float32(i+1)), bottom, 2.0, float32(i))
		}
		for i := 0; i < parts; i++ {
			wallVertices = append(wallVertices, left, bottom+(finePerPart*float32(i)), 3.0, float32(i))
			wallVertices = append(wallVertices, left, bottom+(finePerPart*float32(i+1)), 3.0, float32(i))
		}
		wallVertices = append(wallVertices, left, top, 0.0, 3.0)
		wallVertices = append(wallVertices, right, bottom, 0.0, 3.0)
		wallVertices = append(wallVertices, right, top, 1.0, 3.0)
		wallVertices = append(wallVertices, left, bottom, 1.0, 3.0)

		gl.BindBuffer(opengl.ARRAY_BUFFER, grid.vertexPositionBuffer)
		gl.BufferData(opengl.ARRAY_BUFFER, len(wallVertices)*4, wallVertices, opengl.STATIC_DRAW)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	}
	{
		grid.slopeTicksByType = [][]int{
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
		grid.floorTickStarts = []int32{0, 6, 12, 18}
		grid.ceilingTickStarts = []int32{3, 9, 15, 21}

	}

	return grid
}

// Dispose releases any internal resources.
func (grid *MapGrid) Dispose() {
	gl := grid.context.OpenGL
	gl.DeleteProgram(grid.program)
	gl.DeleteBuffers([]uint32{grid.vertexPositionBuffer, grid.tickVertexPositionBuffer})
	grid.vao.Dispose()
}

func (grid *MapGrid) calculateTickVertices() []float32 {
	dotHalf := float32(0.05)
	dotBase := float32(0.5) - (dotHalf * 2.0)

	top := dotBase + dotHalf
	topEnd := dotBase - dotHalf
	left := -dotBase - dotHalf
	leftEnd := -dotBase + dotHalf

	right := top
	rightEnd := topEnd
	bottom := left
	bottomEnd := leftEnd

	return []float32{
		left, top, 3.0, 3.0, leftEnd, top, 3.0, 3.0, left, topEnd, 3.0, 3.0,
		leftEnd, top, 3.0, 3.0, leftEnd, topEnd, 3.0, 3.0, left, topEnd, 3.0, 3.0,

		right, top, 3.0, 3.0, right, topEnd, 3.0, 3.0, rightEnd, top, 3.0, 3.0,
		rightEnd, top, 3.0, 3.0, right, topEnd, 3.0, 3.0, rightEnd, topEnd, 3.0, 3.0,

		right, bottom, 3.0, 3.0, rightEnd, bottom, 3.0, 3.0, right, bottomEnd, 3.0, 3.0,
		right, bottomEnd, 3.0, 3.0, rightEnd, bottom, 3.0, 3.0, rightEnd, bottomEnd, 3.0, 3.0,

		left, bottomEnd, 3.0, 3.0, leftEnd, bottom, 3.0, 3.0, left, bottom, 3.0, 3.0,
		left, bottomEnd, 3.0, 3.0, leftEnd, bottomEnd, 3.0, 3.0, leftEnd, bottom, 3.0, 3.0,
	}
}

// Render renders the grid.
func (grid *MapGrid) Render(columns, rows int, mapper TileMapper) {
	gl := grid.context.OpenGL

	grid.vao.OnShader(func() {
		modelMatrix := mgl.Ident4()
		lineColor := [4]float32{0.0, 0.8, 0.0, 1.0}

		grid.viewMatrixUniform.Set(gl, grid.context.ViewMatrix)
		grid.projectionMatrixUniform.Set(gl, &grid.context.ProjectionMatrix)

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

		for y := 0; y < rows; y++ {
			for x := 0; x < columns; x++ {
				tileType, slopeControl, wallHeights := mapper.MapGridInfo(level.TilePosition{X: byte(x), Y: byte(y)})
				if tileType == level.TileTypeSolid {
					continue
				}

				modelMatrix = mgl.Ident4().
					Mul4(mgl.Translate3D((float32(x)+0.5)*level.FineCoordinatesPerTileSide, (float32(y)+0.5)*level.FineCoordinatesPerTileSide, 0.0)).
					Mul4(mgl.Scale3D(level.FineCoordinatesPerTileSide, level.FineCoordinatesPerTileSide, 1.0))
				grid.modelMatrixUniform.Set(gl, &modelMatrix)
				var alphaMatrix mgl.Mat4
				alphaMatrix[15] = 1.0 // for tick marks

				gl.EnableVertexAttribArray(uint32(grid.vertexPositionAttrib))
				gl.BindBuffer(opengl.ARRAY_BUFFER, grid.vertexPositionBuffer)
				gl.VertexAttribOffset(uint32(grid.vertexPositionAttrib), 4, opengl.FLOAT, false, 0, 0)

				alphaMatrix[0+0] = heightFactor[int(wallHeights.North[0])+32]
				alphaMatrix[0+1] = heightFactor[int(wallHeights.North[1])+32]
				alphaMatrix[0+2] = heightFactor[int(wallHeights.North[2])+32]

				alphaMatrix[4+0] = heightFactor[int(wallHeights.East[0])+32]
				alphaMatrix[4+1] = heightFactor[int(wallHeights.East[1])+32]
				alphaMatrix[4+2] = heightFactor[int(wallHeights.East[2])+32]

				alphaMatrix[8+0] = heightFactor[int(wallHeights.South[0])+32]
				alphaMatrix[8+1] = heightFactor[int(wallHeights.South[1])+32]
				alphaMatrix[8+2] = heightFactor[int(wallHeights.South[2])+32]

				alphaMatrix[12+0] = heightFactor[int(wallHeights.West[0])+32]
				alphaMatrix[12+1] = heightFactor[int(wallHeights.West[1])+32]
				alphaMatrix[12+2] = heightFactor[int(wallHeights.West[2])+32]
				if tileType == level.TileTypeDiagonalOpenNorthEast || tileType == level.TileTypeDiagonalOpenSouthWest {
					alphaMatrix[0+3] = 1.0
				}
				if tileType == level.TileTypeDiagonalOpenNorthWest || tileType == level.TileTypeDiagonalOpenSouthEast {
					alphaMatrix[4+3] = 1.0
				}

				grid.colorUniform.Set(gl, &lineColor)
				grid.alphaMatrixUniform.Set(gl, &alphaMatrix)
				gl.DrawArrays(opengl.LINES, 0, 28)

				grid.renderSlopeTicks(gl, tileType, slopeControl)
			}
		}
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	})
}

func (grid *MapGrid) renderSlopeTicks(gl opengl.OpenGL, tileType level.TileType, slopeControl level.TileSlopeControl) {
	var floorTicks []int
	var ceilingTicks []int

	if slopeControl != level.TileSlopeControlFloorFlat {
		floorTicks = grid.slopeTicksByType[tileType]
	}
	if slopeControl == level.TileSlopeControlCeilingMirrored {
		ceilingTicks = grid.slopeTicksByType[tileType]
	} else if slopeControl != level.TileSlopeControlCeilingFlat {
		ceilingTicks = grid.slopeTicksByType[tileType.Info().SlopeInvertedType]
	}

	gl.EnableVertexAttribArray(uint32(grid.vertexPositionAttrib))
	gl.BindBuffer(opengl.ARRAY_BUFFER, grid.tickVertexPositionBuffer)
	gl.VertexAttribOffset(uint32(grid.vertexPositionAttrib), 4, opengl.FLOAT, false, 0, 0)

	floorTickColor := [4]float32{0.0, 0.8, 0.0, 0.9}
	grid.colorUniform.Set(gl, &floorTickColor)
	for _, index := range floorTicks {
		gl.DrawArrays(opengl.TRIANGLES, grid.floorTickStarts[index], 3)
	}
	ceilingTickColor := [4]float32{0.8, 0.0, 0.0, 0.9}
	grid.colorUniform.Set(gl, &ceilingTickColor)
	for _, index := range ceilingTicks {
		gl.DrawArrays(opengl.TRIANGLES, grid.ceilingTickStarts[index], 3)
	}
}
