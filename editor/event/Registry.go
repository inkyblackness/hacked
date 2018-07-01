package event

import "reflect"

// Registry describes an entity where an event handler can be registered.
type Registry interface {
	// RegisterHandler must be called with a concrete structural type (implementing the Event interface),
	// and a function that takes one argument that is of the given type.
	// This function panics if this is not fulfilled.
	// The same handler function can be registered several times for the same type.
	// It will be called for each registration.
	//
	// The returned function can be used to unregister the handler again.
	RegisterHandler(eType reflect.Type, handlerFunc interface{}) func()
}
