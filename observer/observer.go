package observer

import (
	"github.com/jochasinga/grx"
	"github.com/jochasinga/bases"
	"github.com/jochasinga/grx/handlers"
)

type Observer struct {
	observable  *grx.Subject
	NextHandler handlers.EventHandler
	ErrHandler  handlers.EventHandler
	DoneHandler handlers.EventHandler
}

// DefaultObserver is a default Observable used by the constructor New.
// It makes sure no attribute is instantiated with nil that can cause a panic.
var DefaultObserver = func() *Observer {
	ob := &Observer{
		NextHandler: handlers.NextFunc(func(e bases.Emitter){}),
		ErrHandler: handlers.ErrFunc(func(err error){}),
		DoneHandler: handlers.DoneFunc(func(){})
		notifier:    bang.New(),
	}
	ob.observable = subject.New(func(s *subject.Subject) {
		s.Observer = o
	})
	return o
}()

// Apply makes Observer implements handlers.EventHandler
func (ob *Observer) Apply(e bases.Emitter) {
	if ob.observable.HasNext() {
		item, err := e.Emit() 
		if err != nil {
			if item == nil {
				ob.ErrHandler(err)
			}
		} else {
			if item != nil {
				ob.NextHandler(item)
			}
		}
	} else {
		observer.DoneHandler()
	}
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
	if ob.NextHandler != nil {
		ob.NextHandler(item)
	}
}

// OnError applies Observer's ErrHandler to an error
func (ob *Observer) OnError(err error) {
	if ob.ErrHandler != nil {
		ob.ErrHandler(err)
	}
}

// OnDone terminates the Observer's internal Observable
func (ob *Observer) OnDone() {
	ob.observable.Done()
	if ob.DoneHandler != nil {
		ob.DoneHandler()
	}
}
