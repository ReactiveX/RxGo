package rxgo

import (
	"context"
	"log"
	"sync"
	"sync/atomic"

	"golang.org/x/sync/errgroup"
)

// Flattens an Observable-of-Observables by applying combineLatest when the Observable-of-Observables completes.
func CombineLatestAll[T any, R any](project func(values []T) R) OperatorFunc[Observable[T], R] {
	return func(source Observable[Observable[T]]) Observable[R] {
		return newObservable(func(subscriber Subscriber[R]) {
			var (
				mu       = new(sync.RWMutex)
				upStream = source.SubscribeOn()
				buffer   = make([]Observable[T], 0)
				err      error
			)

		outerLoop:
			for {
				select {
				case <-subscriber.Closed():
					upStream.Stop()
					return

				case item, ok := <-upStream.ForEach():
					if !ok {
						break outerLoop
					}

					if err = item.Err(); err != nil {
						break outerLoop
					}

					if item.Done() {
						break outerLoop
					}

					buffer = append(buffer, item.Value())
				}
			}

			if err != nil {
				Error[R](err).Send(subscriber)
				return
			}

			var (
				noOfBuffer   = len(buffer)
				emitCount    = new(atomic.Uint32)
				latestValues = make([]T, noOfBuffer)
			)

			// To ensure the output array always has the same length,
			// combineLatest will actually wait for all input Observables
			// to emit at least once, before it starts emitting results.
			onNext := func() {
				if emitCount.Load() == uint32(noOfBuffer) {
					mu.RLock()
					Next(project(latestValues)).Send(subscriber)
					mu.RUnlock()
				}
			}

			g, ctx := errgroup.WithContext(context.TODO())

			observeStream := func(ctx context.Context, index int, obs Observable[T]) func() error {
				return func() error {
					var (
						emitted  bool
						upStream = obs.SubscribeOn()
					)

				loop:
					for {
						select {
						case <-ctx.Done():
							upStream.Stop()
							break loop

						case <-subscriber.Closed():
							upStream.Stop()
							break loop

						case item, ok := <-upStream.ForEach():
							if !ok {
								break loop
							}

							if err := item.Err(); err != nil {
								return err
							}

							if item.Done() {
								break loop
							}

							//  Passing an empty array will result in an Observable that completes immediately.
							if !emitted {
								emitCount.Add(1)
								emitted = true
							}

							mu.Lock()
							latestValues[index] = item.Value()
							mu.Unlock()
							onNext()
						}
					}

					return nil
				}
			}

			for i, source := range buffer {
				g.Go(observeStream(ctx, i, source))
			}

			if err := g.Wait(); err != nil {
				Error[R](err).Send(subscriber)
				return
			}

			Complete[R]().Send(subscriber)
		})
	}
}

// Create an observable that combines the latest values from all passed observables and the source into arrays and emits them.
func CombineLatestWith[T any](sources ...Observable[T]) OperatorFunc[T, []T] {
	return func(source Observable[T]) Observable[[]T] {
		sources = append([]Observable[T]{source}, sources...)
		return newObservable(func(subscriber Subscriber[[]T]) {
			var (
				mu           = new(sync.RWMutex)
				noOfSource   = len(sources)
				emitCount    = new(atomic.Uint32)
				latestValues = make([]T, noOfSource)
			)

			// To ensure the output array always has the same length,
			// combineLatest will actually wait for all input Observables
			// to emit at least once, before it starts emitting results.
			onNext := func() {
				if emitCount.Load() == uint32(noOfSource) {
					mu.RLock()
					Next(latestValues).Send(subscriber)
					mu.RUnlock()
				}
			}

			observeStream := func(ctx context.Context, index int, obs Observable[T]) func() error {
				return func() error {
					var (
						emitted  bool
						upStream = obs.SubscribeOn()
					)

				loop:
					for {
						select {
						case <-ctx.Done():
							upStream.Stop()
							break loop

						case <-subscriber.Closed():
							upStream.Stop()
							break loop

						case item, ok := <-upStream.ForEach():
							if !ok {
								break loop
							}

							if err := item.Err(); err != nil {
								return err
							}

							if item.Done() {
								break loop
							}

							//  Passing an empty array will result in an Observable that completes immediately.
							if !emitted {
								emitCount.Add(1)
								emitted = true
							}

							mu.Lock()
							latestValues[index] = item.Value()
							mu.Unlock()
							onNext()
						}
					}

					return nil
				}
			}

			g, ctx := errgroup.WithContext(context.TODO())

			for i, source := range sources {
				g.Go(observeStream(ctx, i, source))
			}

			if err := g.Wait(); err != nil {
				Error[[]T](err).Send(subscriber)
				return
			}

			Complete[[]T]().Send(subscriber)
		})
	}
}

