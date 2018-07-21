package external

import "github.com/inkyblackness/hacked/ui/gui"

// Import starts an import dialog series, calling the given callback with a file name.
func Import(machine gui.ModalStateMachine, info string, callback func(string), lastFailed bool) {
	machine.SetState(&importStartState{
		machine:   machine,
		callback:  callback,
		info:      info,
		withError: lastFailed,
	})
}
