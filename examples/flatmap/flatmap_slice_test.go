package main

import (
	"testing"

	"github.com/reactivex/rxgo"
)

func TestFlatMapExample(t *testing.T) {
	// given
	observerMock := rxgo.NewObserverMock()

	// and
	primeSequence := rxgo.Just([]int{2, 3, 5, 7, 11, 13})

	// when
	primeSequence.
		FlatMap(func(primes interface{}) rxgo.Observable {
			return rxgo.Create(func(emitter rxgo.Observer, disposed bool) {
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
