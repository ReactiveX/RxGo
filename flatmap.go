package rxgo

import (
	"context"

	"github.com/reactivex/rxgo/handlers"
	"golang.org/x/sync/semaphore"
)

// transforms emitted items into observables and flattens them into single observable.
// maxInParallel argument controls how many transformed observables are processed in parallel
// For an example please take a look at flatmap_slice_test.go file in the examples directory.
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

	return &observable{
		ch: out,
	}
}

func flatObservedSequence(out chan interface{}, o Observable, apply func(interface{}) Observable, maxInParallel uint) {
	ctx := context.TODO()
	sem := semaphore.NewWeighted(int64(maxInParallel))

	defer close(out)
	emissionObserver := newFlattenEmissionObserver(out)

	for {
		element, err := o.Next()
		if err != nil {
			break
		}
		sequence := apply(element)
		sem.Acquire(ctx, 1)
		go func() {
			defer sem.Release(1)
			sequence.Subscribe(emissionObserver).Block()
		}()
	}

	sem.Acquire(ctx, int64(maxInParallel))
}

func newFlattenEmissionObserver(out chan interface{}) Observer {
	return NewObserver(handlers.NextFunc(func(element interface{}) {
		out <- element
	}))
}
