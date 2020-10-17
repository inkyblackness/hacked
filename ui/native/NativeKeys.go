package native

import (
	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/inkyblackness/hacked/ui/input"
)

var keyMap = map[glfw.Key]input.Key{
	glfw.KeyEnter:     input.KeyEnter,
	glfw.KeyKPEnter:   input.KeyEnter,
	glfw.KeyEscape:    input.KeyEscape,
	glfw.KeyBackspace: input.KeyBackspace,
	glfw.KeyTab:       input.KeyTab,

	glfw.KeyDown:  input.KeyDown,
	glfw.KeyLeft:  input.KeyLeft,
	glfw.KeyRight: input.KeyRight,
	glfw.KeyUp:    input.KeyUp,

	glfw.KeyDelete:   input.KeyDelete,
	glfw.KeyEnd:      input.KeyEnd,
	glfw.KeyHome:     input.KeyHome,
	glfw.KeyInsert:   input.KeyInsert,
	glfw.KeyPageDown: input.KeyPageDown,
	glfw.KeyPageUp:   input.KeyPageUp,

	glfw.KeyLeftAlt:      input.KeyAlt,
	glfw.KeyLeftControl:  input.KeyControl,
	glfw.KeyLeftShift:    input.KeyShift,
	glfw.KeyLeftSuper:    input.KeySuper,
	glfw.KeyRightAlt:     input.KeyAlt,
	glfw.KeyRightControl: input.KeyControl,
	glfw.KeyRightShift:   input.KeyShift,
	glfw.KeyRightSuper:   input.KeySuper,

	glfw.KeyPause:       input.KeyPause,
	glfw.KeyPrintScreen: input.KeyPrintScreen,

	glfw.KeyCapsLock:   input.KeyCapsLock,
	glfw.KeyScrollLock: input.KeyScrollLock,

	glfw.KeyF1:  input.KeyF1,
	glfw.KeyF10: input.KeyF10,
	glfw.KeyF11: input.KeyF11,
	glfw.KeyF12: input.KeyF12,
	glfw.KeyF13: input.KeyF13,
	glfw.KeyF14: input.KeyF14,
	glfw.KeyF15: input.KeyF15,
	glfw.KeyF16: input.KeyF16,
	glfw.KeyF17: input.KeyF17,
	glfw.KeyF18: input.KeyF18,
	glfw.KeyF19: input.KeyF19,
	glfw.KeyF2:  input.KeyF2,
	glfw.KeyF20: input.KeyF20,
	glfw.KeyF21: input.KeyF21,
	glfw.KeyF22: input.KeyF22,
	glfw.KeyF23: input.KeyF23,
	glfw.KeyF24: input.KeyF24,
	glfw.KeyF25: input.KeyF25,
	glfw.KeyF3:  input.KeyF3,
	glfw.KeyF4:  input.KeyF4,
	glfw.KeyF5:  input.KeyF5,
	glfw.KeyF6:  input.KeyF6,
	glfw.KeyF7:  input.KeyF7,
	glfw.KeyF8:  input.KeyF8,
	glfw.KeyF9:  input.KeyF9,
}
