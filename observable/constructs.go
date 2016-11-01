package observable

import (
	"github.com/jochasinga/grx/event"
	"github.com/jochasinga/grx/observer"
)

// Observable is a stream of events
type Observable struct {
	Stream    chan *event.Event
	Observer  *observer.Observer
}

// To query a channel's length, this method should block.
func (observable *Observable) isCompleted() bool {
	if len(observable.Stream) > 0 {
		return false
	}
	return true
}

// New constructs an empty Observable with 0 or more buffer length.
// myStream := observable.New(1)
func New(buf ...int) *Observable {
	bufferLen := 0
	if len(buf) > 0 {
		bufferLen = buf[len(buf)-1]
	}
	return &Observable{ Stream: make(chan *event.Event, bufferLen) }
}
