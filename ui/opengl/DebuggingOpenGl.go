package opengl

type debuggingOpenGl struct {
	gl OpenGl

	onEntry DebuggingEntryFunc
	onExit  DebuggingExitFunc
	onError DebuggingErrorFunc
}

func (debugging *debuggingOpenGl) recordEntry(name string, param ...interface{}) {
	debugging.onEntry(name, param...)
}

func (debugging *debuggingOpenGl) recordExit(name string, result ...interface{}) {
	var errorCodes []uint32

	for errorCode := debugging.gl.GetError(); errorCode != NO_ERROR; errorCode = debugging.gl.GetError() {
		errorCodes = append(errorCodes, errorCode)
	}
	debugging.onExit(name, result...)
	if len(errorCodes) > 0 {
		debugging.onError(name, errorCodes)
	}
}

// ActiveTexture implements the opengl.OpenGl interface.
func (debugging *debuggingOpenGl) ActiveTexture(texture uint32) {
	debugging.recordEntry("ActiveTexture", texture)
	debugging.gl.ActiveTexture(texture)
	debugging.recordExit("ActiveTexture")
}

// AttachShader implements the OpenGl interface.
func (debugging *debuggingOpenGl) AttachShader(program uint32, shader uint32) {
	debugging.recordEntry("AttachShader", program, shader)
	debugging.gl.AttachShader(program, shader)
	debugging.recordExit("AttachShader")
}

// BindAttribLocation implements the OpenGl interface.
func (debugging *debuggingOpenGl) BindAttribLocation(program uint32, index uint32, name string) {
	debugging.recordEntry("BindAttribLocation", program, index, name)
	debugging.gl.BindAttribLocation(program, index, name)
	debugging.recordExit("BindAttribLocation")
}

// BindBuffer implements the OpenGl interface.
func (debugging *debuggingOpenGl) BindBuffer(target uint32, buffer uint32) {
	debugging.recordEntry("BindBuffer", target, buffer)
	debugging.gl.BindBuffer(target, buffer)
	debugging.recordExit("BindBuffer")
}

// BindTexture implements the opengl.OpenGl interface.
func (debugging *debuggingOpenGl) BindTexture(target uint32, texture uint32) {
	debugging.recordEntry("BindTexture", target, texture)
	debugging.gl.BindTexture(target, texture)
	debugging.recordExit("BindTexture")
}

// BindVertexArray implements the OpenGl interface.
func (debugging *debuggingOpenGl) BindVertexArray(array uint32) {
	debugging.recordEntry("BindVertexArray", array)
	debugging.gl.BindVertexArray(array)
	debugging.recordExit("BindVertexArray")
}

// BlendFunc implements the OpenGl interface.
func (debugging *debuggingOpenGl) BlendFunc(sfactor uint32, dfactor uint32) {
	debugging.recordEntry("BlendFunc", sfactor, dfactor)
	debugging.gl.BlendFunc(sfactor, dfactor)
	debugging.recordExit("BlendFunc")
}

// BufferData implements the OpenGl interface.
func (debugging *debuggingOpenGl) BufferData(target uint32, size int, data interface{}, usage uint32) {
	debugging.recordEntry("BufferData", target, size, data, usage)
	debugging.gl.BufferData(target, size, data, usage)
	debugging.recordExit("BufferData")
}

// Clear implements the OpenGl interface.
func (debugging *debuggingOpenGl) Clear(mask uint32) {
	debugging.recordEntry("Clear", mask)
	debugging.gl.Clear(mask)
	debugging.recordExit("Clear")
}

// ClearColor implements the OpenGl interface.
func (debugging *debuggingOpenGl) ClearColor(red float32, green float32, blue float32, alpha float32) {
	debugging.recordEntry("ClearColor", red, green, blue, alpha)
	debugging.gl.ClearColor(red, green, blue, alpha)
	debugging.recordExit("ClearColor")
}

// CompileShader implements the OpenGl interface.
func (debugging *debuggingOpenGl) CompileShader(shader uint32) {
	debugging.recordEntry("CompileShader", shader)
	debugging.gl.CompileShader(shader)
	debugging.recordExit("CompileShader")
}

// CreateProgram implements the OpenGl interface.
func (debugging *debuggingOpenGl) CreateProgram() uint32 {
	debugging.recordEntry("CreateProgram")
	result := debugging.gl.CreateProgram()
	debugging.recordExit("CreateProgram", result)
	return result
}

// CreateShader implements the OpenGl interface.
func (debugging *debuggingOpenGl) CreateShader(shaderType uint32) uint32 {
	debugging.recordEntry("CreateShader", shaderType)
	result := debugging.gl.CreateShader(shaderType)
	debugging.recordExit("CreateShader", result)
	return result
}

// DeleteBuffers implements the OpenGl interface.
func (debugging *debuggingOpenGl) DeleteBuffers(buffers []uint32) {
	debugging.recordEntry("DeleteBuffers", buffers)
	debugging.gl.DeleteBuffers(buffers)
	debugging.recordExit("DeleteBuffers")
}

