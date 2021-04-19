package opengl

type debuggingOpenGL struct {
	gl OpenGL

	onEntry DebuggingEntryFunc
	onExit  DebuggingExitFunc
	onError DebuggingErrorFunc
}

func (debugging *debuggingOpenGL) recordEntry(name string, param ...interface{}) {
	debugging.onEntry(name, param...)
}

func (debugging *debuggingOpenGL) recordExit(name string, result ...interface{}) {
	var errorCodes []uint32

	for errorCode := debugging.gl.GetError(); errorCode != NO_ERROR; errorCode = debugging.gl.GetError() {
		errorCodes = append(errorCodes, errorCode)
	}
	debugging.onExit(name, result...)
	if len(errorCodes) > 0 {
		debugging.onError(name, errorCodes)
	}
}

// ActiveTexture implements the opengl.OpenGL interface.
func (debugging *debuggingOpenGL) ActiveTexture(texture uint32) {
	debugging.recordEntry("ActiveTexture", texture)
	debugging.gl.ActiveTexture(texture)
	debugging.recordExit("ActiveTexture")
}

// AttachShader implements the OpenGL interface.
func (debugging *debuggingOpenGL) AttachShader(program uint32, shader uint32) {
	debugging.recordEntry("AttachShader", program, shader)
	debugging.gl.AttachShader(program, shader)
	debugging.recordExit("AttachShader")
}

// BindAttribLocation implements the OpenGL interface.
func (debugging *debuggingOpenGL) BindAttribLocation(program uint32, index uint32, name string) {
	debugging.recordEntry("BindAttribLocation", program, index, name)
	debugging.gl.BindAttribLocation(program, index, name)
	debugging.recordExit("BindAttribLocation")
}

// BindBuffer implements the OpenGL interface.
func (debugging *debuggingOpenGL) BindBuffer(target uint32, buffer uint32) {
	debugging.recordEntry("BindBuffer", target, buffer)
	debugging.gl.BindBuffer(target, buffer)
	debugging.recordExit("BindBuffer")
}

// BindFramebuffer implements the OpenGL interface.
func (debugging *debuggingOpenGL) BindFramebuffer(target uint32, buffer uint32) {
	debugging.recordEntry("BindFramebuffer", target, buffer)
	debugging.gl.BindFramebuffer(target, buffer)
	debugging.recordExit("BindFramebuffer")
}

// BindRenderbuffer implements the OpenGL interface.
func (debugging *debuggingOpenGL) BindRenderbuffer(target uint32, buffer uint32) {
	debugging.recordEntry("BindRenderbuffer", target, buffer)
	debugging.gl.BindRenderbuffer(target, buffer)
	debugging.recordExit("BindRenderbuffer")
}

// BindSampler implements the OpenGL interface.
func (debugging *debuggingOpenGL) BindSampler(unit uint32, sampler uint32) {
	debugging.recordEntry("BindSampler", unit, sampler)
	debugging.gl.BindSampler(unit, sampler)
	debugging.recordExit("BindSampler")
}

// BindTexture implements the OpenGL interface.
func (debugging *debuggingOpenGL) BindTexture(target uint32, texture uint32) {
	debugging.recordEntry("BindTexture", target, texture)
	debugging.gl.BindTexture(target, texture)
	debugging.recordExit("BindTexture")
}

// BindVertexArray implements the OpenGL interface.
func (debugging *debuggingOpenGL) BindVertexArray(array uint32) {
	debugging.recordEntry("BindVertexArray", array)
	debugging.gl.BindVertexArray(array)
	debugging.recordExit("BindVertexArray")
}

// BlendEquation implements the OpenGL interface.
func (debugging *debuggingOpenGL) BlendEquation(mode uint32) {
	debugging.recordEntry("BlendEquation", mode)
	debugging.gl.BlendEquation(mode)
	debugging.recordExit("BlendEquation")
}

// BlendEquationSeparate implements the OpenGL interface.
func (debugging *debuggingOpenGL) BlendEquationSeparate(modeRGB uint32, modeAlpha uint32) {
	debugging.recordEntry("BlendEquationSeparate", modeRGB, modeAlpha)
	debugging.gl.BlendEquationSeparate(modeRGB, modeAlpha)
	debugging.recordExit("BlendEquationSeparate")
}