// Converts a higher-order Observable into a first-order Observable by  concatenating the inner Observables in order.
func ConcatAll[T any]() OperatorFunc[Observable[T], T] {
	return func(source Observable[Observable[T]]) Observable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(1)

			var (
				index      uint
				upStream   = source.SubscribeOn(wg.Done)
				downStream Subscriber[T]
			)

			unsubscribeAll := func() {
				upStream.Stop()
				if downStream != nil {
					downStream.Stop()
				}
			}

		outerLoop:
			for {
				select {
				case <-subscriber.Closed():
					upStream.Stop()
					break outerLoop

				case item, ok := <-upStream.ForEach():
					// If the upstream closed, we break
					if !ok {
						break outerLoop
					}

					if err := item.Err(); err != nil {
						Error[T](err).Send(subscriber)
						break outerLoop
					}

					if item.Done() {
						Complete[T]().Send(subscriber)
						break outerLoop
					}

					wg.Add(1)
					// we should wait the projection to complete
					downStream = item.Value().SubscribeOn(wg.Done)

				innerLoop:
					for {
						select {
						case <-subscriber.Closed():
							unsubscribeAll()
							break outerLoop

						case item, ok := <-downStream.ForEach():
							if !ok {
								break innerLoop
							}

							if item.Err() != nil {
								unsubscribeAll()
								item.Send(subscriber)
								break outerLoop
							}

							if item.Done() {
								break innerLoop
							}

							item.Send(subscriber)
						}
					}

					index++
				}
			}

			wg.Wait()
		})
	}
}

// Emits all of the values from the source observable, then, once it completes, subscribes to each observable source provided, one at a time, emitting all of their values, and not subscribing to the next one until it completes.
func ConcatWith[T any](sources ...Observable[T]) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		sources = append([]Observable[T]{source}, sources...)
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				wg  = new(sync.WaitGroup)
				err error
			)

		outerLoop:
			for len(sources) > 0 {
				wg.Add(1)
				firstSource := sources[0]
				upStream := firstSource.SubscribeOn(wg.Done)

			innerLoop:
				for {
					select {
					case <-subscriber.Closed():
						upStream.Stop()
						return

					case item, ok := <-upStream.ForEach():
						if !ok {
							break innerLoop
						}

						if item.Done() {
							// start another loop
							break innerLoop
						}

						if err = item.Err(); err != nil {
							break outerLoop
						}

						item.Send(subscriber)
					}
				}

				sources = sources[1:]
			}

			if err != nil {
				Error[T](err).Send(subscriber)
			} else {
				Complete[T]().Send(subscriber)
			}

			wg.Wait()
		})
	}
}

// Accepts an Array of ObservableInput or a dictionary Object of ObservableInput and returns an Observable that emits either an array of values in the exact same order as the passed array, or a dictionary of values in the same shape as the passed dictionary.
func ForkJoin[T any](sources ...Observable[T]) Observable[[]T] {
	return newObservable(func(subscriber Subscriber[[]T]) {
		var (
			noOfSource = len(sources)
		)

		// forkJoin is an operator that takes any number of input observables which can be passed either as an array or a dictionary of input observables. If no input observables are provided (e.g. an empty array is passed), then the resulting stream will complete immediately.
		if noOfSource < 1 {
			Complete[[]T]().Send(subscriber)
			return
		}

		var (
			emitCount    = new(atomic.Uint32)
			mu           = new(sync.RWMutex)
			g, ctx       = errgroup.WithContext(context.TODO())
			latestValues = make([]T, noOfSource)
		)

		// In order for the resulting array to have the same length as the number of input observables, whenever any of the given observables completes without emitting any value, forkJoin will complete at that moment as well and it will not emit anything either, even if it already has some last values from other observables.
		onNext := func(index int, v T) {
			mu.Lock()
			latestValues[index] = v
			mu.Unlock()
		}

		observeStream := func(ctx context.Context, index int, obs Observable[T]) func() error {
			return func() error {
				var (
					emitted  bool
					upStream = obs.SubscribeOn()
				)

			observe:
				for {
					select {
					case <-ctx.Done():
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
							return err
						}

						if item.Done() {
							break observe
						}

						if !emitted {
							emitCount.Add(1)
							emitted = true
						}

						// forkJoin will wait for all passed observables to emit and complete and then it will emit an array or an object with last values from corresponding observables.
						onNext(index, item.Value())
					}
				}

				return nil
			}
		}

		for i, source := range sources {
			g.Go(observeStream(ctx, i, source))
		}

		if err := g.Wait(); err != nil {
			Error[[]T](err).Send(subscriber)
			return
		}

		if emitCount.Load() == uint32(len(sources)) {
			mu.RLock()
			Next(latestValues).Send(subscriber)
			mu.RUnlock()
		}

		Complete[[]T]().Send(subscriber)
	})
}

