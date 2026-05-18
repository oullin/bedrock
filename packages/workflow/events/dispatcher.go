package events

// Listener handles a workflow event.
type Listener[T any] func(Event[T])

// Dispatcher routes typed workflow events to listeners keyed by string event name.
type Dispatcher[T any] struct {
	listeners map[string][]Listener[T]
}

func NewDispatcher[T any]() *Dispatcher[T] {
	return &Dispatcher[T]{listeners: map[string][]Listener[T]{}}
}

// On registers a listener for the given event name.
func (d *Dispatcher[T]) On(name string, listener Listener[T]) {
	d.listeners[name] = append(d.listeners[name], listener)
}

// OnGuard is a typed convenience for guard listeners.
func (d *Dispatcher[T]) OnGuard(name string, listener func(*GuardEvent[T])) {
	d.On(name, func(event Event[T]) {
		guard, ok := event.(*GuardEvent[T])

		if ok {
			listener(guard)
		}
	})
}

// Dispatch fires the event to every listener registered for `name`.
func (d *Dispatcher[T]) Dispatch(name string, event Event[T]) {
	if d == nil {
		return
	}

	for _, listener := range d.listeners[name] {
		listener(event)
	}
}
