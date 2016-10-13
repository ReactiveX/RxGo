package rx

import (
//	"errors"
	"sync"
	"fmt"
)

type Observer struct {
	OnNext func(e Event)
	OnError func(e Event)
	OnCompleted func(e Event)
}

// Context is a general context map
type Context map[string]interface{}

func (cxt Context) Add(key string, val interface{}) string {
	cxt[key] = val
	return key
}

type Event struct {
	//Value     interface{}
	Value     interface{}
	Error     error
	Completed bool
}

/*
func (e Event) getValue() (Context, error) {
	val := e.Value
	if val == nil {
		return nil, errors.New("Event.Value is nil.")
	}

	typecxt := make(Context)

	switch val := val.(type) {
	default:
		return nil, errors.New("Impossible to get type of Event.Value.")
	case bool:
		typecxt.Add("b", val)
	case string:
		typecxt.Add("s", val)
	case int:
		typecxt.Add("i", val)
	case int32:
		typecxt.Add("i32", val)
	case int64:
		typecxt.Add("i64", val)
	case float32:
		typecxt.Add("f32", val)
	case float64:
		typecxt.Add("f64", val)
	}

	return typecxt, nil
}
*/

type Observable struct {
	Name string
	//Duplex map[string]chan Event
	Stream   chan Event
	//Observer Observer
	Observer *Observer
}

func (o *Observable) Add(ev Event) {
	go func() {
		o.Stream <- ev
	}()
}

// Just create an <Observable> with only one item and emit it "as-is".
func Just(item interface{}) *Observable {
	o := &Observable{
		Name:   "",
		Stream: make(chan Event, 1),
	}

	o.Stream <- Event{Value: item}
	close(o.Stream)
	return o
}

func From(iterable []interface{}) *Observable {
	o := &Observable{
		Name:   "",
		Stream: make(chan Event, len(iterable)),
	}
	var wg sync.WaitGroup
	wg.Add(len(iterable))
	for _, item := range iterable {
		go func(item interface{}) {
			defer wg.Done()
			o.Stream <-Event{Value: item}
		}(item)
	}
	wg.Wait()
	close(o.Stream)

	return o
}

func Start(fx func() Event) *Observable {
	o := &Observable{ Stream: make(chan Event, 1) }
	o.Stream <- fx()
	close(o.Stream)
	return o
}

func (o *Observable) Subscribe(ob *Observer) {
	//var wg sync.WaitGroup
	o.Observer = ob
	
	for ev := range o.Stream {
		/*
		wg.Add(1)
		go func(ev *Event) {
			defer wg.Done()
			o.Observer(e)
		}(e)
                */
		if ev.Error != nil {
			o.Observer.OnError(ev)
		} else if ev.Completed {
			o.Observer.OnCompleted(ev)
		} else {
			fmt.Println(ev)
			o.Observer.OnNext(ev)
		}
	}
	//wg.Wait()
}

func NewObservable(name string) *Observable {
	return &Observable{
		Name: name,
		//Duplex: make(map[string]chan Event),
		Stream: make(chan Event),
	}
}
