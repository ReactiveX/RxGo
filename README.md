# Go Reactive Extension (GRX)
Reactive extension for Go.

## Why?
Many have the tendency to stay with the "raw" Go and ward off any attempts of abstracting the language to a higher level. This mentality creates a divide between the experienced system programmers and new programmers who want to learn to code, possibly in Go. With `GRX`, concurrency primitives (goroutines, channels, syncs in Go) are abstracted in favor of a more intuitive, graphical nature of event streams.

## What is Reactive Programming?
**Reactive** is a buzzword. Yes, you heard right. But it is also a yet-to-be-refined programming *paradigm*--way of doing things--like the MVC was decades ago. A chat application like Slack using websocket and any real-time game server are reactive. Reactive is how a system is always "polling" for any change in any variable that when it occurs will be instanteneously and asynchronously (without waiting) *propagated* or reflected, very possibly to the user of the system or affecting another variable.

While there is a ![Reactive Manifesto](http://www.reactivemanifesto.org/), it is still open to interpretations.

## ReactiveX 
The most well-known interpretation being ![ReactiveX](http://reactivex.io)'s, which consists of concepts like streams, observables, events, and etc. (This is what **GRX** is based on). There are other great interpretations, such as data-flow concept, in which if you're interested do check out ![go-flow](https://github.com/trustmaster/goflow). Pretty interesting especially if you know interactive programming in Max/MSP or Pure Data.

## 'Nuf said, show me code
Not before I warn you that this is  a work in progress. Here's some very basic examples.

```go

package main
import (
	"fmt"
	"grx"
	"reflect"
)

func main() {
	nums := []int{1, 2, 3}
	xnums := []int{}

	// Create an observable (stream) from a slice of integers.
	numStream := grx.From(nums)

	// Create an Observer object.
	observer := &Observer{ 
	
		// While there is more events down the stream, times two to the value emitted and append to xnums.
		OnNext: func(e Event) { xnums = append(xnums, e.Value.(int) * 2) },
		
		// If an event emits an error, panic (in this case there is no way an error can occur).
		OnError: func(e Event) { panic(e.Error) },
		
		// If the observable is about to end (on last event), append 0 to xnums.
		OnCompleted: func(e Event) { xnums = append(xnums, 0) },
	}

	// Start watching the stream
	numStream.Subscribe(observer)
	
	fmt.Println(reflect.DeepEqual(xnums, []int{2, 4, 6, 0})) // true
}

```


