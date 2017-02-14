package observer

import (
	"github.com/reactivex/rxgo"
	"github.com/reactivex/rxgo/handlers"
)

// Observer represents a group of EventHandlers.
type Observer struct {
	NextHandler handlers.NextFunc
	ErrHandler  handlers.ErrFunc
	DoneHandler handlers.DoneFunc
}

// DefaultObserver guarantees any handler won't be nil.
var DefaultObserver = Observer{
	NextHandler: func(interface{}) {},
	ErrHandler:  func(err error) {},
	DoneHandler: func() {},
}

// Handle registers Observer to EventHandler.
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
func New(eventHandlers ...rx.EventHandler) Observer {
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
