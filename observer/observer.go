package observer

import (
	"github.com/jochasinga/grx/bases"
	"github.com/jochasinga/grx/handlers"
)

type Observer struct {
	//observable *subject.Subject
	sink        <-chan bases.Emitter
	NextHandler handlers.NextFunc
	ErrHandler  handlers.ErrFunc
	DoneHandler handlers.DoneFunc
}

/*
// DefaultObserver is a default Observable used by the constructor New.
// It makes sure no attribute is instantiated with nil that can cause a panic.
var DefaultObserver = func() *Observer {
	ob := &Observer{
		NextHandler: handlers.NextFunc(func(e bases.Item) {}),
		ErrHandler:  handlers.ErrFunc(func(err error) {}),
		DoneHandler: handlers.DoneFunc(func() {}),
	}
	ob.observable = subject.New(func(s *subject.Subject) {
		s.Sentinel = ob
	})
	return ob
}()
*/

// Apply makes Observer implements handlers.EventHandler
func (ob Observer) Handle(e bases.Emitter) {
	_, err := e.Emit()
	if err != nil {
		ob.ErrHandler.Handle(e)
		return
	}
	ob.NextHandler.Handle(e)
}

// New constructs a new Observer instance with default Observable
/*
func New(fs ...func(*Observer)) *Observer {
	ob := DefaultObserver
	if len(fs) > 0 {
		for _, f := range fs {
			f(ob)
		}
	}
	return ob
}
*/

// OnNext applies Observer's NextHandler to an Item
/*
func (ob *Observer) OnNext(item bases.Item) {
	if ob.NextHandler != nil {
		ob.NextHandler(item)
	}

		if handle, ok := ob.NextHandler.(handlers.NextFunc); ok {
			handle(item)
		}
}
*/

/*
// OnError applies Observer's ErrHandler to an error
func (ob *Observer) OnError(err error) {
	if ob.ErrHandler != nil {
		ob.ErrHandler(err)
	}

		if handle, ok := ob.ErrHandler.(handlers.ErrFunc); ok {
			handle(err)
		}

	return
}
*/

/*
// OnDone terminates the Observer's internal Observable
func (ob *Observer) OnDone() {
	ob.observable.Done()
	if ob.DoneHandler != nil {
		ob.DoneHandler()
	}

		if handle, ok := ob.DoneHandler.(handlers.DoneFunc); ok {
			handle()
		}
}
*/
