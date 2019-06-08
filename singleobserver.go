package rxgo

import (
	"github.com/reactivex/rxgo/handlers"
)

// SingleObserver represents a group of EventHandlers.
type SingleObserver interface {
	handlers.EventHandler
	Disposable

	Block() (interface{}, error)
	OnSuccess(item interface{})
	OnError(err error)
}

type singleObserver struct {
	disposed    chan struct{}
	nextHandler handlers.NextFunc
	errHandler  handlers.ErrFunc
	done        chan interface{}
}

// NewSingleObserver constructs a new SingleObserver instance with default SingleObserver and accept
// any number of EventHandler
func NewSingleObserver(eventHandlers ...handlers.EventHandler) SingleObserver {
	ob := singleObserver{
		disposed: make(chan struct{}),
	}

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
	close(o.disposed)
}

func (o *singleObserver) IsDisposed() bool {
	select {
	case <-o.disposed:
		return true
	default:
		return false
	}
}

// Block terminates the SingleObserver's internal Observable
func (o *singleObserver) Block() (interface{}, error) {
	if !o.IsDisposed() {
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

// OnNext applies SingleObserver's NextHandler to an Item
func (o *singleObserver) OnSuccess(item interface{}) {
	if !o.IsDisposed() {
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
	if !o.IsDisposed() {
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
