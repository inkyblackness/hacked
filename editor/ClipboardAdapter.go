package editor

import "github.com/inkyblackness/hacked/ui/opengl"

type clipboardAdapter struct {
	window opengl.Window
}

func (adapter clipboardAdapter) String() (string, error) {
	return adapter.window.ClipboardString()
}

func (adapter clipboardAdapter) SetString(value string) {
	adapter.window.SetClipboardString(value)
}
