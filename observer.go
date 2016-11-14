package grx

import (
        //"github.com/jochasinga/grx"
        //"github.com/jochasinga/grx/observable"
)

// Observable is a built-in "Sentinel"
type Observer struct {
        Observable  *Observable
        NextHandler NextFunc
        ErrHandler  ErrFunc
        DoneHandler DoneFunc
}

// New constructs a new Observer instance with default Observable
func NewObserver() *Observer {
        ob := &Observer{
                Observable: &Observable{
                        Stream: make(chan interface{}),
                        Observer: new(Observer),
                },
        }
        ob.Observable.Observer = ob
        return ob
}

// OnNext chucks a value into the Observer's internal Observable.
func (ob *Observer) OnNext(v interface{}) {
        ob.Observable.Stream <- v
}

// OnError chucks an error into the Observer's internal Observable.
func (ob *Observer) OnError(err error) {
        ob.Observable.Stream <- err
}

// OnDone terminates the Observer's internal Observable.
func (ob *Observer) OnDone() {
        close(ob.Observable.Stream)
}
