package observer

import (
	"github.com/reactivex/rxgo/handlers"
)


// Observer represents a group of EventHandlers.
type (
	Observer interface {
		OnNext(interface{})
		OnError(error)
		OnDone()
	}

	wrappingObserver struct {
		onNextHandler handlers.NextFunc
		onErrorHandler handlers.ErrFunc
		onDoneHandler handlers.DoneFunc
	}
)

func (o wrappingObserver) OnNext(item interface{}) {
	o.onNextHandler(item)
}

func (o wrappingObserver) OnError(err error) {
	o.onErrorHandler(err)
}
func (o wrappingObserver) OnDone() {
	o.onDoneHandler()
}

// New constructs a new Observer instance with default Observer and accept
// any number of EventHandler
func New(eventHandlers ...interface{}) Observer {
	ob := wrappingObserver{
		onNextHandler: func(item interface{}) {},
		onErrorHandler: func(err error) {},
		onDoneHandler: func() {},
	}

	if len(eventHandlers) > 0 {
		for _, handler := range eventHandlers {
			if customFunc, ok := handlers.AsNextFunc(handler); ok {
				ob.onNextHandler = customFunc
			} else if customFunc, ok := handlers.AsErrFunc(handler); ok {
				ob.onErrorHandler = customFunc
			} else if customFunc, ok := handlers.AsDoneFunc(handler); ok {
				ob.onDoneHandler = customFunc
			} else if customObserver, ok := handler.(Observer); ok {
				return customObserver
			}
		}
	}
	return ob
}
