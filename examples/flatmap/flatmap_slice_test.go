package main

import (
	"github.com/reactivex/rxgo"
	"testing"
)

func TestFlatMapExample(t *testing.T) {
	// given
	observerMock := rx.NewObserverMock()

	// and
	primeSequence := rx.Just([]int{2, 3, 5, 7, 11, 13})

	// when
	primeSequence.
		FlatMap(func(primes interface{}) rx.Observable {
			return rx.Create(func(emitter rx.Observer, disposed bool) {
				for _, prime := range primes.([]int) {
					emitter.OnNext(prime)
				}
				emitter.OnDone()
			})
		}, 1).
		Last().
		Subscribe(observerMock.Capture()).Block()

	// then
	observerMock.AssertCalled(t, "OnNext", 13)
}