// BlendFunc implements the OpenGL interface.
func (debugging *debuggingOpenGL) BlendFunc(sfactor uint32, dfactor uint32) {
	debugging.recordEntry("BlendFunc", sfactor, dfactor)
	debugging.gl.BlendFunc(sfactor, dfactor)
	debugging.recordExit("BlendFunc")
}

// BlendFuncSeparate implements the OpenGL interface.
func (debugging *debuggingOpenGL) BlendFuncSeparate(srcRGB uint32, dstRGB uint32, srcAlpha uint32, dstAlpha uint32) {
	debugging.recordEntry("BlendFuncSeparate", srcRGB, dstRGB, srcAlpha, dstAlpha)
	debugging.gl.BlendFuncSeparate(srcRGB, dstRGB, srcAlpha, dstAlpha)
	debugging.recordExit("BlendFuncSeparate")
}

// BufferData implements the OpenGL interface.
func (debugging *debuggingOpenGL) BufferData(target uint32, size int, data interface{}, usage uint32) {
	debugging.recordEntry("BufferData", target, size, data, usage)
	debugging.gl.BufferData(target, size, data, usage)
	debugging.recordExit("BufferData")
}

// CheckFramebufferStatus implements the OpenGL interface.
func (debugging *debuggingOpenGL) CheckFramebufferStatus(target uint32) uint32 {
	debugging.recordEntry("CheckFramebufferStatus", target)
	result := debugging.gl.CheckFramebufferStatus(target)
	debugging.recordExit("CheckFramebufferStatus")
	return result
}

// Clear implements the OpenGL interface.
func (debugging *debuggingOpenGL) Clear(mask uint32) {
	debugging.recordEntry("Clear", mask)
	debugging.gl.Clear(mask)
	debugging.recordExit("Clear")
}

// ClearColor implements the OpenGL interface.
func (debugging *debuggingOpenGL) ClearColor(red float32, green float32, blue float32, alpha float32) {
	debugging.recordEntry("ClearColor", red, green, blue, alpha)
	debugging.gl.ClearColor(red, green, blue, alpha)
	debugging.recordExit("ClearColor")
}

// CompileShader implements the OpenGL interface.
func (debugging *debuggingOpenGL) CompileShader(shader uint32) {
	debugging.recordEntry("CompileShader", shader)
	debugging.gl.CompileShader(shader)
	debugging.recordExit("CompileShader")
}

// CreateProgram implements the OpenGL interface.
func (debugging *debuggingOpenGL) CreateProgram() uint32 {
	debugging.recordEntry("CreateProgram")
	result := debugging.gl.CreateProgram()
	debugging.recordExit("CreateProgram", result)
	return result
}

// CreateShader implements the OpenGL interface.
func (debugging *debuggingOpenGL) CreateShader(shaderType uint32) uint32 {
	debugging.recordEntry("CreateShader", shaderType)
	result := debugging.gl.CreateShader(shaderType)
	debugging.recordExit("CreateShader", result)
	return result
}

// DeleteBuffers implements the OpenGL interface.
func (debugging *debuggingOpenGL) DeleteBuffers(buffers []uint32) {
	debugging.recordEntry("DeleteBuffers", buffers)
	debugging.gl.DeleteBuffers(buffers)
	debugging.recordExit("DeleteBuffers")
}

// DeleteFramebuffers implements the OpenGL interface.
func (debugging *debuggingOpenGL) DeleteFramebuffers(buffers []uint32) {
	debugging.recordEntry("DeleteFramebuffers", buffers)
	debugging.gl.DeleteFramebuffers(buffers)
	debugging.recordExit("DeleteFramebuffers")
}

// DeleteProgram implements the OpenGL interface.
func (debugging *debuggingOpenGL) DeleteProgram(program uint32) {
	debugging.recordEntry("DeleteProgram", program)
	debugging.gl.DeleteProgram(program)
	debugging.recordExit("DeleteProgram")
}

// DeleteShader implements the OpenGL interface.
func (debugging *debuggingOpenGL) DeleteShader(shader uint32) {
	debugging.recordEntry("DeleteShader", shader)
	debugging.gl.DeleteShader(shader)
	debugging.recordExit("DeleteShader")
}

// DeleteTextures implements the OpenGL interface.
func (debugging *debuggingOpenGL) DeleteTextures(textures []uint32) {
	debugging.recordEntry("DeleteTextures", textures)
	debugging.gl.DeleteTextures(textures)
	debugging.recordExit("DeleteTextures")
}

