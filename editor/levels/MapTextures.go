package levels

import (
	"fmt"
	"math"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/hacked/editor/graphics"
	"github.com/inkyblackness/hacked/editor/render"
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ui/opengl"
)

var mapTexturesVertexShaderSource = `
#version 150
precision mediump float;

in vec3 vertexPosition;

uniform mat4 modelMatrix;
uniform mat4 viewMatrix;
uniform mat4 projectionMatrix;
uniform mat4 uvMatrix;

out vec2 uv;

void main(void) {
	gl_Position = projectionMatrix * viewMatrix * modelMatrix * vec4(vertexPosition, 1.0);

	uv = (uvMatrix * vec4(vertexPosition, 1.0)).xy;
}
`

var mapTexturesFragmentShaderSource = `
#version 150
precision mediump float;

in vec2 uv;

uniform sampler2D palette;
uniform sampler2D bitmap;

out vec4 fragColor;

void main(void) {
	vec4 pixel = texture(bitmap, uv);

	fragColor = texture(palette, vec2(pixel.r, 0.5));
}
`

// TextureQuery is a getter function to retrieve the texture for the given
// level texture index.
type TextureQuery func(index int) (*graphics.BitmapTexture, error)

// TileTextureQuery is a getter function to retrieve properties for rendering a texture of a tile.
type TileTextureQuery func(x, y int) (tileType level.TileType, textureIndex int, textureRotations int)

// MapTextures is a renderable for textures.
type MapTextures struct {
	context *render.Context

	program                 uint32
	vao                     *opengl.VertexArrayObject
	vertexPositionBuffer    uint32
	vertexPositionAttrib    int32
	modelMatrixUniform      opengl.Matrix4Uniform
	viewMatrixUniform       opengl.Matrix4Uniform
	projectionMatrixUniform opengl.Matrix4Uniform
	uvMatrixUniform         opengl.Matrix4Uniform

	paletteUniform int32
	bitmapUniform  int32

	textureQuery TextureQuery

	lastTileType level.TileType
}

var uvRotations map[int]*mgl.Mat4

func init() {
	uvRotations = make(map[int]*mgl.Mat4)
	for i := 0; i < 4; i++ {
		matrix := mgl.Translate3D(0.5, 0.5, 0.0).
			Mul4(mgl.HomogRotate3DZ(float32(math.Pi * float32(i) / -2.0))).
			Mul4(mgl.Translate3D(-0.5, -0.5, 0.0)).
			Mul4(mgl.Scale3D(1.0, -1.0, 1.0))
		uvRotations[i] = &matrix
	}
}

// NewMapTextures returns a new instance of a renderable for tile map textures.
func NewMapTextures(context *render.Context, textureQuery TextureQuery) *MapTextures {
	gl := context.OpenGL
	program, programErr := opengl.LinkNewStandardProgram(gl, mapTexturesVertexShaderSource, mapTexturesFragmentShaderSource)

	if programErr != nil {
		panic(fmt.Errorf("MapTextures shader failed: %v", programErr))
	}
	renderable := &MapTextures{
		context: context,
		program: program,

		vao:                     opengl.NewVertexArrayObject(gl, program),
		vertexPositionBuffer:    gl.GenBuffers(1)[0],
		vertexPositionAttrib:    gl.GetAttribLocation(program, "vertexPosition"),
		modelMatrixUniform:      opengl.Matrix4Uniform(gl.GetUniformLocation(program, "modelMatrix")),
		viewMatrixUniform:       opengl.Matrix4Uniform(gl.GetUniformLocation(program, "viewMatrix")),
		projectionMatrixUniform: opengl.Matrix4Uniform(gl.GetUniformLocation(program, "projectionMatrix")),
		uvMatrixUniform:         opengl.Matrix4Uniform(gl.GetUniformLocation(program, "uvMatrix")),
		paletteUniform:          gl.GetUniformLocation(program, "palette"),
		bitmapUniform:           gl.GetUniformLocation(program, "bitmap"),
		textureQuery:            textureQuery,
		lastTileType:            level.TileTypeSolid,
	}

	renderable.vao.WithSetter(func(gl opengl.OpenGL) {
		gl.EnableVertexAttribArray(uint32(renderable.vertexPositionAttrib))
		gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexPositionBuffer)
		gl.VertexAttribOffset(uint32(renderable.vertexPositionAttrib), 3, opengl.FLOAT, false, 0, 0)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)
	})

	return renderable
}

