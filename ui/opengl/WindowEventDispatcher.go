package opengl

import (
	"github.com/inkyblackness/hacked/ui/input"
)

type keyDeferrer struct {
	window *WindowEventDispatcher
}

func (def *keyDeferrer) KeyPress(key input.Key, modifier input.Modifier) {
	def.window.CallKeyPress(key, modifier)
}

func (def *keyDeferrer) KeyRelease(key input.Key, modifier input.Modifier) {
	def.window.CallKeyRelease(key, modifier)
}

func (def *keyDeferrer) Modifier(modifier input.Modifier) {
	def.window.CallModifier(modifier)
}

// WindowEventDispatcher implements the common, basic functionality of Window.
type WindowEventDispatcher struct {
	CallClosing           ClosingCallback
	CallClosed            ClosedCallback
	CallRender            RenderCallback
	CallResize            ResizeCallback
	CallOnMouseMove       MouseMoveCallback
	CallOnMouseButtonUp   MouseButtonCallback
	CallOnMouseButtonDown MouseButtonCallback
	CallOnMouseScroll     MouseScrollCallback
	CallModifier          ModifierCallback
	CallKeyPress          KeyCallback
	CallKeyRelease        KeyCallback
	CallCharCallback      CharCallback
	CallFileDropCallback  FileDropCallback
}

// NullWindowEventDispatcher returns an initialized instance with empty callbacks.
func NullWindowEventDispatcher() WindowEventDispatcher {
	return WindowEventDispatcher{
		CallClosing:           func() {},
		CallClosed:            func() {},
		CallRender:            func() {},
		CallResize:            func(int, int) {},
		CallOnMouseMove:       func(float32, float32) {},
		CallOnMouseButtonUp:   func(uint32, input.Modifier) {},
		CallOnMouseButtonDown: func(uint32, input.Modifier) {},
		CallOnMouseScroll:     func(float32, float32) {},
		CallKeyPress:          func(input.Key, input.Modifier) {},
		CallKeyRelease:        func(input.Key, input.Modifier) {},
		CallModifier:          func(input.Modifier) {},
		CallCharCallback:      func(rune) {},
		CallFileDropCallback:  func([]string) {},
	}
}

// StickyKeyListener returns an instance of a listener acting as an adapter
// for the key-down/-up callbacks.
func (window *WindowEventDispatcher) StickyKeyListener() input.StickyKeyListener {
	return &keyDeferrer{window: window}
}

// OnClosing implements the Window interface.
func (window *WindowEventDispatcher) OnClosing(callback ClosingCallback) {
	window.CallClosing = callback
}

// OnClosed implements the Window interface.
func (window *WindowEventDispatcher) OnClosed(callback ClosedCallback) {
	window.CallClosed = callback
}

// OnRender implements the Window interface.
func (window *WindowEventDispatcher) OnRender(callback RenderCallback) {
	window.CallRender = callback
}

// OnResize implements the Window interface.
func (window *WindowEventDispatcher) OnResize(callback ResizeCallback) {
	window.CallResize = callback
}

// OnMouseMove implements the Window interface.
func (window *WindowEventDispatcher) OnMouseMove(callback MouseMoveCallback) {
	window.CallOnMouseMove = callback
}

// OnMouseButtonDown implements the Window interface.
func (window *WindowEventDispatcher) OnMouseButtonDown(callback MouseButtonCallback) {
	window.CallOnMouseButtonDown = callback
}

// OnMouseButtonUp implements the Window interface.
func (window *WindowEventDispatcher) OnMouseButtonUp(callback MouseButtonCallback) {
	window.CallOnMouseButtonUp = callback
}

// OnMouseScroll implements the Window interface.
func (window *WindowEventDispatcher) OnMouseScroll(callback MouseScrollCallback) {
	window.CallOnMouseScroll = callback
}

// OnKeyPress implements the Window interface
func (window *WindowEventDispatcher) OnKeyPress(callback KeyCallback) {
	window.CallKeyPress = callback
}

// OnKeyRelease implements the Window interface
func (window *WindowEventDispatcher) OnKeyRelease(callback KeyCallback) {
	window.CallKeyRelease = callback
}

// OnModifier implements the WindowEventDispatcher interface
func (window *WindowEventDispatcher) OnModifier(callback ModifierCallback) {
	window.CallModifier = callback
}

// OnCharCallback implements the Window interface
func (window *WindowEventDispatcher) OnCharCallback(callback CharCallback) {
	window.CallCharCallback = callback
}

// OnFileDropCallback implements the Window interface
func (window *WindowEventDispatcher) OnFileDropCallback(callback FileDropCallback) {
	window.CallFileDropCallback = callback
}
