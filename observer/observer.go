package observer

import (
	"github.com/jochasinga/grx/bases"
	"github.com/jochasinga/grx/handlers"
	"github.com/jochasinga/grx/subject"
)

type Observer struct {
	observable  *subject.Subject
	NextHandler bases.EventHandler
	ErrHandler  bases.EventHandler
	DoneHandler bases.EventHandler
}

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

// Apply makes Observer implements handlers.EventHandler
func (ob *Observer) Handle(e bases.Emitter) {
	_, err := e.Emit()
	if err != nil {
		ob.ErrHandler.Handle(e)
	}
	ob.NextHandler.Handle(e)
}

// New constructs a new Observer instance with default Observable
func New(fs ...func(*Observer)) *Observer {
	ob := DefaultObserver
	if len(fs) > 0 {
		for _, f := range fs {
			f(ob)
		}
	}
	return ob
}

// OnNext applies Observer's NextHandler to an Item
func (ob *Observer) OnNext(item bases.Item) {
	if handle, ok := ob.NextHandler.(handlers.NextFunc); ok {
		handle(item)
	}
}

// OnError applies Observer's ErrHandler to an error
func (ob *Observer) OnError(err error) {
	if handle, ok := ob.ErrHandler.(handlers.ErrFunc); ok {
		handle(err)
	}
	return
}

// OnDone terminates the Observer's internal Observable
func (ob *Observer) OnDone() {
	ob.observable.Done()
	if handle, ok := ob.DoneHandler.(handlers.DoneFunc); ok {
		handle()
	}
}
