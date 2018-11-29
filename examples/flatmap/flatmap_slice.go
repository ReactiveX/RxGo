package main

import (
	"fmt"

	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/observable"
)

func main() {
	primeSequence := observable.Just([]int{2, 3, 5, 7, 11, 13})

	<-primeSequence.
		FlatMap(func(primes interface{}) observable.Observable {
			return observable.Create(func(emitter observable.Observer, disposed bool) {
				for _, prime := range primes.([]int) {
					emitter.OnNext(prime)
				}
				emitter.OnDone()
			})
		}, 1).
		Last().
		Subscribe(handlers.NextFunc(func(prime interface{}) {
			fmt.Println("Prime -> ", prime)
		}))
}
