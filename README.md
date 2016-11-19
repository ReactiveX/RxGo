# Go Reactive Extensions (grx)
Reactive Extensions (Rx) for the Go Language

## Getting Started
[ReactiveX](http://reactivex.io/), or Rx for short, is an API for programming with observable event streams. This is a port of ReactiveX to Go.

Rx is a new, alternative way of asychronous programming to callbacks, promises and deferred. It is about processing streams of events, with events being any occurances or changes within the system, either influenced by the external factors (i.e. users or another remote service) or internal components (i.e. logs).

The pattern is that you `Subscribe` to an `Observable` using an `Observer`:

```go

subscription := observable.Subscribe(observer)

```

**NOTE**: Observables are not active in themselves. They need to be subscribed to make something happen. Simply having an Observable lying around doesn't make anything happen, like sitting and watching time flies.

## Install

```bash

go get -u github.com/jochasinga/grx

```

## Importing the Rx package
```go

import "github.com/jochasinga/grx"

```
## Simple Usage
```go

watcher := &grx.Observer{

	// Register a handler function for every emitted value.
	NextHandler: grx.NextFunc(func (x interface{}) {
		fmt.Printf("Got: %s\n", x.(int))
	}),
}

source := grx.From([]int{1, 2, 3, 4, 5})
_, err := source.Subscribe(watcher)
if err != nil {

	// This error is return right away if Observable is nil.
	panic(err)
}

```

The above will print the format string for every number in the slice.

```go

package main
import (
	"fmt"
	"time"

	"github.com/jochasinga/grx"
)

func main() {

	score := 9

	watcher := &Observer{
	
		// Register a handler function for each emitted value.
		NextHandler: grx.NextFunc(func(v interface{}) {
			score +=  v.(int)
		}),
		
		// Register a "clean-up" function when the observable terminates.
		DoneHandler: grx.Donefunc(func() {
			score *= 2
		}),
	}

	// Create an observable from a single item and subscribe to the observer.
	_, _ := observable.Just(1).Subscribe(myObserver)

	// Block/wait here a bit for score to update before printing out.
	<-time.After(100 * time.Millisecond)

	fmt.Println(score) // 20
}

```

An `Observable` is a synchronous stream of evens which can emit a value of type `interface{}`, `error`,
or notify as completed. Below is what an `Observable` looks like:

```bash

              time -->

(*)-------------(e)--------------|>
 |               |               |
Start        Event with         Done
             value = 1

```

An `Observer` watches over an `Observable` with a set of handlers: `NextHandler`, `ErrHandler`, and `DoneHandler`, which get called for each event it encounters. `NextHandler` is called zero or more times to handle each event that emits a value, that is when there is still a next event. `ErrHandler` or `DoneHandler` is called, respectively, when an event emits an error or the Observable is finished. When you subscribe an `Observer` to an `Observable`, it starts another non-blocking goroutine. Thus, that is why in the previous code it was necessary to block with `<-time.After()` before printing the score.

There is a few ways you can create an `Observable` and subscribe it. Here is a different one using the operator `CreateObservable` and subscribing a set of handler functions with `SubscribeWith` instead of an `Observer`.

```go

source := grx.CreateObservable(func(ob *grx.Observer) {
	ob.OnNext("Hello")
}

nextf := func(v interface{}) { 
	fmt.Println(v)
}

_, _ = source.SubscribeFunc(nextf, nil, nil)

```
Most Observable methods and operators will return the Observable itself, making it chainable.

```go

f1 := func() interface{} {
	time.Sleep(2 * time.Second)
	return 1 
}

f2 := func() interface{} {
	time.Sleep(time.Second)
	return 2
}

myObserver := &grx.Observer{
	NextHandler: grx.NextFunc(func(v interface{}) { 
		fmt.Println(v) 
	}),
}

myStream := observable.Start(f1, f2).Subscribe(myObserver)

// Block/wait for Observer to interact with the stream
<-time.After(100 * time.Millisecond)

// 2 printed
// 1 printed

```

## Is this Idiomatic Go?
It depends. This is certainly not for [purists who may regard this as a violation
to Go's core philosophy](https://www.reddit.com/r/golang/comments/5d1erq/reactive_extension_for_go/). 
Channels are the underlying implementations of almost all methods and operators. 
In fact, `Observable` is basically a channel. The goal of this extension is just to expose
a set of APIs complying to ReactiveX's way of programming instead of managing concurrency with
primitives like channels and goroutines. However, they can always be used alongside one another 
(check the examples).

## This is a very early project and there will be breaking changes.
