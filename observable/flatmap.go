package observable

import (
	"sync"

	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/observer"
)

// transforms emitted items into observables and flattens them into single observable.
// maxInParallel argument controls how many transformed observables are processed in parallel
// For an example please take a look at flatmap_slice_test.go file in the examples directory.
func (o Observable) FlatMap(apply func(interface{}) Observable, maxInParallel uint) Observable {
	return o.flatMap(apply, maxInParallel, flatObservedSequence)
}

func (o Observable) flatMap(
	apply func(interface{}) Observable,
	maxInParallel uint,
	flatteningFunc func(out chan interface{}, o Observable, apply func(interface{}) Observable, maxInParallel uint)) Observable {

	out := make(chan interface{})

	if maxInParallel < 1 {
		maxInParallel = 1
	}

	go flatteningFunc(out, o, apply, maxInParallel)

	return Observable(out)
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
	for element := range o {
		sequence = apply(element)
		count++
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-(sequence.Subscribe(*emissionObserver))
		}()

		if count%maxInParallel == 0 {
			wg.Wait()
		}
	}

	wg.Wait()
}

func newFlattenEmissionObserver(out chan interface{}) *observer.Observer {
	ob := observer.New(handlers.NextFunc(func(element interface{}) {
		out <- element
	}))
	return &ob
}
