package event

import (
	"reflect"
)

type pendingHandlerAction struct {
	add bool
	val reflect.Value
}

type handlerEntry struct {
	dispatching bool
	list        []reflect.Value

	pending []pendingHandlerAction
}

// Dispatcher is a distributor of events to registered handlers.
// A handler is a event-type-specific function that takes the concrete type of an event as a parameter.
type Dispatcher struct {
	handlers map[reflect.Type]*handlerEntry
}

// NewDispatcher returns a new instance.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[reflect.Type]*handlerEntry),
	}
}

// RegisterHandler must be called with a concrete structural type (implementing the Event interface),
// and a function that takes one argument that is of the given type.
// This function panics if this is not fulfilled.
// The same handler function can be registered several times for the same type.
// It will be called for each registration.
//
// The returned function can be used to unregister the handler again. It is a closure over
// UnregisterHandler(eType, handlerFunc).
func (dispatcher *Dispatcher) RegisterHandler(eType reflect.Type, handlerFunc interface{}) func() {
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

	entry, existing := dispatcher.handlers[eType]
	if !existing {
		entry = &handlerEntry{}
		dispatcher.handlers[eType] = entry
	}
	handlerValue := reflect.ValueOf(handlerFunc)
	if !entry.dispatching {
		entry.list = append(entry.list, handlerValue)
	} else {
		entry.pending = append(entry.pending, pendingHandlerAction{add: true, val: handlerValue})
	}

	return func() { dispatcher.UnregisterHandler(eType, handlerFunc) }
}

// UnregisterHandler removes a handler that was previously registered.
// If there was no registration done, this call is ignored.
// If the same handler was registered multiple times, all registrations are removed.
func (dispatcher *Dispatcher) UnregisterHandler(eType reflect.Type, handlerFunc interface{}) {
	handlerValue := reflect.ValueOf(handlerFunc)
	entry := dispatcher.handlers[eType]
	if !entry.dispatching {
		dispatcher.removeHandlerFromList(entry, handlerValue)
	} else {
		entry.pending = append(entry.pending, pendingHandlerAction{add: false, val: handlerValue})
	}
}

// Event dispatches the given event to all currently registered handlers.
func (dispatcher *Dispatcher) Event(e Event) {
	entry, existing := dispatcher.handlers[reflect.TypeOf(e)]
	if !existing {
		return
	}
	if entry.dispatching {
		panic("event distribution during event of same type not supported.")
	}
	args := []reflect.Value{reflect.ValueOf(e)}
	entry.dispatching = true
	for _, handler := range entry.list {
		if dispatcher.isHandlerStillRegistered(entry, handler) {
			handler.Call(args)
		}
	}
	for _, pending := range entry.pending {
		if pending.add {
			entry.list = append(entry.list, pending.val)
		} else {
			dispatcher.removeHandlerFromList(entry, pending.val)
		}
	}
	entry.pending = nil
	entry.dispatching = false
}

func (dispatcher *Dispatcher) removeHandlerFromList(entry *handlerEntry, handlerValue reflect.Value) {
	listLen := len(entry.list)
	removed := 0
	for i := listLen - 1; i >= 0; i-- {
		if entry.list[i] == handlerValue {
			removed++
			copy(entry.list[i:listLen-removed], entry.list[i+1:listLen-removed+1])
		}
	}
	if removed > 0 {
		entry.list = entry.list[:listLen-removed]
	}
}

func (dispatcher *Dispatcher) isHandlerStillRegistered(entry *handlerEntry, handler reflect.Value) bool {
	for _, pending := range entry.pending {
		if !pending.add && (pending.val == handler) {
			return false
		}
	}
	return true
}
