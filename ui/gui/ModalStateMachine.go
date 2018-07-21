package gui

// ModalStateMachine allows to set a new state.
type ModalStateMachine interface {
	SetState(state ModalState)
}
