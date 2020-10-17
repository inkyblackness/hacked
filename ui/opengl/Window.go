package opengl

import (
	"encoding/json"
	"io/ioutil"

	"github.com/inkyblackness/hacked/ui/input"
)

// ClosingCallback is the function to handle close requests by the user.
type ClosingCallback func()

// ClosedCallback is the function to clean up resources when the window is being closed.
type ClosedCallback func()

// RenderCallback is the function to receive render events. When the callback
// returns, the window will swap the internal buffer.
type RenderCallback func()

// MouseMoveCallback is the function to receive the current mouse coordinate while moving.
// Movement is reported while the cursor is within the client area of the window, and
// beyond the window as long as at least one captured button is pressed.
// Reported values are with sub-pixel precision, if possible.
type MouseMoveCallback func(x float32, y float32)

// MouseButtonCallback is the function to receive button up/down events.
// An Up event is sent for every reported Down event, even if the mouse cursor is outside
// the client area.
type MouseButtonCallback func(buttonMask uint32, modifier input.Modifier)

// MouseScrollCallback is the function to receive scroll events.
// Delta values are right-hand oriented: positive values go right/down/far.
type MouseScrollCallback func(dx float32, dy float32)

// ResizeCallback is called for a change of window dimensions.
type ResizeCallback func(width int, height int)

// CharCallback is called for typing a character.
type CharCallback func(char rune)

// KeyCallback is called for pressing or releasing a key on the keyboard.
type KeyCallback func(key input.Key, modifier input.Modifier)

// ModifierCallback is called when the currently active modifier changed.
type ModifierCallback func(modifier input.Modifier)

// FileDropCallback is called when one or more files were dropped into the window.
type FileDropCallback func(filePaths []string)

// WindowState describes how the window is currently presented to the user.
type WindowState struct {
	// Maximized is set to true if the window is currently resized to fill the screen.
	Maximized bool
	// Top marks the screen position of the upper border of the window.
	Top int
	// Left marks the screen position of the left border of the window.
	Left int
	// Width marks the width of the non-maximized window.
	Width int
	// Height marks the height of the non-maximized window.
	Height int
}

// WindowStateFromFile attempts to load the state from a file with given name.
// If this does not work, the given default state will be returned.
func WindowStateFromFile(filename string, defaultState WindowState) WindowState {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return defaultState
	}
	var state WindowState
	err = json.Unmarshal(data, &state)
	if err != nil {
		return defaultState
	}
	return state
}

// SaveTo stores the state in a file with given filename.
func (state WindowState) SaveTo(filename string) error {
	bytes, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, bytes, 0640)
}

// Window represents an OpenGL render surface.
type Window interface {
	// StateSnapshot returns the current state of the window.
	StateSnapshot() WindowState
	// RestoreState puts the window into the same state when StateSnapshot() was called.
	RestoreState(state WindowState)

	// ClipboardString returns the current value of the clipboard, if it is compatible with UTF-8.
	ClipboardString() (string, error)
	// SetClipboardString sets the current value of the clipboard as UTF-8 string.
	SetClipboardString(value string)

	// OnClosing registers a callback function which shall be called when the user requests to close the window.
	OnClosing(callback ClosingCallback)
	// OnClosed registers a callback function which shall be called when the window is being closed.
	OnClosed(callback ClosedCallback)
	// SetCloseRequest notifies the window system whether the window shall be closed.
	// It can be called during the callback of OnClosing() to abort a close request by the user.
	SetCloseRequest(shouldClose bool)

	// OpenGL returns the OpenGL API wrapper for this window.
	OpenGL() OpenGL
	// OnRender registers a callback function which shall be called to update the scene.
	OnRender(callback RenderCallback)

	// OnResize registers a callback function for sizing events.
	OnResize(callback ResizeCallback)
	// Size returns the dimensions of the window display area in pixel.
	Size() (width int, height int)
	// SetFullScreen sets the full screen state of the window.
	SetFullScreen(on bool)
	// SetProjectModified sets an indicator in the window frame that the project has not been saved.
	SetProjectModified(modified bool)

	// SetCursorVisible controls whether the mouse cursor is currently visible.
	SetCursorVisible(visible bool)

	// OnMouseMove registers a callback function for mouse move events.
	OnMouseMove(callback MouseMoveCallback)
	// OnMouseButtonDown registers a callback function for mouse button down events.
	OnMouseButtonDown(callback MouseButtonCallback)
	// OnMouseButtonUp registers a callback function for mouse button up events.
	OnMouseButtonUp(callback MouseButtonCallback)
	// OnMouseScroll registers a callback function for mouse scroll events.
	OnMouseScroll(callback MouseScrollCallback)

	// OnKey registers a callback function for key events.
	OnKey(callback KeyCallback)
	// OnModifier registers a callback function for change of modifier events.
	OnModifier(callback ModifierCallback)
	// OnCharCallback registers a callback function for typed characters.
	OnCharCallback(callback CharCallback)

	// OnFileDropCallback registers a callback function for dropped files.
	OnFileDropCallback(callback FileDropCallback)
}
