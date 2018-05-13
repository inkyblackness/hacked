package opengl

// AttributeSetter is a callback function for a VertexArrayObject.
type AttributeSetter func(OpenGl)

// VertexArrayObject registers common vertex properties which are pre-set
// before regular rendering loops.
// Depending on the supported level of the Open GL API, these settings are
// either buffered (using a VAO - vertex array object), or replayed for
// every render call.
type VertexArrayObject struct {
	gl      OpenGl
	program uint32
	handle  uint32
	setter  AttributeSetter
}

// NewVertexArrayObject returns a new instance.
func NewVertexArrayObject(gl OpenGl, program uint32) *VertexArrayObject {
	vao := &VertexArrayObject{
		gl:      gl,
		program: program,
		handle:  gl.GenVertexArrays(1)[0],
		setter:  func(OpenGl) {}}

	return vao
}

// Dispose drops any associated resources.
func (vao *VertexArrayObject) Dispose() {
	vao.gl.DeleteVertexArrays([]uint32{vao.handle})
	vao.handle = 0
}

// WithSetter registers the setter function for this object.
// Depending on the support of the OpenGL API, this setter may be called
// once immediately, or for each render call.
func (vao *VertexArrayObject) WithSetter(setter AttributeSetter) {
	if vao.handle != 0 {
		vao.OnShader(func() { setter(vao.gl) })
	}
	vao.setter = setter
}

// OnShader executes the provided task with an activated shader program
// and the vertex array object initialized.
func (vao *VertexArrayObject) OnShader(task func()) {
	gl := vao.gl

	gl.UseProgram(vao.program)
	gl.BindVertexArray(vao.handle)

	defer func() {
		gl.BindVertexArray(0)
		gl.UseProgram(0)
	}()

	if vao.handle == 0 {
		vao.setter(gl)
	}
	task()
}
