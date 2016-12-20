package grx

// Handler is an interface which all Handler functions implement.
type EventHandler interface {
	Apply(interface{})
}

type (
	NextFunc func(v interface{})
	ErrFunc  func(err error)
	DoneFunc func()
)

// Do calls nextf(v)
func (nextf NextFunc) Apply(v interface{}) {
	nextf(v)
}

// Do calls errf(v) if v is an error
func (errf ErrFunc) Apply(v interface{}) {
	if err, ok := v.(error); ok {
		errf(err)
	}
}

// Do calls donef() without any argument.
func (donef DoneFunc) Apply(v interface{}) {
	donef()
}

func (observer *BaseObserver) Apply(v interface{}) {
	if observer._observable.HasNext() {
		switch next := v.(type) {
		case error:
			if observer.ErrHandler != nil {
				observer.ErrHandler(next)
			}
		default:
			if observer.NextHandler != nil {
				observer.NextHandler(next)
			}
		}
	} else {
		if observer.DoneHandler != nil {
			observer.DoneHandler()
		}
	}
}
