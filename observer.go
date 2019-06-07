package rxgo

import (
	"github.com/reactivex/rxgo/options"

	"github.com/reactivex/rxgo/handlers"
)

// Observer represents a group of EventHandlers.
type Observer interface {
	handlers.EventHandler
	Disposable

	OnNext(item interface{})
	OnError(err error)
	OnDone()

	Block() error
	setChannel(chan interface{})
	getChannel() chan interface{}
	setBackpressureStrategy(options.BackpressureStrategy)
	getBackpressureStrategy() options.BackpressureStrategy
}

type observer struct {
	disposed             chan struct{}
	nextHandler          handlers.NextFunc
	errHandler           handlers.ErrFunc
	doneHandler          handlers.DoneFunc
	done                 chan error
	channel              chan interface{}
	backpressureStrategy options.BackpressureStrategy
	buffer               int
}

func (o *observer) setBackpressureStrategy(strategy options.BackpressureStrategy) {
	o.backpressureStrategy = strategy
}

func (o *observer) getBackpressureStrategy() options.BackpressureStrategy {
	return o.backpressureStrategy
}

func (o *observer) setChannel(ch chan interface{}) {
	o.channel = ch
}

func (o *observer) getChannel() chan interface{} {
	return o.channel
}

// NewObserver constructs a new Observer instance with default Observer and accept
// any number of EventHandler
func NewObserver(eventHandlers ...handlers.EventHandler) Observer {
	ob := observer{
		disposed: make(chan struct{}),
	}

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
	close(o.disposed)
}

func (o *observer) IsDisposed() bool {
	select {
	case <-o.disposed:
		return true
	default:
		return false
	}
}

// OnNext applies Observer's NextHandler to an Item
func (o *observer) OnNext(item interface{}) {
	if !o.IsDisposed() {
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

// OnDone terminates the Observer's internal Observable
func (o *observer) OnDone() {
	if !o.IsDisposed() {
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
	if !o.IsDisposed() {
		for v := range o.done {
			return v
		}
	}
	return nil
}