// DeleteProgram implements the OpenGl interface.
func (debugging *debuggingOpenGl) DeleteProgram(program uint32) {
	debugging.recordEntry("DeleteProgram", program)
	debugging.gl.DeleteProgram(program)
	debugging.recordExit("DeleteProgram")
}

// DeleteShader implements the OpenGl interface.
func (debugging *debuggingOpenGl) DeleteShader(shader uint32) {
	debugging.recordEntry("DeleteShader", shader)
	debugging.gl.DeleteShader(shader)
	debugging.recordExit("DeleteShader")
}

// DeleteTextures implements the OpenGl interface.
func (debugging *debuggingOpenGl) DeleteTextures(textures []uint32) {
	debugging.recordEntry("DeleteTextures", textures)
	debugging.gl.DeleteTextures(textures)
	debugging.recordExit("DeleteTextures")
}

// DeleteVertexArrays implements the OpenGl interface.
func (debugging *debuggingOpenGl) DeleteVertexArrays(arrays []uint32) {
	debugging.recordEntry("DeleteVertexArrays", arrays)
	debugging.gl.DeleteVertexArrays(arrays)
	debugging.recordExit("DeleteVertexArrays")
}

// Disable implements the OpenGl interface.
func (debugging *debuggingOpenGl) Disable(cap uint32) {
	debugging.recordEntry("Disable", cap)
	debugging.gl.Disable(cap)
	debugging.recordExit("Disable")
}

// DrawArrays implements the OpenGl interface.
func (debugging *debuggingOpenGl) DrawArrays(mode uint32, first int32, count int32) {
	debugging.recordEntry("DrawArrays", first, count)
	debugging.gl.DrawArrays(mode, first, count)
	debugging.recordExit("DrawArrays")
}

// Enable implements the OpenGl interface.
func (debugging *debuggingOpenGl) Enable(cap uint32) {
	debugging.recordEntry("Enable", cap)
	debugging.gl.Enable(cap)
	debugging.recordExit("Enable")
}

// EnableVertexAttribArray implements the OpenGl interface.
func (debugging *debuggingOpenGl) EnableVertexAttribArray(index uint32) {
	debugging.recordEntry("EnableVertexAttribArray", index)
	debugging.gl.EnableVertexAttribArray(index)
	debugging.recordExit("EnableVertexAttribArray")
}

// GenerateMipmap implements the opengl.OpenGl interface.
func (debugging *debuggingOpenGl) GenerateMipmap(target uint32) {
	debugging.recordEntry("GenerateMipmap", target)
	debugging.gl.GenerateMipmap(target)
	debugging.recordExit("GenerateMipmap")
}

// GenBuffers implements the OpenGl interface.
func (debugging *debuggingOpenGl) GenBuffers(n int32) []uint32 {
	debugging.recordEntry("GenBuffers", n)
	result := debugging.gl.GenBuffers(n)
	debugging.recordExit("GenBuffers", result)
	return result
}

// GenTextures implements the opengl.OpenGl interface.
func (debugging *debuggingOpenGl) GenTextures(n int32) []uint32 {
	debugging.recordEntry("GenTextures", n)
	result := debugging.gl.GenTextures(n)
	debugging.recordExit("GenTextures", result)
	return result
}

// GenVertexArrays implements the OpenGl interface.
func (debugging *debuggingOpenGl) GenVertexArrays(n int32) []uint32 {
	debugging.recordEntry("GenVertexArrays", n)
	result := debugging.gl.GenVertexArrays(n)
	debugging.recordExit("GenVertexArrays", result)
	return result
}

// GetAttribLocation implements the OpenGl interface.
func (debugging *debuggingOpenGl) GetAttribLocation(program uint32, name string) int32 {
	debugging.recordEntry("GetAttribLocation", program, name)
	result := debugging.gl.GetAttribLocation(program, name)
	debugging.recordExit("GetAttribLocation", result)
	return result
}

// GetError implements the OpenGl interface.
func (debugging *debuggingOpenGl) GetError() uint32 {
	debugging.recordEntry("GetError")
	result := debugging.gl.GetError()
	debugging.recordExit("GetError", result)
	return result
}

// GetProgramInfoLog implements the OpenGl interface.
func (debugging *debuggingOpenGl) GetProgramInfoLog(program uint32) string {
	debugging.recordEntry("GetProgramInfoLog", program)
	result := debugging.gl.GetProgramInfoLog(program)
	debugging.recordExit("GetProgramInfoLog", result)
	return result
}

// GetProgramParameter implements the OpenGl interface.
func (debugging *debuggingOpenGl) GetProgramParameter(program uint32, param uint32) int32 {
	debugging.recordEntry("GetProgramParameter", program, param)
	result := debugging.gl.GetProgramParameter(program, param)
	debugging.recordExit("GetProgramParameter", result)
	return result
}

