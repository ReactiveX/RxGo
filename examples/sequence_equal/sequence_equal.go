package main

import (
	"fmt"

	"github.com/reactivex/rxgo"
	"github.com/reactivex/rxgo/handlers"
)

func main() {
	sequence := rxgo.Just(2, "3", 5, 7, 11, 13)

	rxgo.
		Just(2, "3", 5, 7, 11, 13).
		SequenceEqual(sequence).
		Subscribe(handlers.NextFunc(func(result interface{}) {
			if result.(bool) {
				fmt.Println("Sequences are equal")
			} else {
				fmt.Println("Sequences are unequal")
			}
		})).Block()
}
