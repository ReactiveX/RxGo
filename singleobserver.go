package rx

import (
	"github.com/reactivex/rxgo/handlers"
)

// SingleObserver represents a group of EventHandlers.
type SingleObserver interface {
	EventHandler
	Disposable

	OnSuccess(item interface{})
	OnError(err error)

	Block() (interface{}, error)
}

type singleObserver struct {
	disposed    bool
	nextHandler handlers.NextFunc
	errHandler  handlers.ErrFunc
	done        chan interface{}
}

// NewSinglesingleObserver constructs a new SingleObserver instance with default SingleObserver and accept
// any number of EventHandler
func NewSingleObserver(eventHandlers ...EventHandler) SingleObserver {
	ob := singleObserver{}

	if len(eventHandlers) > 0 {
		for _, handler := range eventHandlers {
			switch handler := handler.(type) {
			case handlers.NextFunc:
				ob.nextHandler = handler
			case handlers.ErrFunc:
				ob.errHandler = handler
			case *singleObserver:
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
	ob.done = make(chan interface{}, 1)

	return &ob
}

// Handle registers SingleObserver to EventHandler.
func (o *singleObserver) Handle(item interface{}) {
	switch item := item.(type) {
	case error:
		o.errHandler(item)
		return
	default:
		o.nextHandler(item)
	}
}

func (o *singleObserver) Dispose() {
	o.disposed = true
}

func (o *singleObserver) IsDisposed() bool {
	return o.disposed
}

// OnNext applies SingleObserver's NextHandler to an Item
func (o *singleObserver) OnSuccess(item interface{}) {
	if !o.disposed {
		switch item := item.(type) {
		case error:
			return
		default:
			if o.nextHandler != nil {
				o.nextHandler(item)
				o.Dispose()
				if o.done != nil {
					o.done <- item
					close(o.done)
				}
			}
		}
	} else {
		// TODO
	}
}

// OnError applies SingleObserver's ErrHandler to an error
func (o *singleObserver) OnError(err error) {
	if !o.disposed {
		if o.errHandler != nil {
			o.errHandler(err)
			o.Dispose()
			if o.done != nil {
				o.done <- err
				close(o.done)
			}
		}
	} else {
		// TODO
	}
}

// OnDone terminates the SingleObserver's internal Observable
func (o *singleObserver) Block() (interface{}, error) {
	if !o.disposed {
		for v := range o.done {
			switch v := v.(type) {
			case error:
				return nil, v
			default:
				return v, nil
			}
		}
	}
	return nil, nil
}
