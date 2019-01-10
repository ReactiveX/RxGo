package rxgo

import (
	"sync"

	"github.com/reactivex/rxgo/v2/handlers"
)

// SingleObserver represents a group of EventHandlers.
type SingleObserver interface {
	handlers.EventHandler
	Disposable

	OnSuccess(item interface{})
	OnError(err error)

	Block() (interface{}, error)
}

type singleObserver struct {
	disposedMutex sync.Mutex
	disposed      bool
	nextHandler   handlers.NextFunc
	errHandler    handlers.ErrFunc
	done          chan interface{}
}

// NewSinglesingleObserver constructs a new SingleObserver instance with default SingleObserver and accept
// any number of EventHandler
func NewSingleObserver(eventHandlers ...handlers.EventHandler) SingleObserver {
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
	o.disposedMutex.Lock()
	o.disposed = true
	o.disposedMutex.Unlock()
}

func (o *singleObserver) IsDisposed() bool {
	o.disposedMutex.Lock()
	defer o.disposedMutex.Unlock()
	return o.disposed
}

// OnNext applies SingleObserver's NextHandler to an Item
func (o *singleObserver) OnSuccess(item interface{}) {
	o.disposedMutex.Lock()
	disposed := o.disposed
	o.disposedMutex.Unlock()
	if !disposed {
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
	o.disposedMutex.Lock()
	disposed := o.disposed
	o.disposedMutex.Unlock()
	if !disposed {
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
	o.disposedMutex.Lock()
	disposed := o.disposed
	o.disposedMutex.Unlock()
	if !disposed {
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
