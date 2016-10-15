# Go Reactive Extension (GRX)
Reactive extension for Go.

## Why?
I started `GRX` as a study into Reactive Programming, or specifically ![ReactiveX](http://reactivex.io)'s interpretation of Reactive Programming.

## Who?
If you want to learn or use Reactive Programming, and you want to do it in Go, this might interest you. If you do not mind those concurrency primitives such as goroutines, channels, and syncs being abstracted away in favor of the ![Stream mantra](https://camo.githubusercontent.com/e581baffb3db3e4f749350326af32de8d5ba4363/687474703a2f2f692e696d6775722e636f6d2f4149696d5138432e6a7067), you may find this interesting.

## What is Reactive Programming?
**Reactive** is a yet-to-be-refined programming *paradigm*--or way of doing things--like the MVC was decades ago. One can say that a chat application (like Slack) or a real-time game server is "reactive". Then you might ask, "what's new?" Nothing is. It is the unique way of looking an it that is.

When you write code, it is hard to break the procedural perception (code being executed/interpreted from the first line all the way down). Even with an event-driven language like JavaScript, it still is. This is simply because humans' perception of reading from top to bottom.

If you have programmed in visual flow-based language like ![Max/MSP](https://cycling74.com/products/max/#.WAGV0dwgd0I) or ![Pure Data](https://puredata.info/), the difference becomes clear. In such systems, data flows around the system via cords, which connect units, patches, or functions (whatever you'd like to call them). The idea is the system is always updated live, and a change in anything within the system will be propagated accordingly without the need to refresh or recompile.

> If you are interested in this flow-based interpretation of Reactive Programming in Go, > check out ![go-flow](https://github.com/trustmaster/goflow).

## ReactiveX
This package is based on ![ReactiveX](http://reactivex.io)'s Observable/Stream Pattern, which consists of concepts like event streams or observables, observers, events, and etc.

## 'Nuf said, show me the code
Not before I warn you that this is a work in progress. How long I keep working on this depends on the stars and forks from you!

**Here's some basic examples.**

```go

package main
import (
	"fmt"
	"grx"
)

func main() {
	// Create an observable (stream) from a single item.
	observable := grx.Just(1)
	score := 9

	// Create an observer
	observer := &Observer{
		OnNext: func(e *Event) {
			score = score + e.Value.(int)
		},
	}

	// Subscribe the observer to watch over the stream
	observable.Subscribe(observer)

	fmt.Println(score) // 10
}

```

An *Observable* is a synchronous stream of Events which can emit a value, error,
or notify as completed. Below is what `Just(item interface{})` just did:

```bash

              time -->

(*)-------------(e)--------------|>
 |               |               |
Start        Event *e*      Terminated
             with val=1

```

An `Observer` watches over an `Observable` with `OnNext`, `OnError`, and `OnCompleted` methods, which get called by each `Event` it encounters. `OnNext` is called zero or more times to handle each event that emits a value. `OnError` and `OnCompleted` are called, respectively, when an event emits an error and a completed signal.

These are done while you do something else. When you subscribe an Observer to an Observable, it is asynchronous and you can forget about it.

Here is another example:

```go

package main
import (
	"fmt"
	"grx"
	"reflect"
)

func main() {
	nums := []interface{}{1, 2, 3}
	xnums := []int{}

	// Create an observable from a slice of integers.
	numStream := grx.From(nums)

	// Create an Observer object.
	observer := &Observer{

		// While there is more events down the stream, times two to the value emitted and append to xnums.
		OnNext: func(e Event) { xnums = append(xnums, e.Value.(int) * 2) },

		// If an event emits an error, panic (in this case it's not possible to get an error).
		OnError: func(e Event) { panic(e.Error) },

		// If the observable is about to end (on last event), append 0 to xnums.
		OnCompleted: func(e Event) { xnums = append(xnums, 0) },
	}

	// Start watching the stream
	numStream.Subscribe(observer)

        // Do something else meanwhile

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
