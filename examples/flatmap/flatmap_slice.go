package main

import (
	"fmt"

	"github.com/reactivex/rxgo"
	"github.com/reactivex/rxgo/handlers"
)

func main() {
	primeSequence := rx.Just([]int{2, 3, 5, 7, 11, 13})

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
		Subscribe(handlers.NextFunc(func(prime interface{}) {
			fmt.Println("Prime -> ", prime)
		})).Block()
}
