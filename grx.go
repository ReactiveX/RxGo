package grx

// Observer is a "sentinel" object consisting of three methods to handle event stream.
type Observer struct {
	OnNext      func(e Event)
	OnError     func(e Event)
	OnCompleted func(e Event)
}

// Event either emits a value, an error, or notify as completed.
type Event struct {
	Value     interface{}
	Error     error
	Completed bool
}

// Observable is a stream of events
// TODO: Consider removing Name
type Observable struct {
	Name      string
	Stream    chan Event
	Observer  *Observer
}

func (o *Observable) isCompleted() bool {
	if len(o.Stream) > 0 {
		return false
	}
	return true
}

// NewObservable constructs an empty Observable
func NewObservable(name string) *Observable {
	return &Observable{
		Name: name,
		Stream: make(chan Event),
	}
}


// Add adds an Event to the Observable
func (o *Observable) Add(ev Event) {
	go func() {
		o.Stream <- ev
	}()
}

// Just creates an observable with only one item and emit "as-is".
func Just(item interface{}) *Observable {
	o := &Observable{
		Stream: make(chan Event, 1),
	}

	go func() {
		o.Stream <- Event{Value: item}
		close(o.Stream)
	}()
	
	return o
}

// From creates an observable from a slice of items and emit them in order.
func From(items []interface{}) *Observable {
	o := &Observable{
		Stream: make(chan Event, len(items)),
	}
	go func() {
		for _, item := range items {
			o.Stream <- Event{Value: item}
		}
		close(o.Stream)
	}()

	return o
}

// Start creates an Observable from a directive-like function's returned event.
func Start(fx func() Event) *Observable {
	o := &Observable{ Stream: make(chan Event, 1) }
	go func() {
		o.Stream <- fx()
		close(o.Stream)
	}()
	
	return o
}

// Subscribe subscribes an Observer to the Observable and starts it.
func (o *Observable) Subscribe(ob *Observer) {
	o.Observer = ob

	// Loop over the Observable's stream.
	for ev := range o.Stream {
		if o.isCompleted() {
			o.Observer.OnNext(ev)
			o.Observer.OnCompleted(ev)
		} else {
			if ev.Error != nil && ev.Value == nil {
				o.Observer.OnError(ev)
			} else if ev.Value != nil && ev.Error == nil {
				o.Observer.OnNext(ev)
			}
		}
	}
}



