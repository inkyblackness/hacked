package opengl

import (
	"fmt"
)

var errorStrings = map[uint32]string{
	NO_ERROR:                      "NO_ERROR",
	INVALID_ENUM:                  "INVALID_ENUM",
	INVALID_VALUE:                 "INVALID_VALUE",
	INVALID_OPERATION:             "INVALID_OPERATION",
	STACK_OVERFLOW:                "STACK_OVERFLOW",
	STACK_UNDERFLOW:               "STACK_UNDERFLOW",
	OUT_OF_MEMORY:                 "OUT_OF_MEMORY",
	INVALID_FRAMEBUFFER_OPERATION: "INVALID_FRAMEBUFFER_OPERATION"}

// ErrorString returns a readable version of the provided error code. If the
// code is unknown, the textual representation of the hexadecimal value is
// returned.
func ErrorString(errorCode uint32) string {
	value, exists := errorStrings[errorCode]

	if !exists {
		value = fmt.Sprintf("0x%04X", errorCode)
	}

	return value
}

// ShaderError describes a problem with a log.
type ShaderError struct {
	Log string
}

// Error returns the log.
func (err ShaderError) Error() string {
	return err.Log
}

// NamedShaderError is an error with a name.
type NamedShaderError struct {
	Name   string
	Nested error
}

// Error returns the name with the nested error information.
func (err NamedShaderError) Error() string {
	return fmt.Sprintf("%s failed: %v", err.Name, err.Nested)
}

// Unwrap returns the nested error.
func (err NamedShaderError) Unwrap() error {
	return err.Nested
}