// DeleteVertexArrays implements the OpenGL interface.
func (debugging *debuggingOpenGL) DeleteVertexArrays(arrays []uint32) {
	debugging.recordEntry("DeleteVertexArrays", arrays)
	debugging.gl.DeleteVertexArrays(arrays)
	debugging.recordExit("DeleteVertexArrays")
}

// Disable implements the OpenGL interface.
func (debugging *debuggingOpenGL) Disable(cap uint32) {
	debugging.recordEntry("Disable", cap)
	debugging.gl.Disable(cap)
	debugging.recordExit("Disable")
}

// DrawArrays implements the OpenGL interface.
func (debugging *debuggingOpenGL) DrawArrays(mode uint32, first int32, count int32) {
	debugging.recordEntry("DrawArrays", first, count)
	debugging.gl.DrawArrays(mode, first, count)
	debugging.recordExit("DrawArrays")
}

// DrawBuffers implements the OpenGL interface.
func (debugging *debuggingOpenGL) DrawBuffers(buffers []uint32) {
	debugging.recordEntry("DrawBuffers", buffers)
	debugging.gl.DrawBuffers(buffers)
	debugging.recordExit("DrawBuffers")
}

// DrawElements implements the OpenGL interface.
func (debugging *debuggingOpenGL) DrawElements(mode uint32, count int32, elementType uint32, indices uintptr) {
	debugging.recordEntry("DrawElements", mode, count, elementType, indices)
	debugging.gl.DrawElements(mode, count, elementType, indices)
	debugging.recordExit("DrawElements")
}

// Enable implements the OpenGL interface.
func (debugging *debuggingOpenGL) Enable(cap uint32) {
	debugging.recordEntry("Enable", cap)
	debugging.gl.Enable(cap)
	debugging.recordExit("Enable")
}

// EnableVertexAttribArray implements the OpenGL interface.
func (debugging *debuggingOpenGL) EnableVertexAttribArray(index uint32) {
	debugging.recordEntry("EnableVertexAttribArray", index)
	debugging.gl.EnableVertexAttribArray(index)
	debugging.recordExit("EnableVertexAttribArray")
}

// FramebufferRenderbuffer implements the OpenGL interface.
func (debugging *debuggingOpenGL) FramebufferRenderbuffer(target uint32, attachment uint32, renderbuffertarget uint32, renderbuffer uint32) {
	debugging.recordEntry("FramebufferRenderbuffer", target, attachment, renderbuffertarget, renderbuffer)
	debugging.gl.FramebufferRenderbuffer(target, attachment, renderbuffertarget, renderbuffer)
	debugging.recordExit("FramebufferRenderbuffer")
}

// FramebufferTexture implements the OpenGL interface.
func (debugging *debuggingOpenGL) FramebufferTexture(target uint32, attachment uint32, texture uint32, level int32) {
	debugging.recordEntry("FramebufferTexture", target, attachment, texture, level)
	debugging.gl.FramebufferTexture(target, attachment, texture, level)
	debugging.recordExit("FramebufferTexture")
}

// GenerateMipmap implements the opengl.OpenGL interface.
func (debugging *debuggingOpenGL) GenerateMipmap(target uint32) {
	debugging.recordEntry("GenerateMipmap", target)
	debugging.gl.GenerateMipmap(target)
	debugging.recordExit("GenerateMipmap")
}

// GenBuffers implements the OpenGL interface.
func (debugging *debuggingOpenGL) GenBuffers(n int32) []uint32 {
	debugging.recordEntry("GenBuffers", n)
	result := debugging.gl.GenBuffers(n)
	debugging.recordExit("GenBuffers", result)
	return result
}

// GenFramebuffers implements the OpenGL interface.
func (debugging *debuggingOpenGL) GenFramebuffers(n int32) []uint32 {
	debugging.recordEntry("GenFramebuffers", n)
	result := debugging.gl.GenFramebuffers(n)
	debugging.recordExit("GenFramebuffers", result)
	return result
}

// GenRenderbuffers implements the OpenGL interface.
func (debugging *debuggingOpenGL) GenRenderbuffers(n int32) []uint32 {
	debugging.recordEntry("GenRenderbuffers", n)
	result := debugging.gl.GenRenderbuffers(n)
	debugging.recordExit("GenRenderbuffers", result)
	return result
}

