package native

import (
	"runtime"
	"time"

	"github.com/inkyblackness/hacked/ui/opengl"
)

// Run creates a native OpenGL window, initializes it with the given function and
// then runs the event loop until the window shall be closed.
func Run(initializer func(opengl.Window) error, title string, framesPerSecond float64) error {
	runtime.LockOSThread()

	window, err := NewOpenGLWindow(title, framesPerSecond)
	if err != nil {
		return err
	}
	defer window.Close()

	err = initializer(window)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(time.Millisecond)
	for !window.ShouldClose() {
		<-ticker.C
		window.Update()
	}
	ticker.Stop()

	return nil
}
