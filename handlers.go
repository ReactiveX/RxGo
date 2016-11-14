package grx

// Handler is an interface which all Handler functions implement.
type Handler interface {
        Do(interface{})
}

type (
        // NextFunc handles emitted value.
        NextFunc func(v interface{})

        // ErrFunc handles emitted error.
        ErrFunc func(err error)

        // DoneFunc handles Observable's termination event.
        DoneFunc func()
)

// Do calls nextf(v).
func (nextf NextFunc) Do(v interface{}) {
        nextf(v)
}

// Do calls errf(v) if v is an error.
func (errf ErrFunc) Do(v interface{}) {
        if err, ok := v.(error); ok {
                errf(err)
        }
}

// Do calls donef() without any argument.
func (donef DoneFunc) Do(v interface{}) {
        donef()
}
