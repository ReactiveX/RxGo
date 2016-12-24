package eventstream

import (
	"github.com/jochasinga/grx/bases"
)

type EventStream chan bases.Emitter

func (evs EventStream) isDone() bool {
	_, ok := <-evs
	return ok
}

// Next returns the next Event on the EventStream
func (evs EventStream) Next() bases.Emitter {
	if !evs.isDone() {
		return <-o.EventStream
	}
	return nil, NewError(grx.EndOfIteratorError)
}

// HasNext returns true if there is the next Event and false otherwise
func (evs EventStream) HasNext() bool {
	return !evs.isDone()
}

// New creates a new EventStream from one or more Event
func New(emitters ...Emitter) EventStream {
	es := make(EventStream)
	go func() {
		for _, emitter := range emitters {
			es <- emitter
		}
		close(es)
	}()
	return es
}

// From creates a new EventStream from an Iterator
func From(iter Iterator) EventStream {
	es := make(EventStream)
	go func() {
		for iter.HasNext() {
			es <- iter.Next()
		}
		close(es)
	}()
	return es
}
