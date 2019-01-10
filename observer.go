package rxgo

import (
	"sync"

	"github.com/reactivex/rxgo/v2/handlers"
)

// Observer represents a group of EventHandlers.
type Observer interface {
	handlers.EventHandler
	Disposable

	OnNext(item interface{})
	OnError(err error)
	OnDone()

	Block() error
}

type observer struct {
	disposedMutex sync.Mutex
	disposed      bool
	nextHandler   handlers.NextFunc
	errHandler    handlers.ErrFunc
	doneHandler   handlers.DoneFunc
	done          chan error
}

// NewObserver constructs a new Observer instance with default Observer and accept
// any number of EventHandler
func NewObserver(eventHandlers ...handlers.EventHandler) Observer {
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
	ob.done = make(chan error, 1)

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

func (o *observer) Dispose() {
	o.disposedMutex.Lock()
	o.disposed = true
	o.disposedMutex.Unlock()
}

func (o *observer) IsDisposed() bool {
	o.disposedMutex.Lock()
	defer o.disposedMutex.Unlock()
	return o.disposed
}

// OnNext applies Observer's NextHandler to an Item
func (o *observer) OnNext(item interface{}) {
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
			}
		}
	} else {
		// TODO
	}
}

// OnError applies Observer's ErrHandler to an error
func (o *observer) OnError(err error) {
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

// OnDone terminates the Observer's internal Observable
func (o *observer) OnDone() {
	o.disposedMutex.Lock()
	disposed := o.disposed
	o.disposedMutex.Unlock()
	if !disposed {
		if o.doneHandler != nil {
			o.doneHandler()
			o.Dispose()
			if o.done != nil {
				o.done <- nil
				close(o.done)
			}
		}
	} else {
		// TODO
	}
}

// OnDone terminates the Observer's internal Observable
func (o *observer) Block() error {
	o.disposedMutex.Lock()
	disposed := o.disposed
	o.disposedMutex.Unlock()
	if !disposed {
		for v := range o.done {
			return v
		}
	}
	return nil
}
