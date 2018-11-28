package observable

import (
	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/observer"
)

// Creates observable from based on source function. Keep it mind to call emitter.OnDone()
// to signal sequence's end.
// Example:
// - emitting none elements
// observable.Create(emitter observer.Observer, disposed bool) { emitter.OnDone() })
// - emitting one element
// observable.Create(func(emitter observer.Observer, disposed bool) {
//		emitter.OnNext("one element")
//		emitter.OnDone()
// })
func Create(source func(emitter observer.Observer, disposed bool)) Observable {
	emitted := make(chan interface{})
	emitter := observer.New(
		handlers.NextFunc(func(el interface{}) {
			if !isClosed(emitted) {
				emitted <- el
			}
		}), handlers.ErrFunc(func(err error) {
			// decide how to deal with errors
			if !isClosed(emitted) {
				close(emitted)
			}
		}), handlers.DoneFunc(func() {
			if !isClosed(emitted) {
				close(emitted)
			}
		}),
	)

	go func() {
		source(emitter, isClosed(emitted))
	}()

	return &observator{
		ch: emitted,
	}
}

func isClosed(ch <-chan interface{}) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}