// GetShaderInfoLog implements the OpenGl interface.
func (debugging *debuggingOpenGl) GetShaderInfoLog(shader uint32) string {
	debugging.recordEntry("GetShaderInfoLog", shader)
	result := debugging.gl.GetShaderInfoLog(shader)
	debugging.recordExit("GetShaderInfoLog", result)
	return result
}

// GetShaderParameter implements the OpenGl interface.
func (debugging *debuggingOpenGl) GetShaderParameter(shader uint32, param uint32) int32 {
	debugging.recordEntry("GetShaderParameter", shader, param)
	result := debugging.gl.GetShaderParameter(shader, param)
	debugging.recordExit("GetShaderParameter", result)
	return result
}

// GetUniformLocation implements the OpenGl interface.
func (debugging *debuggingOpenGl) GetUniformLocation(program uint32, name string) int32 {
	debugging.recordEntry("GetUniformLocation", program, name)
	result := debugging.gl.GetUniformLocation(program, name)
	debugging.recordExit("GetUniformLocation", result)
	return result
}

// LinkProgram implements the OpenGl interface.
func (debugging *debuggingOpenGl) LinkProgram(program uint32) {
	debugging.recordEntry("LinkProgram", program)
	debugging.gl.LinkProgram(program)
	debugging.recordExit("LinkProgram")
}

// ReadPixels implements the OpenGl interface.
func (debugging *debuggingOpenGl) ReadPixels(x int32, y int32, width int32, height int32, format uint32, pixelType uint32, pixels interface{}) {
	debugging.recordEntry("ReadPixels", x, y, width, height, format, pixelType, pixels)
	debugging.gl.ReadPixels(x, y, width, height, format, pixelType, pixels)
	debugging.recordExit("ReadPixels")
}

// ShaderSource implements the OpenGl interface.
func (debugging *debuggingOpenGl) ShaderSource(shader uint32, source string) {
	debugging.recordEntry("ShaderSource", shader, source)
	debugging.gl.ShaderSource(shader, source)
	debugging.recordExit("ShaderSource")
}

// TexImage2D implements the opengl.OpenGl interface.
func (debugging *debuggingOpenGl) TexImage2D(target uint32, level int32, internalFormat uint32, width int32, height int32,
	border int32, format uint32, xtype uint32, pixels interface{}) {
	debugging.recordEntry("TexImage2D", target, level, internalFormat, width, height, border, format, xtype, pixels)
	debugging.gl.TexImage2D(target, level, internalFormat, width, height, border, format, xtype, pixels)
	debugging.recordExit("TexImage2D")
}

// TexParameteri implements the opengl.OpenGl interface.
func (debugging *debuggingOpenGl) TexParameteri(target uint32, pname uint32, param int32) {
	debugging.recordEntry("TexParameteri", target, pname, param)
	debugging.gl.TexParameteri(target, pname, param)
	debugging.recordExit("TexParameteri")
}

// Uniform1i implements the opengl.OpenGl interface.
func (debugging *debuggingOpenGl) Uniform1i(location int32, value int32) {
	debugging.recordEntry("Uniform1i", location, value)
	debugging.gl.Uniform1i(location, value)
	debugging.recordExit("Uniform1i")
}

// Uniform4fv implements the opengl.OpenGl interface.
func (debugging *debuggingOpenGl) Uniform4fv(location int32, value *[4]float32) {
	debugging.recordEntry("Uniform4fv", location, value)
	debugging.gl.Uniform4fv(location, value)
	debugging.recordExit("Uniform4fv")
}

// UniformMatrix4fv implements the OpenGl interface.
func (debugging *debuggingOpenGl) UniformMatrix4fv(location int32, transpose bool, value *[16]float32) {
	debugging.recordEntry("UniformMatrix4fv", location, transpose, value)
	debugging.gl.UniformMatrix4fv(location, transpose, value)
	debugging.recordExit("UniformMatrix4fv")
}

// UseProgram implements the OpenGl interface.
func (debugging *debuggingOpenGl) UseProgram(program uint32) {
	debugging.recordEntry("UseProgram", program)
	debugging.gl.UseProgram(program)
	debugging.recordExit("UseProgram")
}

// VertexAttribOffset implements the OpenGl interface.
func (debugging *debuggingOpenGl) VertexAttribOffset(index uint32, size int32, attribType uint32, normalized bool, stride int32, offset int) {
	debugging.recordEntry("VertexAttribOffset", index, size, attribType, normalized, stride, offset)
	debugging.gl.VertexAttribOffset(index, size, attribType, normalized, stride, offset)
	debugging.recordExit("VertexAttribOffset")
}

// Viewport implements the OpenGl interface.
func (debugging *debuggingOpenGl) Viewport(x int32, y int32, width int32, height int32) {
	debugging.recordEntry("Viewport", x, y, width, height)
	debugging.gl.Viewport(x, y, width, height)
	debugging.recordExit("Viewport")
}
