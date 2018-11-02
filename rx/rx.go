package rx

import (
	"context"
	"time"
)

// rx contains high-level interface types that can be extended

// Iterator IS anything that CAN get its Next value and return it and an error
// suggesting the end of its sequence.
type Iterator interface {
	Next() (interface{}, error)
}

// An Emitter IS anything that CAN Emit a value asynchronously at any point in time
type Emitter interface {
	context.Context
	SubscribeAt() time.Time
	// StartOn() time.Time
	// Deadline() (time.Time, bool)
	// Done() <-chan struct{}
	// Err() error
}

// Observable IS an Iterator that CAN Subscribe one or more
// Observer(s) and immediately return an Emitter (see Emitter).
type Observable interface {
	Iterator
	//Iterable
	Subscribe(Observer) Emitter
}

// An Observer IS a Handler that CAN Observe an Observable and Handle its emitted value(s).
// Example:
//
// type ObserverA struct { }
// func (oa *ObserverA) Handle(v interface{}) {
// 	// for every emission, just print fizz
// 	fmt.Println("fizz")
// }

// type ObserverB func(int) int
// var funcObserver = ObserverB(func(n int) int {
//         return n * 2
// })
// func (ob ObserverB) Handle(v interface{}) {
//         // for every emission that's an int, multiply it by 2 and print the product
//         if n, ok := v.(int); ok {
//                 num := ob(n)
//                 fmt.Println(num)
//         }
// }
//
//

type Observer interface {
	Handler
	// Observe(Observable)
	OnNext(interface{})
	OnError(error)
	OnDone()
}

// A Handler IS anything that CAN Handle a value and create one or more side effect(s).
type Handler interface {
	Handle(interface{})
}
