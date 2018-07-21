package external

import "github.com/inkyblackness/hacked/ui/gui"

// Export starts an export dialog series, calling the given callback with a folder name.
func Export(machine gui.ModalStateMachine, info string, callback func(string), lastFailed bool) {
	machine.SetState(&exportStartState{
		machine:   machine,
		callback:  callback,
		info:      info,
		withError: lastFailed,
	})
}