// Dispose releases any internal resources
func (renderable *MapTextures) Dispose() {
	gl := renderable.context.OpenGL

	renderable.vao.Dispose()
	gl.DeleteProgram(renderable.program)
	gl.DeleteBuffers([]uint32{renderable.vertexPositionBuffer})
}

// Render renders
func (renderable *MapTextures) Render(columns, rows int, tileTextureQuery TileTextureQuery, paletteTexture *graphics.PaletteTexture) {
	gl := renderable.context.OpenGL

	renderable.vao.OnShader(func() {
		renderable.viewMatrixUniform.Set(gl, renderable.context.ViewMatrix)
		renderable.projectionMatrixUniform.Set(gl, &renderable.context.ProjectionMatrix)

		textureUnit := int32(0)
		gl.ActiveTexture(opengl.TEXTURE0 + uint32(textureUnit))
		gl.BindTexture(opengl.TEXTURE_2D, paletteTexture.Handle())
		gl.Uniform1i(renderable.paletteUniform, textureUnit)

		textureUnit = 1
		gl.ActiveTexture(opengl.TEXTURE0 + uint32(textureUnit))

		scaling := mgl.Scale3D(fineCoordinatesPerTileSide, fineCoordinatesPerTileSide, 1.0)
		for y := 0; y < rows; y++ {
			for x := 0; x < columns; x++ {
				tileType, textureIndex, textureRotations := tileTextureQuery(x, y)
				if tileType != level.TileTypeSolid {
					texture, _ := renderable.textureQuery(textureIndex)
					if texture != nil {
						modelMatrix := mgl.Translate3D(float32(x)*fineCoordinatesPerTileSide, float32(y)*fineCoordinatesPerTileSide, 0.0).
							Mul4(scaling)

						uvMatrix := uvRotations[textureRotations]
						renderable.uvMatrixUniform.Set(gl, uvMatrix)
						renderable.modelMatrixUniform.Set(gl, &modelMatrix)
						vertexCount := renderable.ensureTileType(tileType)
						gl.BindTexture(opengl.TEXTURE_2D, texture.Handle())
						gl.Uniform1i(renderable.bitmapUniform, textureUnit)

						gl.DrawArrays(opengl.TRIANGLES, 0, int32(vertexCount))
					}
				}
			}
		}

		gl.BindTexture(opengl.TEXTURE_2D, 0)
	})
}

func (renderable *MapTextures) ensureTileType(tileType level.TileType) (vertexCount int) {
	displayedType := level.TileTypeOpen

	vertexCount = 6
	if tileType == level.TileTypeDiagonalOpenNorthEast || tileType == level.TileTypeDiagonalOpenNorthWest ||
		tileType == level.TileTypeDiagonalOpenSouthEast || tileType == level.TileTypeDiagonalOpenSouthWest {
		displayedType = tileType
		vertexCount = 3
	}
	if renderable.lastTileType != displayedType {
		gl := renderable.context.OpenGL
		var vertices []float32
		limit := float32(1.0)

		switch displayedType {
		case level.TileTypeDiagonalOpenNorthEast:
			vertices = []float32{
				0.0, limit, 0.0,
				limit, limit, 0.0,
				limit, 0.0, 0.0}
		case level.TileTypeDiagonalOpenNorthWest:
			vertices = []float32{
				0.0, limit, 0.0,
				limit, limit, 0.0,
				0.0, 0.0, 0.0}
		case level.TileTypeDiagonalOpenSouthEast:
			vertices = []float32{
				limit, limit, 0.0,
				limit, 0.0, 0.0,
				0.0, 0.0, 0.0}
		case level.TileTypeDiagonalOpenSouthWest:
			vertices = []float32{
				0.0, limit, 0.0,
				limit, 0.0, 0.0,
				0.0, 0.0, 0.0}
		case level.TileTypeOpen:
			vertices = []float32{
				0.0, 0.0, 0.0,
				limit, 0.0, 0.0,
				limit, limit, 0.0,

				limit, limit, 0.0,
				0.0, limit, 0.0,
				0.0, 0.0, 0.0}
		}
		gl.BindBuffer(opengl.ARRAY_BUFFER, renderable.vertexPositionBuffer)
		gl.BufferData(opengl.ARRAY_BUFFER, len(vertices)*4, vertices, opengl.STATIC_DRAW)
		gl.BindBuffer(opengl.ARRAY_BUFFER, 0)

		renderable.lastTileType = displayedType
	}

	return
}
