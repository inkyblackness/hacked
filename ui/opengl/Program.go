package opengl

import "fmt"

// LinkNewProgram creates a new shader program based on the provided shaders.
func LinkNewProgram(gl OpenGl, shaders ...uint32) (program uint32, err error) {
	program = gl.CreateProgram()

	for _, shader := range shaders {
		gl.AttachShader(program, shader)
	}
	gl.LinkProgram(program)

	if gl.GetProgramParameter(program, LINK_STATUS) == 0 {
		err = fmt.Errorf("%v", gl.GetProgramInfoLog(program))
		gl.DeleteProgram(program)
		program = 0
	}

	return
}

// LinkNewStandardProgram creates a new shader based on two shader sources.
func LinkNewStandardProgram(gl OpenGl, vertexShaderSource, fragmentShaderSource string) (program uint32, err error) {
	vertexShader, vertexErr := CompileNewShader(gl, VERTEX_SHADER, vertexShaderSource)
	defer gl.DeleteShader(vertexShader)
	fragmentShader, fragmentErr := CompileNewShader(gl, FRAGMENT_SHADER, fragmentShaderSource)
	defer gl.DeleteShader(fragmentShader)

	if (vertexErr == nil) && (fragmentErr == nil) {
		program, err = LinkNewProgram(gl, vertexShader, fragmentShader)
	} else {
		err = fmt.Errorf("vertexShader: %v\nfragmentShader: %v", vertexErr, fragmentErr)
	}

	return
}
