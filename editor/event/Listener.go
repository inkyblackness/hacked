package event

// Listener receives events.
type Listener interface {
	Event(e Event)
}
