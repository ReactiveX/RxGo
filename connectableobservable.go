package rxgo

import (
	"sync"

	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/options"
)

type ConnectableObservable interface {
	Iterable
	Connect() Observer
	Subscribe(handler handlers.EventHandler, opts ...options.Option) Observer
}

type connectableObservable struct {
	iterator   Iterator
	observable Observable
	observers  []Observer
}

func newConnectableObservableFromObservable(observable Observable) ConnectableObservable {
	return &connectableObservable{
		observable: observable,
		iterator:   observable.Iterator(),
	}
}

func (c *connectableObservable) Iterator() Iterator {
	return c.iterator
}

func (c *connectableObservable) Subscribe(handler handlers.EventHandler, opts ...options.Option) Observer {
	ob := CheckEventHandler(handler)
	c.observers = append(c.observers, ob)
	return ob
}

func (c *connectableObservable) Connect() Observer {
	source := make([]interface{}, 0)

	it := c.iterator
	for it.Next() {
		item := it.Value()
		source = append(source, item)
	}

	var wg sync.WaitGroup

	for _, ob := range c.observers {
		wg.Add(1)
		local := make([]interface{}, len(source))
		copy(local, source)

		go func(ob Observer) {
			defer wg.Done()
			var e error
		OuterLoop:
			for _, item := range local {
				switch item := item.(type) {
				case error:
					ob.OnError(item)

					// Record error
					e = item
					break OuterLoop
				default:
					ob.OnNext(item)
				}
			}

			if e == nil {
				ob.OnDone()
			} else {
				ob.OnError(e)
			}
		}(ob)
	}

	ob := NewObserver()
	go func() {
		wg.Wait()
		ob.OnDone()
	}()

	return ob
}
