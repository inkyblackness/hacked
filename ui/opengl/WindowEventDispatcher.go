package opengl

import (
	"github.com/inkyblackness/hacked/ui/input"
)

type keyDeferrer struct {
	window *WindowEventDispatcher
}

func (def *keyDeferrer) Key(key input.Key, modifier input.Modifier) {
	def.window.CallKey(key, modifier)
}

func (def *keyDeferrer) Modifier(modifier input.Modifier) {
	def.window.CallModifier(modifier)
}

// WindowEventDispatcher implements the common, basic functionality of WindowEventDispatcher.
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
	CallKey               KeyCallback
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
		CallKey:               func(input.Key, input.Modifier) {},
		CallModifier:          func(input.Modifier) {},
		CallCharCallback:      func(rune) {},
		CallFileDropCallback:  func([]string) {},
	}
}

// StickyKeyListener returns an instance of a listener acting as an adapter
// for the key-down/-up callbacks.
func (window *WindowEventDispatcher) StickyKeyListener() input.StickyKeyListener {
	return &keyDeferrer{window}
}

// OnClosing implements the WindowEventDispatcher interface.
func (window *WindowEventDispatcher) OnClosing(callback ClosingCallback) {
	window.CallClosing = callback
}

// OnClosed implements the WindowEventDispatcher interface.
func (window *WindowEventDispatcher) OnClosed(callback ClosedCallback) {
	window.CallClosed = callback
}

// OnRender implements the WindowEventDispatcher interface.
func (window *WindowEventDispatcher) OnRender(callback RenderCallback) {
	window.CallRender = callback
}

// OnResize implements the WindowEventDispatcher interface.
func (window *WindowEventDispatcher) OnResize(callback ResizeCallback) {
	window.CallResize = callback
}

// OnMouseMove implements the WindowEventDispatcher interface.
func (window *WindowEventDispatcher) OnMouseMove(callback MouseMoveCallback) {
	window.CallOnMouseMove = callback
}

// OnMouseButtonDown implements the WindowEventDispatcher interface.
func (window *WindowEventDispatcher) OnMouseButtonDown(callback MouseButtonCallback) {
	window.CallOnMouseButtonDown = callback
}

// OnMouseButtonUp implements the WindowEventDispatcher interface.
func (window *WindowEventDispatcher) OnMouseButtonUp(callback MouseButtonCallback) {
	window.CallOnMouseButtonUp = callback
}

// OnMouseScroll implements the WindowEventDispatcher interface.
func (window *WindowEventDispatcher) OnMouseScroll(callback MouseScrollCallback) {
	window.CallOnMouseScroll = callback
}

// OnKey implements the WindowEventDispatcher interface
func (window *WindowEventDispatcher) OnKey(callback KeyCallback) {
	window.CallKey = callback
}

// OnModifier implements the WindowEventDispatcher interface
func (window *WindowEventDispatcher) OnModifier(callback ModifierCallback) {
	window.CallModifier = callback
}

// OnCharCallback implements the WindowEventDispatcher interface
func (window *WindowEventDispatcher) OnCharCallback(callback CharCallback) {
	window.CallCharCallback = callback
}

// OnFileDropCallback implements the WindowEventDispatcher interface
func (window *WindowEventDispatcher) OnFileDropCallback(callback FileDropCallback) {
	window.CallFileDropCallback = callback
}
