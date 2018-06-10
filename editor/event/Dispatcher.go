package event

import "reflect"

// Dispatcher is a distributor of events to registered handlers.
// A handler is a event-type-specific function that takes the concrete type of an event as a parameter.
type Dispatcher struct {
	handlers map[reflect.Type][]reflect.Value
}

// NewDispatcher returns a new instance.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[reflect.Type][]reflect.Value),
	}
}

// RegisterHandler must be called with a concrete structural type (implementing the Event interface),
// and a function that takes one argument that is of the given type.
// This function panics if this is not fulfilled.
func (dispatcher *Dispatcher) RegisterHandler(eType reflect.Type, handlerFunc interface{}) {
	if eType.Kind() != reflect.Struct {
		panic("event type must be a structure")
	}
	handlerType := reflect.TypeOf(handlerFunc)
	if handlerType.Kind() != reflect.Func {
		panic("handler must be a function")
	}
	if handlerType.NumIn() != 1 {
		panic("handler takes wrong number of arguments")
	}
	if handlerType.In(0) != eType {
		panic("handler does not take given type as argument")
	}

	handlerList := append(dispatcher.handlers[eType], reflect.ValueOf(handlerFunc))
	dispatcher.handlers[eType] = handlerList
}

// UnregisterHandler removes a handler that was previously registered.
// If there was no registration done, this call is ignored.
func (dispatcher *Dispatcher) UnregisterHandler(eType reflect.Type, handler interface{}) {
	handlerValue := reflect.ValueOf(handler)
	handlerList := dispatcher.handlers[eType]
	listLen := len(handlerList)
	removed := 0
	for i := listLen - 1; i >= 0; i-- {
		if handlerList[i] == handlerValue {
			removed++
			copy(handlerList[i:listLen-removed], handlerList[i+1:listLen-removed+1])
		}
	}
	if removed > 0 {
		dispatcher.handlers[eType] = handlerList[:listLen-removed]
	}
}

// Event dispatches the given event to all currently registered handlers.
func (dispatcher *Dispatcher) Event(e Event) {
	args := []reflect.Value{reflect.ValueOf(e)}
	eType := reflect.TypeOf(e)
	handlerList := dispatcher.handlers[eType]
	for _, handler := range handlerList {
		handler.Call(args)
	}
}
