package opengl

// LinkNewProgram creates a new shader program based on the provided shaders.
func LinkNewProgram(gl OpenGL, shaders ...uint32) (program uint32, err error) {
	program = gl.CreateProgram()

	for _, shader := range shaders {
		gl.AttachShader(program, shader)
	}
	gl.LinkProgram(program)

	if gl.GetProgramParameter(program, LINK_STATUS) == 0 {
		err = ShaderError{Log: gl.GetProgramInfoLog(program)}
		gl.DeleteProgram(program)
		program = 0
	}

	return
}

// LinkNewStandardProgram creates a new shader based on two shader sources.
func LinkNewStandardProgram(gl OpenGL, vertexShaderSource, fragmentShaderSource string) (program uint32, err error) {
	vertexShader, vertexErr := CompileNewShader(gl, VERTEX_SHADER, vertexShaderSource)
	if vertexErr != nil {
		return 0, NamedShaderError{Name: "vertexShader", Nested: vertexErr}
	}
	defer gl.DeleteShader(vertexShader)
	fragmentShader, fragmentErr := CompileNewShader(gl, FRAGMENT_SHADER, fragmentShaderSource)
	if fragmentErr != nil {
		return 0, NamedShaderError{Name: "fragmentShader", Nested: fragmentErr}
	}
	defer gl.DeleteShader(fragmentShader)

	return LinkNewProgram(gl, vertexShader, fragmentShader)
}
