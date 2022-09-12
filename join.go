package rxgo

import (
	"errors"
	"sync"
	"sync/atomic"
)

// Create an observable that combines the latest values from all passed observables
// and the source into arrays and emits them.
func CombineLatestWith[T any](sources ...Observable[T]) OperatorFunc[T, []T] {
	return func(source Observable[T]) Observable[[]T] {
		sources = append([]Observable[T]{source}, sources...)
		return newObservable(func(subscriber Subscriber[[]T]) {
			var (
				noOfSource   = len(sources)
				emitCount    = new(atomic.Uint32)
				errCh        = make(chan error, 1)
				stopCh       = make(chan struct{})
				wg           = new(sync.WaitGroup)
				latestValues = make([]T, noOfSource)
			)

			wg.Add(noOfSource)

			// To ensure the output array always has the same length,
			// combineLatest will actually wait for all input Observables
			// to emit at least once, before it starts emitting results.
			onNext := func() {
				if emitCount.Load() == uint32(noOfSource) {
					Next(latestValues).Send(subscriber)
				}
			}

			observeStream := func(index int, upStream Subscriber[T]) {
				var (
					emitted bool
				)

			loop:
				for {
					select {
					case <-subscriber.Closed():
						upStream.Stop()
						break loop

					case <-stopCh:
						upStream.Stop()
						break loop

					case item, ok := <-upStream.ForEach():
						if !ok {
							break loop
						}

						if err := item.Err(); err != nil {
							errCh <- err
							break loop
						}

						if item.Done() {
							break loop
						}

						//  Passing an empty array will result in an Observable that completes immediately.
						if !emitted {
							emitCount.Add(1)
							emitted = true
						}
						latestValues[index] = item.Value()
						onNext()
					}
				}
			}

			go func() {
				select {
				case <-subscriber.Closed():
					return
				case _, ok := <-errCh:
					if !ok {
						return
					}
					close(stopCh)
				}
			}()

			for i, source := range sources {
				subscriber := source.SubscribeOn(wg.Done)
				go observeStream(i, subscriber)
			}

			wg.Wait()

			select {
			case <-errCh:
			default:
				// Close error channel gracefully
				close(errCh)
			}

			Complete[[]T]().Send(subscriber)
		})
	}
}

// Accepts an Array of ObservableInput or a dictionary Object of ObservableInput
// and returns an Observable that emits either an array of values in the exact same
// order as the passed array, or a dictionary of values in the same shape as the
// passed dictionary.
func ForkJoin[T any](sources ...Observable[T]) Observable[[]T] {
	return newObservable(func(subscriber Subscriber[[]T]) {
		var (
			noOfSource = len(sources)
		)
		// forkJoin is an operator that takes any number of input observables which can be
		// passed either as an array or a dictionary of input observables. If no input
		// observables are provided (e.g. an empty array is passed), then the resulting
		// stream will complete immediately.
		if noOfSource < 1 {
			Complete[[]T]().Send(subscriber)
			return
		}

		var (
			wg = new(sync.WaitGroup)
			// mu  = new(sync.Mutex)
			err = new(atomic.Pointer[error])
			// Single buffered channel will not accept more than one signal
			errCh         = make(chan error, 1)
			stopCh        = make(chan struct{})
			latestValues  = make([]T, noOfSource)
			subscriptions = make([]Subscriber[T], noOfSource)
		)

		wg.Add(noOfSource)

		go func() {
			select {
			case <-subscriber.Closed():
				return
			case v, ok := <-errCh:
				if !ok {
					return
				}
				err.Swap(&v)
				close(stopCh)
			}
		}()

		// In order for the resulting array to have the same length as the number of
		// input observables, whenever any of the given observables completes without
		// emitting any value, forkJoin will complete at that moment as well and it
		// will not emit anything either, even if it already has some last values
		// from other observables.
		onNext := func(index int, v T) {
			// mu.Lock()
			// defer mu.Unlock()
			latestValues[index] = v
		}

		observeStream := func(index int, upStream Subscriber[T]) {
		observe:
			for {
				select {
				case <-stopCh:
					upStream.Stop()
					break observe

				case <-subscriber.Closed():
					upStream.Stop()
					break observe

				case item, ok := <-upStream.ForEach():
					if !ok {
						break observe
					}

					// if one error, everything error
					if err := item.Err(); err != nil {
						errCh <- err
						break observe
					}

					if item.Done() {
						break observe
					}

					// forkJoin will wait for all passed observables to emit and complete
					// and then it will emit an array or an object with last values from
					// corresponding observables.
					onNext(index, item.Value())
				}
			}
		}

		for i, source := range sources {
			subscriber := source.SubscribeOn(wg.Done)
			subscriptions[i] = subscriber
			go observeStream(i, subscriber)
		}

		wg.Wait()

		// Remove dangling go-routine
		select {
		case <-errCh:
		default:
			// Close error channel gracefully
			close(errCh)
		}

		for _, sub := range subscriptions {
			sub.Stop()
		}

		if exception := err.Load(); exception != nil {
			if errors.Is(*exception, ErrEmpty) {
				Complete[[]T]().Send(subscriber)
				return
			}

			Error[[]T](*exception).Send(subscriber)
			return
		}

		Next(latestValues).Send(subscriber)
		Complete[[]T]().Send(subscriber)
	})
}

// Combines multiple Observables to create an Observable whose values are calculated
// from the values, in order, of each of its input Observables.
func Zip[T any](sources ...Observable[T]) Observable[[]T] {
	return newObservable(func(subscriber Subscriber[[]T]) {
		var (
			wg         = new(sync.WaitGroup)
			noOfSource = len(sources)
			observers  = make([]Subscriber[T], 0, noOfSource)
		)

		wg.Add(noOfSource)

		for _, source := range sources {
			observers = append(observers, source.SubscribeOn(wg.Done))
		}

		unsubscribeAll := func() {
			for _, obs := range observers {
				obs.Stop()
			}
		}

		var (
			result    = make([]T, noOfSource)
			completes = [2]uint{} // true | false
		)
	loop:
		for {
			for i, obs := range observers {
				select {
				case <-subscriber.Closed():
					unsubscribeAll()
					break loop

				case item, ok := <-obs.ForEach():
					if !ok {
						completes[1]++
					} else if ok || item.Done() {
						completes[0]++
					}

					if item != nil {
						result[i] = item.Value()
					}
				}
			}

			if completes[1] >= uint(noOfSource) {
				Complete[[]T]().Send(subscriber)
				break loop
			}

			Next(result).Send(subscriber)

			// Reset the values for next loop
			result = make([]T, noOfSource)
			completes = [2]uint{}
		}

		wg.Wait()
	})
}