// Converts a higher-order Observable into a first-order Observable by dropping inner Observables while the previous inner Observable has not yet completed.
func ExhaustAll[T any]() OperatorFunc[Observable[T], T] {
	return ExhaustMap(func(value Observable[T], _ uint) Observable[T] {
		return value
	})
}

// FIXME: Merge the values from all observables to a single observable result.
func MergeWith[T any](input Observable[T], inputs ...Observable[T]) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		inputs = append([]Observable[T]{source, input}, inputs...)
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				wg                  = new(sync.WaitGroup)
				mu                  = new(sync.RWMutex)
				activeSubCount      = new(atomic.Int32)
				noOfInputs          = len(inputs)
				activeSubscriptions = make([]Subscriber[T], noOfInputs)
				err                 = new(atomic.Pointer[error])
				stopCh              = make(chan struct{})
				errCh               = make(chan error, 1)
			)

			onError := func(err error) {
				mu.Lock()
				defer mu.Unlock()
				select {
				case errCh <- err:
				default:
				}
			}

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

			observeStream := func(index int, stream Subscriber[T]) {
				defer activeSubCount.Add(-1)

			observe:
				for {
					select {
					case <-subscriber.Closed():
						stream.Stop()
						break observe

					case <-stopCh:
						stream.Stop()
						break observe

					case item, ok := <-stream.ForEach():
						if !ok {
							break observe
						}

						if err := item.Err(); err != nil {
							onError(err)
							break observe
						}

						if item.Done() {
							break observe
						}

						item.Send(subscriber)
					}
				}
			}

			wg.Add(noOfInputs)
			// activeSubCount.Store(int32(noOfInputs))

			for i, input := range inputs {
				activeSubscriptions[i] = input.SubscribeOn(wg.Done)
				go observeStream(i, activeSubscriptions[i])
			}

			wg.Wait()

			// Remove dangling go-routine
			select {
			case <-errCh:
			default:
				mu.Lock()
				// Close error channel gracefully
				close(errCh)
				mu.Unlock()
			}

			// Stop all stream
			for _, sub := range activeSubscriptions {
				sub.Stop()
			}

			if exception := err.Load(); exception != nil {
				Error[T](*exception).Send(subscriber)
				return
			}

			Complete[T]().Send(subscriber)
		})
	}
}

// FIXME: Creates an Observable that mirrors the first source Observable to emit a
// next, error or complete notification from the combination of the Observable
// to which the operator is applied and supplied Observables.
func RaceWith[T any](sources ...Observable[T]) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		sources = append([]Observable[T]{source}, sources...)
		return newObservable(func(subscriber Subscriber[T]) {

		})
	}
}

// Converts a higher-order Observable into a first-order Observable producing values only from the most recent observable sequence
func SwitchAll[T any]() OperatorFunc[Observable[T], T] {
	return func(source Observable[Observable[T]]) Observable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(1)

			var (
				upStream = source.SubscribeOn(wg.Done)
				// observables = make([]Observable[T], 0)
			)

		outerLoop:
			for {
				select {
				case <-subscriber.Closed():
					upStream.Stop()
					break outerLoop

				case item, ok := <-upStream.ForEach():
					wg.Add(1)
					item.Value().SubscribeOn(wg.Done)
					log.Println(item, ok)
				}
			}
		})
	}
}

