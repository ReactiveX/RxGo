package rxgo

import (
	"sync"

	"github.com/reactivex/rxgo/handlers"
)

// Transforms emitted items into Observables and flattens them into a single Observable.
// The maxInParallel argument controls how many transformed Observables are processed in parallel.
// For an example, please take a look at flatmap_slice_test.go in the examples directory.
func (o *observable) FlatMap(apply func(interface{}) Observable, maxInParallel uint) Observable {
	return o.flatMap(apply, maxInParallel, flatObservedSequence)
}

func (o *observable) flatMap(
	apply func(interface{}) Observable,
	maxInParallel uint,
	flatteningFunc func(out chan interface{}, o Observable, apply func(interface{}) Observable, maxInParallel uint)) Observable {

	out := make(chan interface{})

	if maxInParallel < 1 {
		maxInParallel = 1
	}

	go flatteningFunc(out, o, apply, maxInParallel)

	return newObservableFromChannel(out)
}

func flatObservedSequence(out chan interface{}, o Observable, apply func(interface{}) Observable, maxInParallel uint) {
	var (
		sequence Observable
		wg       sync.WaitGroup
		count    uint
	)

	defer close(out)
	emissionObserver := newFlattenEmissionObserver(out)

	count = 0

	it := o.Iterator()
	for {
		if item, err := it.Next(); err == nil {
			sequence = apply(item)
			count++
			wg.Add(1)
			go func(s Observable) {
				defer wg.Done()
				s.Subscribe(emissionObserver).Block()
			}(sequence)

			if count%maxInParallel == 0 {
				wg.Wait()
			}
		} else {
			break
		}
	}

	wg.Wait()
}

func newFlattenEmissionObserver(out chan interface{}) Observer {
	return NewObserver(handlers.NextFunc(func(element interface{}) {
		out <- element
	}))
}
