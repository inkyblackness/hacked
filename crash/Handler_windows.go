package crash

//#include <Windows.h>
import "C"
import (
	"fmt"
	"runtime/debug"
	"unsafe"
)

func init() {
	Handler = windowsHandler
}

func windowsHandler(versionInfo string) {
	if e := recover(); e != nil {
		text := `There was a crash in the application. This is bad and I am sorry.

Perhaps you can make something with the details of the error below.
If you can reproduce this, please make a screenshot of this box and report it with details on the http://www.systemshock.org forums.
Thank you!

`
		text += fmt.Sprintf("%s:\n%s", e, debug.Stack())
		messageBox(text, "Something unexpected happened - "+versionInfo)
	}
}

func messageBox(text, caption string) {
	textArg, textFin := wrapString(text)
	defer textFin()
	captionArg, captionFin := wrapString(caption)
	defer captionFin()
	var hwnd unsafe.Pointer
	C.MessageBoxExA((*C.struct_HWND__)(hwnd), textArg, captionArg, C.UINT(0), C.WORD(0))
}

func wrapString(value string) (wrapped *C.char, finisher func()) {
	wrapped = C.CString(value)
	finisher = func() { C.free(unsafe.Pointer(wrapped)) } // nolint: gas
	return
}