// GenTextures implements the opengl.OpenGL interface.
func (debugging *debuggingOpenGL) GenTextures(n int32) []uint32 {
	debugging.recordEntry("GenTextures", n)
	result := debugging.gl.GenTextures(n)
	debugging.recordExit("GenTextures", result)
	return result
}

// GenVertexArrays implements the OpenGL interface.
func (debugging *debuggingOpenGL) GenVertexArrays(n int32) []uint32 {
	debugging.recordEntry("GenVertexArrays", n)
	result := debugging.gl.GenVertexArrays(n)
	debugging.recordExit("GenVertexArrays", result)
	return result
}

// GetAttribLocation implements the OpenGL interface.
func (debugging *debuggingOpenGL) GetAttribLocation(program uint32, name string) int32 {
	debugging.recordEntry("GetAttribLocation", program, name)
	result := debugging.gl.GetAttribLocation(program, name)
	debugging.recordExit("GetAttribLocation", result)
	return result
}

// GetError implements the OpenGL interface.
func (debugging *debuggingOpenGL) GetError() uint32 {
	debugging.recordEntry("GetError")
	result := debugging.gl.GetError()
	debugging.recordExit("GetError", result)
	return result
}

// GetIntegerv implements the OpenGL interface.
func (debugging *debuggingOpenGL) GetIntegerv(name uint32, data *int32) {
	debugging.recordEntry("GetIntegerv", name, data)
	debugging.gl.GetIntegerv(name, data)
	debugging.recordExit("GetIntegerv")
}

// GetProgramInfoLog implements the OpenGL interface.
func (debugging *debuggingOpenGL) GetProgramInfoLog(program uint32) string {
	debugging.recordEntry("GetProgramInfoLog", program)
	result := debugging.gl.GetProgramInfoLog(program)
	debugging.recordExit("GetProgramInfoLog", result)
	return result
}

// GetProgramParameter implements the OpenGL interface.
func (debugging *debuggingOpenGL) GetProgramParameter(program uint32, param uint32) int32 {
	debugging.recordEntry("GetProgramParameter", program, param)
	result := debugging.gl.GetProgramParameter(program, param)
	debugging.recordExit("GetProgramParameter", result)
	return result
}

// GetShaderInfoLog implements the OpenGL interface.
func (debugging *debuggingOpenGL) GetShaderInfoLog(shader uint32) string {
	debugging.recordEntry("GetShaderInfoLog", shader)
	result := debugging.gl.GetShaderInfoLog(shader)
	debugging.recordExit("GetShaderInfoLog", result)
	return result
}

// GetShaderParameter implements the OpenGL interface.
func (debugging *debuggingOpenGL) GetShaderParameter(shader uint32, param uint32) int32 {
	debugging.recordEntry("GetShaderParameter", shader, param)
	result := debugging.gl.GetShaderParameter(shader, param)
	debugging.recordExit("GetShaderParameter", result)
	return result
}

// GetUniformLocation implements the OpenGL interface.
func (debugging *debuggingOpenGL) GetUniformLocation(program uint32, name string) int32 {
	debugging.recordEntry("GetUniformLocation", program, name)
	result := debugging.gl.GetUniformLocation(program, name)
	debugging.recordExit("GetUniformLocation", result)
	return result
}

// IsEnabled implements the OpenGL interface.
func (debugging *debuggingOpenGL) IsEnabled(cap uint32) bool {
	debugging.recordEntry("IsEnabled", cap)
	result := debugging.gl.IsEnabled(cap)
	debugging.recordExit("IsEnabled", result)
	return result
}

// LinkProgram implements the OpenGL interface.
func (debugging *debuggingOpenGL) LinkProgram(program uint32) {
	debugging.recordEntry("LinkProgram", program)
	debugging.gl.LinkProgram(program)
	debugging.recordExit("LinkProgram")
}

// PixelStorei implements the OpenGL interface.
func (debugging *debuggingOpenGL) PixelStorei(name uint32, param int32) {
	debugging.recordEntry("PixelStorei", name, param)
	debugging.gl.PixelStorei(name, param)
	debugging.recordExit("PixelStorei")
}

// PolygonMode implements the OpenGL interface.
func (debugging *debuggingOpenGL) PolygonMode(face uint32, mode uint32) {
	debugging.recordEntry("PolygonMode", face, mode)
	debugging.gl.PolygonMode(face, mode)
	debugging.recordExit("PolygonMode")
}