// Collects all observable inner sources from the source, once the source completes, it will subscribe to all inner sources, combining their values by index and emitting them.
func ZipAll[T any]() OperatorFunc[Observable[T], []T] {
	return func(source Observable[Observable[T]]) Observable[[]T] {
		return newObservable(func(subscriber Subscriber[[]T]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(1)

			var (
				runNext     bool
				upStream    = source.SubscribeOn(wg.Done)
				observables = make([]Observable[T], 0)
			)

			// Collects all observable inner sources from the source, once the
			// source completes, it will subscribe to all inner sources,
			// combining their values by index and emitting them.
		loop:
			for {
				select {
				case <-subscriber.Closed():
					upStream.Stop()
					break loop

				case item, ok := <-upStream.ForEach():
					if !ok {
						break loop
					}

					if err := item.Err(); err != nil {
						Error[[]T](err).Send(subscriber)
						break loop
					}

					if item.Done() {
						runNext = true
						break loop
					}

					observables = append(observables, item.Value())
				}
			}

			var noOfObservables = len(observables)
			if runNext && noOfObservables > 0 {
				var (
					observers = make([]Subscriber[T], 0, noOfObservables)
					result    []T
					completed uint
				)

				setupValues := func() {
					result = make([]T, noOfObservables)
					completed = 0
				}

				unsubscribeAll := func() {
					for _, obs := range observers {
						obs.Stop()
					}
				}

				setupValues()
				wg.Add(noOfObservables)
				for _, obs := range observables {
					observers = append(observers, obs.SubscribeOn(wg.Done))
				}

			outerLoop:
				for {
				innerLoop:
					for i, obs := range observers {
						select {
						case <-subscriber.Closed():
							unsubscribeAll()
							break outerLoop

						case item, ok := <-obs.ForEach():
							if !ok || item.Done() {
								completed++
								unsubscribeAll()
								break innerLoop
							}

							if err := item.Err(); err != nil {
								unsubscribeAll()
								Error[[]T](err).Send(subscriber)
								break outerLoop
							}

							if item != nil {
								result[i] = item.Value()
							}
						}
					}

					// any of the stream completed, we will escape
					if completed > 0 {
						Complete[[]T]().Send(subscriber)
						break outerLoop
					}

					Next(result).Send(subscriber)

					// reset the values for next loop
					setupValues()
				}
			}

			wg.Wait()
		})
	}
}

// Combines multiple Observables to create an Observable whose values are calculated from the values, in order, of each of its input Observables.
func ZipWith[T any](input Observable[T], inputs ...Observable[T]) OperatorFunc[T, []T] {
	return func(source Observable[T]) Observable[[]T] {
		inputs = append([]Observable[T]{source, input}, inputs...)
		return newObservable(func(subscriber Subscriber[[]T]) {
			var (
				wg         = new(sync.WaitGroup)
				noOfSource = uint(len(inputs))
				observers  = make([]Subscriber[T], 0, noOfSource)
			)

			wg.Add(int(noOfSource))

			for _, input := range inputs {
				observers = append(observers, input.SubscribeOn(wg.Done))
			}

			unsubscribeAll := func() {
				for _, obs := range observers {
					obs.Stop()
				}
			}

			var (
				result    []T
				completed uint
			)

			setupValues := func() {
				result = make([]T, noOfSource)
				completed = 0
			}

			setupValues()

		outerLoop:
			for {

			innerLoop:
				for i, obs := range observers {
					select {
					case <-subscriber.Closed():
						unsubscribeAll()
						break outerLoop

					case item, ok := <-obs.ForEach():
						if !ok || item.Done() {
							completed++
							unsubscribeAll()
							break innerLoop
						}

						if err := item.Err(); err != nil {
							unsubscribeAll()
							Error[[]T](err).Send(subscriber)
							break outerLoop
						}

						if item != nil {
							result[i] = item.Value()
						}
					}
				}

				// any of the stream completed, we will escape
				if completed > 0 {
					Complete[[]T]().Send(subscriber)
					break outerLoop
				}

				Next(result).Send(subscriber)

				// reset the values for next loop
				setupValues()
			}

			wg.Wait()
		})
	}
}
