package event

// Registry describes an entity where an event handler can be registered.
type Registry interface {
	// RegisterHandler must be called a function that takes one argument which is
	// a concrete structural type (implementing the Event interface).
	// This function panics if this is not fulfilled.
	// The same handler function can be registered several times for the same type.
	// It will be called for each registration.
	//
	// The returned function can be used to unregister the handler again.
	RegisterHandler(handlerFunc interface{}) func()
}
