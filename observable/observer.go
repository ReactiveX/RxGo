package observable

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

type observer struct {
	nextHandler handlers.NextFunc
	errHandler  handlers.ErrFunc
	doneHandler handlers.DoneFunc
}

// NewObserver constructs a new Observer instance with default Observer and accept
// any number of EventHandler
func NewObserver(eventHandlers ...rx.EventHandler) Observer {
	ob := observer{}

	if len(eventHandlers) > 0 {
		for _, handler := range eventHandlers {
			switch handler := handler.(type) {
			case handlers.NextFunc:
				ob.nextHandler = handler
			case handlers.ErrFunc:
				ob.errHandler = handler
			case handlers.DoneFunc:
				ob.doneHandler = handler
			case *observer:
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
func (o *observer) Handle(item interface{}) {
	switch item := item.(type) {
	case error:
		o.errHandler(item)
		return
	default:
		o.nextHandler(item)
	}
}

// OnNext applies Observer's NextHandler to an Item
func (o *observer) OnNext(item interface{}) {
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
func (o *observer) OnError(err error) {
	if o.errHandler != nil {
		o.errHandler(err)
	}
}

// OnDone terminates the Observer's internal Observable
func (o *observer) OnDone() {
	if o.doneHandler != nil {
		o.doneHandler()
	}
}