// ReadPixels implements the OpenGL interface.
func (debugging *debuggingOpenGL) ReadPixels(x int32, y int32, width int32, height int32, format uint32, pixelType uint32, pixels interface{}) {
	debugging.recordEntry("ReadPixels", x, y, width, height, format, pixelType, pixels)
	debugging.gl.ReadPixels(x, y, width, height, format, pixelType, pixels)
	debugging.recordExit("ReadPixels")
}

// RenderbufferStorage implements the opengl.OpenGL interface.
func (debugging *debuggingOpenGL) RenderbufferStorage(target uint32, internalFormat uint32, width int32, height int32) {
	debugging.recordEntry("RenderbufferStorage", target, internalFormat, width, height)
	debugging.gl.RenderbufferStorage(target, internalFormat, width, height)
	debugging.recordExit("RenderbufferStorage")
}

// Scissor implements the OpenGL interface.
func (debugging *debuggingOpenGL) Scissor(x, y int32, width, height int32) {
	debugging.recordEntry("Scissor", x, y, width, height)
	debugging.gl.Scissor(x, y, width, height)
	debugging.recordExit("Scissor")
}

// ShaderSource implements the OpenGL interface.
func (debugging *debuggingOpenGL) ShaderSource(shader uint32, source string) {
	debugging.recordEntry("ShaderSource", shader, source)
	debugging.gl.ShaderSource(shader, source)
	debugging.recordExit("ShaderSource")
}

// TexImage2D implements the opengl.OpenGL interface.
func (debugging *debuggingOpenGL) TexImage2D(target uint32, level int32, internalFormat uint32, width int32, height int32,
	border int32, format uint32, xtype uint32, pixels interface{}) {
	debugging.recordEntry("TexImage2D", target, level, internalFormat, width, height, border, format, xtype, pixels)
	debugging.gl.TexImage2D(target, level, internalFormat, width, height, border, format, xtype, pixels)
	debugging.recordExit("TexImage2D")
}

// TexParameteri implements the opengl.OpenGL interface.
func (debugging *debuggingOpenGL) TexParameteri(target uint32, pname uint32, param int32) {
	debugging.recordEntry("TexParameteri", target, pname, param)
	debugging.gl.TexParameteri(target, pname, param)
	debugging.recordExit("TexParameteri")
}

// Uniform1i implements the opengl.OpenGL interface.
func (debugging *debuggingOpenGL) Uniform1i(location int32, value int32) {
	debugging.recordEntry("Uniform1i", location, value)
	debugging.gl.Uniform1i(location, value)
	debugging.recordExit("Uniform1i")
}

// Uniform4fv implements the opengl.OpenGL interface.
func (debugging *debuggingOpenGL) Uniform4fv(location int32, value *[4]float32) {
	debugging.recordEntry("Uniform4fv", location, value)
	debugging.gl.Uniform4fv(location, value)
	debugging.recordExit("Uniform4fv")
}

// UniformMatrix4fv implements the OpenGL interface.
func (debugging *debuggingOpenGL) UniformMatrix4fv(location int32, transpose bool, value *[16]float32) {
	debugging.recordEntry("UniformMatrix4fv", location, transpose, value)
	debugging.gl.UniformMatrix4fv(location, transpose, value)
	debugging.recordExit("UniformMatrix4fv")
}

// UseProgram implements the OpenGL interface.
func (debugging *debuggingOpenGL) UseProgram(program uint32) {
	debugging.recordEntry("UseProgram", program)
	debugging.gl.UseProgram(program)
	debugging.recordExit("UseProgram")
}

// VertexAttribOffset implements the OpenGL interface.
func (debugging *debuggingOpenGL) VertexAttribOffset(index uint32, size int32, attribType uint32, normalized bool, stride int32, offset int) {
	debugging.recordEntry("VertexAttribOffset", index, size, attribType, normalized, stride, offset)
	debugging.gl.VertexAttribOffset(index, size, attribType, normalized, stride, offset)
	debugging.recordExit("VertexAttribOffset")
}

// Viewport implements the OpenGL interface.
func (debugging *debuggingOpenGL) Viewport(x int32, y int32, width int32, height int32) {
	debugging.recordEntry("Viewport", x, y, width, height)
	debugging.gl.Viewport(x, y, width, height)
	debugging.recordExit("Viewport")
}
