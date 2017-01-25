package observer

import (
	"github.com/jochasinga/grx/bases"
	"github.com/jochasinga/grx/handlers"
)

// Observer represents
type Observer struct {
	sink        <-chan interface{}
	NextHandler handlers.NextFunc
	ErrHandler  handlers.ErrFunc
	DoneHandler handlers.DoneFunc
}

// DefaultObserver makes sure any handler won't turn up nil.
var DefaultObserver = Observer{
	NextHandler: func(interface{}) {},
	ErrHandler:  func(err error) {},
	DoneHandler: func() {},
}

// Handle makes Observer implements handlers.EventHandler
func (ob Observer) Handle(item interface{}) {
	switch item := item.(type) {
	case error:
		ob.ErrHandler(item)
		return
	default:
		ob.NextHandler(item)
	}
}

// New constructs a new Observer instance with default Observer and accept
// any number of EventHandler
func New(eventHandlers ...bases.EventHandler) Observer {
	ob := DefaultObserver
	if len(eventHandlers) > 0 {
		for _, handler := range eventHandlers {
			switch handler := handler.(type) {
			case handlers.NextFunc:
				ob.NextHandler = handler
			case handlers.ErrFunc:
				ob.ErrHandler = handler
			case handlers.DoneFunc:
				ob.DoneHandler = handler
			case Observer:
				ob = handler
			}
		}
	}
	return ob
}

// OnNext applies Observer's NextHandler to an Item
func (ob Observer) OnNext(item interface{}) {
	switch item := item.(type) {
	case error:
		return
	default:
		if ob.NextHandler != nil {
			ob.NextHandler(item)
		}
	}
}

// OnError applies Observer's ErrHandler to an error
func (ob Observer) OnError(err error) {
	if ob.ErrHandler != nil {
		ob.ErrHandler(err)
	}
}

// OnDone terminates the Observer's internal Observable
func (ob Observer) OnDone() {
	if ob.DoneHandler != nil {
		ob.DoneHandler()
	}
}
