package observable

import (
	"sync"

	"github.com/reactivex/rxgo/handlers"
)

// transforms emitted items into observables and flattens them into single observable.
// maxInParallel argument controls how many transformed observables are processed in parallel
// For an example please take a look at flatmap_slice_test.go file in the examples directory.
func (o *observator) FlatMap(apply func(interface{}) Observable, maxInParallel uint) Observable {
	return o.flatMap(apply, maxInParallel, flatObservedSequence)
}

func (o *observator) flatMap(
	apply func(interface{}) Observable,
	maxInParallel uint,
	flatteningFunc func(out chan interface{}, o Observable, apply func(interface{}) Observable, maxInParallel uint)) Observable {

	out := make(chan interface{})

	if maxInParallel < 1 {
		maxInParallel = 1
	}

	go flatteningFunc(out, o, apply, maxInParallel)

	return &observator{
		ch: out,
	}
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

	for {
		element, err := o.Next()
		if err != nil {
			break
		}
		sequence = apply(element)
		count++
		wg.Add(1)
		go func() {
			defer wg.Done()
			sequence.Subscribe(emissionObserver).Block()
		}()

		if count%maxInParallel == 0 {
			wg.Wait()
		}
	}

	wg.Wait()
}

func newFlattenEmissionObserver(out chan interface{}) Observer {
	return NewObserver(handlers.NextFunc(func(element interface{}) {
		out <- element
	}))
}
