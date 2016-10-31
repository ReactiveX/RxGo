# Go Reactive Extension (GRX)
Simple Reactive Extension (RX) for Go.

## Why?
I started `GRX` as I studied Reactive Programming, or specifically [ReactiveX](http://reactivex.io)'s interpretation of Reactive Programming.

## Is this for you?
If you are someone who wants to learn or use Reactive Programming in Go, and you are open-mind about concurrency primitives such as goroutines, channels, and syncs being abstracted in favor of the [Stream mantra](https://camo.githubusercontent.com/e581baffb3db3e4f749350326af32de8d5ba4363/687474703a2f2f692e696d6775722e636f6d2f4149696d5138432e6a7067), you may find this interesting. 

## What is Reactive Programming?
**Reactive** is a yet-to-be-refined programming *paradigm*--or way of doing things--like the MVC was decades ago. One can say that a chat application (like Slack) or a real-time game server is "reactive". Then you might ask, "what's new?" Nothing is, really. It is the way of looking at the problem and programming that is.

## The Procedural Perception
When you write code, it is hard to break the procedural perception (code being executed/interpreted from the first line all the way down). With an event-driven code in language like JavaScript, the callbacks become hard to maintain and make sense. This is simply because it is natural to humans' perception to read from top to bottom.

However, if you have programmed in a visual flow-based language like [Max/MSP](https://cycling74.com/products/max/#.WAGV0dwgd0I) or [Pure Data](https://puredata.info/), you already have programmed a "reactive" system, so to speak. In such systems, data flows around the system via cords, which connect units, patches, or functions (whatever you'd like to call them) which are self-contained, take one or more inputs, process it, and return one or more outputs. The idea is the system is always updated live, and a change in **anything** within the system will be propagated accordingly without the need to refresh or recompile.

> If you are interested in this flow-based interpretation of Reactive Programming in Go, check out [go-flow](https://github.com/trustmaster/goflow).

## ReactiveX
This package is based on [ReactiveX](http://reactivex.io)'s Observable/Stream Pattern, which consists of concepts like event streams or observables, observers, events, and etc. What is considered Reactive Programming is still murky, but ReactiveX is by far the most interesting one.

## Usage

```go

package main
import (
	"fmt"
	"time"
	
	"github.com/jochasinga/grx/observable"
    "github.com/jochasinga/grx/observer"
	"github.com/jochasinga/grx/event"
)

func main() {

	// Create an observable (stream) from a single item.
	myStream := observable.Just(1)
	score := 9

	// Create an observer
	myObserver := &observer.Observer{
		OnNext: func(e *event.Event) {
			score = score + e.Value.(int)
		},
	}

	// Subscribe the observer to watch over the stream
	myStream.Subscribe(myObserver)

	// Block/wait here a bit for score to update.
	<-time.After(100)

	fmt.Println(score) // 10
}

```

An `Observable` is a synchronous stream of `Event`s which can emit a `Value`, `Error`,
or notify as `Completed`. Below is what an `Observable` created looks like:

```bash

              time -->

(*)-------------(e)--------------|>
 |               |               |
Start        Event with     Terminated
             value = 1

```

An `Observer` watches over an `Observable` with `OnNext`, `OnError`, and `OnCompleted` methods, which get called for each `Event` it encounters. `OnNext` is called zero or more times to handle each `Event` that emits a value. `OnError` or `OnCompleted` is called, respectively, when an `Event` emits an error or a completed signal.

These are done while you're doing something else. When you subscribe an `Observer` to an `Observable`, it is asynchronous. Thus, that is why in our code we waited with `<-time.After()` before printing the score. It didn't block.

Here is another example, this time creating an `Observable` from a slice with `From` operator:

```go

package main
import (
	"fmt"
	"reflect"
	"time"
	
	"github.com/jochasinga/grx/observable"
    "github.com/jochasinga/grx/observer"
	"github.com/jochasinga/grx/event"
)

func main() {
	nums := []interface{}{1, 2, 3}
	xnums := []int{}

	// Create an observable from a slice of integers.
	numStream := observer.From(nums)

	// Create an Observer object.
	myObserver := &Observer{
		OnNext: func(e *event.Event) {
			xnums = append(xnums, e.Value.(int) * 2)
		},
		OnError: func(e *event.Event) {
			panic(e.Error)
		},
		OnCompleted: func(e *Event) {
            xnums = append(xnums, 0)
		},
	}

	// Start watching the stream
	numStream.Subscribe(myObserver)

	// Block/wait for Observer to interact with the stream
	<-time.After(100)

	fmt.Println(reflect.DeepEqual(xnums, []int{2, 4, 6, 0})) // true
}

```

The above code would have been visualized this way:

```bash

                     time -->

(*)-------(e)----------(e)----------(e)-----------|>
 |         |            |            |            |
Start   Event with   Event with    Event with  Terminated
        value = 1    value = 2     value = 3

```

One of the concepts being ReactiveX's stream (or at least what I think) is to abstract all the concurrency primitives in a language, and instead having the user to focus on events and streams, not threads, coroutines, goroutines, etc.

Most Observable methods and operators will return the Observable itself, making it chainable.

```go

f1 := func() *event.Event {
	time.Sleep(2 * time.Second)
	return &event.Event{ Value: 1 }
}

f2 := func() *event.Event {
	time.Sleep(time.Second)
	return &event.Event{ Value: 2 }
}

myObserver := &observer.Observer{
	OnNext: func(e *event.Event) { fmt.Println(e.Value) },
}

myStream := grx.Start(f1, f2).Subscribe(myObserver)

// 2 printed
// 1 printed

```
