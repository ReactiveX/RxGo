package grx

// Observable is a built-in "Sentinel"
type BaseObserver struct {
	_observable Observable
	NextHandler NextFunc
	ErrHandler  ErrFunc
	DoneHandler DoneFunc
}

// New constructs a new Observer instance with default Observable
func NewBaseObserver() *BaseObserver {
	ob := &BaseObserver{
		_observable: Observable(&BaseObservable{
			C:         make(chan interface{}),
			_observer: (Observer)(nil),
		}),
	}
	//ob.observable.observer = ob
	ob._observable.setInnerObserverTo(Observer(ob))
	return ob
}

func (ob *BaseObserver) setInnerObservableTo(o Observable) Observable {
	ob._observable = o
	return ob._observable
}

func (ob *BaseObserver) getInnerObservable() Observable {
	return ob._observable
}

// OnNext chucks a value into the Observer's internal Observable.
func (ob *BaseObserver) OnNext(v interface{}) {
	//ob._observable.C <- v
	ob._observable.addItem(v)
	/*
		if ob.NextHandler != nil {
			ob.NextHandler(v)
		}
	*/
}

// OnError chucks an error into the Observer's internal Observable.
func (ob *BaseObserver) OnError(err error) {
	//ob.observable.C <- err
	//ob._observable.addItem(err)
	if ob.ErrHandler != nil {
		ob.ErrHandler(err)
	}
}

// OnDone terminates the Observer's internal Observable.
func (ob *BaseObserver) OnDone() {
	/*
			ob._observable.done <- struct{}{}
		        close(ob.observable.C)
	*/
	//ob._observable.done()
	if ob.DoneHandler != nil {
		ob.DoneHandler()
	}
}
