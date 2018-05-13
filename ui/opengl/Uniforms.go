package opengl

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

// UniformLocation represents a uniform parameter for a shader program
type UniformLocation int32

// Matrix4Uniform represents one uniform parameter for a 4x4 matrix.
type Matrix4Uniform UniformLocation

// Set stores the provided matrix as the current uniform value.
func (location Matrix4Uniform) Set(gl OpenGl, matrix *mgl.Mat4) {
	matrixArray := (*[16]float32)(matrix)
	gl.UniformMatrix4fv(int32(location), false, matrixArray)
}

// Vector4Uniform represents one uniform parameter for a 4-dimensional vector.
type Vector4Uniform UniformLocation

// Set stores the provided 4-dimensional vector as the current uniform value.
func (location Vector4Uniform) Set(gl OpenGl, value *[4]float32) {
	gl.Uniform4fv(int32(location), value)
}
