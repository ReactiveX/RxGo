package observable

import "github.com/reactivex/rxgo/observer"

// Creates observable from based on source function. Keep it mind to call emitter.OnDone()
// to signal sequence's end.
// Example:
// - emitting none elements
// observable.Create(emitter *observer.Observer) { emitter.OnDone() })
// - emitting one element
// observable.Create(func(emitter chan interface{}) {
//		emitter.OnNext("one element")
//		emitter.OnDone()
// })
func Create(source func(emitter *observer.Observer)) Observable {
	emitted := make(chan interface{})
	emitter := &observer.Observer{
		NextHandler: func(el interface{}) {
			emitted <- el
		},
		ErrHandler: func(err error) {
			// decide how to deal with errors
		},
		DoneHandler: func() {
			close(emitted)
		},
	}

	go func() {
		source(emitter)
	}()

	return emitted
}
