package grx

// Observer is a "sentinel" object consisting of three methods to handle event stream.
type Observer struct {
	OnNext      func(e *Event)
	OnError     func(e *Event)
	OnCompleted func(e *Event)
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
	Stream    chan *Event
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
		Stream: make(chan *Event),
	}
}

// Add adds an Event to the Observable
func (o *Observable) Add(ev *Event) {
	go func() {
		o.Stream <- ev
	}()
}

func Empty() *Observable {
	o := &Observable{}
	return o
}

// Just creates an observable with only one item and emit "as-is".
func Just(item interface{}) *Observable {
	o := &Observable{ Stream: make(chan *Event, 1) }
	go func() {
		o.Stream <- &Event{Value: item}
		close(o.Stream)
	}()
	return o
}

// From creates an observable from a slice of items and emit them in order.
func From(items []interface{}) *Observable {
	o := &Observable{
		Stream: make(chan *Event, len(items)),
	}
	for _, item := range items {
		o.Stream <- &Event{Value: item}
	}
	close(o.Stream)
	return o
}

// Start creates an Observable from a directive-like function's returned event.
func Start(fx func() *Event) *Observable {
	o := &Observable{ Stream: make(chan *Event, 1) }
	go func() {
		o.Stream <- fx()
		close(o.Stream)
	}()
	
	return o
}

// Subscribe subscribes an Observer to the Observable and starts it.
func (o *Observable) Subscribe(ob *Observer) {
	o.Observer = ob

	if o.Stream != nil {
		// Loop over the Observable's stream.
		for ev := range o.Stream {

			// Check if the stream is completed 
			if o.isCompleted() {

				// Not sure if this is even necessary
				ev.Completed = true
				
				// Check if the cap is more than 0 (likely to be a last Event)
				if cap(o.Stream) > 0 {
					o.Observer.OnNext(ev)
					o.Observer.OnCompleted(ev)
					return
				}

				// If this is not a last Event (i.e. an empty Observable),
				// call OnCompleted once and return
				o.Observer.OnCompleted(ev)
				return
			} 
			if ev.Error != nil {
				o.Observer.OnError(ev)
				return
			}

			// Keep calling OnNext
			o.Observer.OnNext(ev)
		}
	}

	// A hack for empty, finite Observable--emit a "terminal" event to signal stream's termination.
	o.Observer.OnCompleted(&Event{ Completed: true })
	return
}



