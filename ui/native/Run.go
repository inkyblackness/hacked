package native

import (
	"runtime"
	"time"

	"github.com/inkyblackness/hacked/ui/opengl"
)

// Run creates a native OpenGL window, initializes it with the given function and
// then runs the event loop until the window shall be closed.
// The provided deferrer is a channel of tasks that can be injected into the event loop.
func Run(initializer func(opengl.Window), title string, framesPerSecond float64, deferrer <-chan func()) (err error) {
	runtime.LockOSThread()

	var window *OpenGlWindow
	window, err = NewOpenGlWindow(title, framesPerSecond)
	if err != nil {
		return
	}

	initializer(window)
	for !window.ShouldClose() {
		select {
		case task := <-deferrer:
			task()
		default:
			time.Sleep(time.Nanosecond)
		}
		window.Update()
	}
	window.Close()

	return
}
