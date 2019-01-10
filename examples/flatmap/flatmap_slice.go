package main

import (
	"fmt"

	"github.com/reactivex/rxgo/v2"
	"github.com/reactivex/rxgo/v2/handlers"
)

func main() {
	primeSequence := rxgo.Just([]int{2, 3, 5, 7, 11, 13})

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
		Subscribe(handlers.NextFunc(func(prime interface{}) {
			fmt.Println("Prime -> ", prime)
		})).Block()
}
