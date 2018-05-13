package opengl

// DebuggingEntryFunc is a function to be called on function entry of
type DebuggingEntryFunc func(name string, param ...interface{})

// DebuggingExitFunc is a function to be called after an OpenGL function has returned.
// result may be several parameters for functions with multiple return values.
type DebuggingExitFunc func(name string, result ...interface{})

// DebuggingErrorFunc is a function which is called if a previous OpenGL call caused an error state.
type DebuggingErrorFunc func(name string, errorCodes []uint32)

// DebugBuilder is a builder for an OpenGl implementation for debugging.
type DebugBuilder struct {
	wrapped OpenGl

	onEntry DebuggingEntryFunc
	onExit  DebuggingExitFunc
	onError DebuggingErrorFunc
}

// NewDebugBuilder wraps the provided OpenGl instance and returns a new builder
// instance.
func NewDebugBuilder(wrapped OpenGl) *DebugBuilder {
	builder := &DebugBuilder{
		wrapped: wrapped,
		onEntry: func(string, ...interface{}) {},
		onExit:  func(string, ...interface{}) {},
		onError: func(string, []uint32) {}}

	return builder
}

// Build creates a new instance of the debugging OpenGl implementation.
// The builder can be resused to create another instance with different parameters.
func (builder *DebugBuilder) Build() OpenGl {
	opengl := &debuggingOpenGl{
		gl: builder.wrapped,

		onEntry: builder.onEntry,
		onExit:  builder.onExit,
		onError: builder.onError}

	return opengl
}

// OnEntry registers a callback function to be called before an OpenGL function is called.
func (builder *DebugBuilder) OnEntry(callback DebuggingEntryFunc) *DebugBuilder {
	builder.onEntry = callback
	return builder
}

// OnExit registers a callback function to be called after an OpenGL function returned.
func (builder *DebugBuilder) OnExit(callback DebuggingExitFunc) *DebugBuilder {
	builder.onExit = callback
	return builder
}

// OnError registers a callback function to be called when an OpenGL function set an error state.
func (builder *DebugBuilder) OnError(callback DebuggingErrorFunc) *DebugBuilder {
	builder.onError = callback
	return builder
}
