package native

import (
	"runtime"
	"time"

	"github.com/inkyblackness/hacked/ui/opengl"
)

// Run creates a native OpenGL window, initializes it with the given function and
// then runs the event loop until the window shall be closed.
// The provided deferrer is a channel of tasks that can be injected into the event loop.
// When the channel is closed, the loop is stopped and the window is closed.
func Run(initializer func(opengl.Window) error, title string, framesPerSecond float64, deferrer <-chan func()) (err error) {
	runtime.LockOSThread()

	var window *OpenGLWindow
	window, err = NewOpenGLWindow(title, framesPerSecond)
	if err != nil {
		return
	}
	defer window.Close()

	err = initializer(window)
	if err != nil {
		return
	}

	stopLoop := false
	for !window.ShouldClose() && !stopLoop {
		select {
		case task, ok := <-deferrer:
			if ok {
				task()
			} else {
				stopLoop = true
			}
		case <-time.After(time.Millisecond):
		}
		window.Update()
	}

	return
}
