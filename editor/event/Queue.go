package event

// Queue stores events to be dispatched at a later time.
type Queue struct {
	events []Event
}

// IsEmpty returns true if there are no more events pending.
func (queue Queue) IsEmpty() bool {
	return len(queue.events) == 0
}

// Event queues the given event at the tail.
// When events are dispatched, they are so first-in, first-out (FIFO).
// nil events are ignored. Typed event pointer that are nil are not ignored.
func (queue *Queue) Event(e Event) {
	if e == nil {
		return
	}
	queue.events = append(queue.events, e)
}

// ForwardTo forwards all currently queued events to the given listener.
// If there are new events added to the queue during this call, they will be put on hold,
// to be dispatched during a next call.
func (queue *Queue) ForwardTo(listener Listener) {
	eventsToDispatch := queue.events
	queue.events = make([]Event, 0, cap(queue.events))
	for _, e := range eventsToDispatch {
		listener.Event(e)
	}
}
