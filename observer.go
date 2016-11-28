package grx

// Observable is a built-in "Sentinel"
type Observer struct {
        observable  *Observable
        NextHandler NextFunc
        ErrHandler  ErrFunc
        DoneHandler DoneFunc
}

// New constructs a new Observer instance with default Observable
func NewObserver() *Observer {
        ob := &Observer{
                observable: &Observable{
                        C: make(chan interface{}),
                        observer: new(Observer),
                },
        }
        ob.observable.observer = ob
        return ob
}

// OnNext chucks a value into the Observer's internal Observable.
func (ob *Observer) OnNext(v interface{}) {
        ob.observable.C <- v
}

// OnError chucks an error into the Observer's internal Observable.
func (ob *Observer) OnError(err error) {
        ob.observable.C <- err
}

// OnDone terminates the Observer's internal Observable.
func (ob *Observer) OnDone() {
	ob.observable.done <- struct{}{}
        close(ob.observable.C)
}
