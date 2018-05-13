package opengl

import "fmt"

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
