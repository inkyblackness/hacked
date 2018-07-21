package gui

// ModalStateWrapper is a wrapper over a series of modal popup states.
// It implements both ModalState and ModalStateMachine to allow nested states, if necessary.
type ModalStateWrapper struct {
	State ModalState
}

// SetState sets the new state.
func (machine *ModalStateWrapper) SetState(state ModalState) {
	machine.State = state
}

// Render renders the current state.
func (machine *ModalStateWrapper) Render() {
	if machine.State != nil {
		machine.State.Render()
	}
}

// HandleFiles forwards the given filenames to the current state.
func (machine *ModalStateWrapper) HandleFiles(filenames []string) {
	if machine.State != nil {
		machine.State.HandleFiles(filenames)
	}
}
