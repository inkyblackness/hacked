package opengl

// CompileNewShader creates a shader of given type and compiles the provided source.
func CompileNewShader(gl OpenGL, shaderType uint32, source string) (shader uint32, err error) {
	shader = gl.CreateShader(shaderType)

	gl.ShaderSource(shader, source)
	gl.CompileShader(shader)

	compileStatus := gl.GetShaderParameter(shader, COMPILE_STATUS)
	if compileStatus == 0 {
		err = ShaderError{Log: gl.GetShaderInfoLog(shader)}
		gl.DeleteShader(shader)
		shader = 0
	}

	return
}
