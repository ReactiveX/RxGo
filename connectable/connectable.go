package connectable

import (
	"sync"

	"github.com/jochasinga/grx/bases"
	"github.com/jochasinga/grx/observable"
	"github.com/jochasinga/grx/observer"
)

type (
	Directive      func() bases.Emitter
	MappableFunc   func(bases.Emitter) bases.Emitter
	CurryableFunc  func(interface{}) MappableFunc
	FilterableFunc func(bases.Emitter) bool
	ReduceableFunc func(bases.Emitter) bases.Emitter
)

type Connectable struct {
	observable.Basic
	observers []observer.Observer
}

func New(buffer uint, observers ...observer.Observer) Connectable {
	return Connectable{
		Basic:     make(chan bases.Emitter, int(buffer)),
		observers: observers,
	}
}

func From(es []bases.Emitter) Connectable {
	source := make(chan bases.Emitter)
	go func() {
		for _, e := range es {
			source <- e
		}
		close(source)
	}()
	return Connectable{Basic: source}
}

func (co Connectable) Subscribe(ob observer.Observer) Connectable {
	co.observers = append(co.observers, ob)
	return co
}

func (co Connectable) Connect() <-chan struct{} {
	done := make(chan struct{}, 1)
	var wg sync.WaitGroup
	for _, ob := range co.observers {
		wg.Add(1)
		go func(ob observer.Observer) {
			for e := range co.Basic {
				ob.NextHandler.Handle(e)
			}
			ob.DoneHandler()
			wg.Done()
		}(ob)
	}
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()
	return done
}

func (co Connectable) Map(fx MappableFunc) Connectable {
	source := make(chan bases.Emitter, len(co.Basic))
	go func() {
		for e := range co.Basic {
			source <- fx(e)
		}
		close(source)
	}()
	return Connectable{Basic: source}
}

func (co Connectable) Filter(fx FilterableFunc) Connectable {
	source := make(chan bases.Emitter, len(co.Basic))
	go func() {
		for e := range co.Basic {
			if fx(e) {
				source <- e
			}
		}
		close(source)
	}()
	return Connectable{Basic: source}
}
