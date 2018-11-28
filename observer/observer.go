package observer

import (
	"github.com/reactivex/rxgo"
	"github.com/reactivex/rxgo/handlers"
)

// Observer represents a group of EventHandlers.
type Observer interface {
	rx.EventHandler

	OnNext(item interface{})
	OnError(err error)
	OnDone()
}

type observator struct {
	nextHandler handlers.NextFunc
	errHandler  handlers.ErrFunc
	doneHandler handlers.DoneFunc
}

// New constructs a new Observer instance with default Observer and accept
// any number of EventHandler
func New(eventHandlers ...rx.EventHandler) Observer {
	ob := observator{}

	if len(eventHandlers) > 0 {
		for _, handler := range eventHandlers {
			switch handler := handler.(type) {
			case handlers.NextFunc:
				ob.nextHandler = handler
			case handlers.ErrFunc:
				ob.errHandler = handler
			case handlers.DoneFunc:
				ob.doneHandler = handler
			case *observator:
				ob = *handler
			}
		}
	}

	if ob.nextHandler == nil {
		ob.nextHandler = func(interface{}) {}
	}
	if ob.errHandler == nil {
		ob.errHandler = func(err error) {}
	}
	if ob.doneHandler == nil {
		ob.doneHandler = func() {}
	}

	return &ob
}

// Handle registers Observer to EventHandler.
func (o *observator) Handle(item interface{}) {
	switch item := item.(type) {
	case error:
		o.errHandler(item)
		return
	default:
		o.nextHandler(item)
	}
}

// OnNext applies Observer's NextHandler to an Item
func (o *observator) OnNext(item interface{}) {
	switch item := item.(type) {
	case error:
		return
	default:
		if o.nextHandler != nil {
			o.nextHandler(item)
		}
	}
}

// OnError applies Observer's ErrHandler to an error
func (o *observator) OnError(err error) {
	if o.errHandler != nil {
		o.errHandler(err)
	}
}

// OnDone terminates the Observer's internal Observable
func (o *observator) OnDone() {
	if o.doneHandler != nil {
		o.doneHandler()
	}
}
