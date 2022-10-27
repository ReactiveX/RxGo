package rxgo

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"
)

// Buffers the source Observable values until closingNotifier emits.
func Buffer[T any, R any](closingNotifier Observable[R]) OperatorFunc[T, []T] {
	return func(source Observable[T]) Observable[[]T] {
		return newObservable(func(subscriber Subscriber[[]T]) {
			var (
				wg          = new(sync.WaitGroup)
				mu          = new(sync.RWMutex)
				errOnce     = new(atomic.Pointer[error])
				ctx, cancel = context.WithCancel(context.TODO())
				buffer      = make([]T, 0)
			)

			wg.Add(2)

			onError := func(err error) {
				errOnce.CompareAndSwap(nil, &err)
				cancel()
			}

			observeStream := func(stream Subscriber[R]) {
			innerLoop:
				for {
					select {
					case <-ctx.Done():
						stream.Stop()
						break innerLoop

					case <-subscriber.Closed():
						stream.Stop()
						break innerLoop

					case <-stream.Closed():
						break innerLoop

					case item, ok := <-stream.ForEach():
						if !ok {
							break innerLoop
						}

						if err := item.Err(); err != nil {
							onError(err)
							break innerLoop
						}

						if item.Done() {
							break innerLoop
						}

						Next(buffer).Send(subscriber)

						mu.Lock()
						// reset buffer
						buffer = make([]T, 0)
						mu.Unlock()
					}
				}
			}

			var (
				upStream     = source.SubscribeOn(wg.Done)
				notifyStream = closingNotifier.SubscribeOn(wg.Done)
			)

			go observeStream(notifyStream)

		outerLoop:
			for {
				select {
				case <-ctx.Done():
					upStream.Stop()
					notifyStream.Stop()
					break outerLoop

				case <-subscriber.Closed():
					upStream.Stop()
					notifyStream.Stop()
					break outerLoop

				case item, ok := <-upStream.ForEach():
					if !ok {
						notifyStream.Stop()
						break outerLoop
					}

					if err := item.Err(); err != nil {
						notifyStream.Stop()
						onError(err)
						break outerLoop
					}

					if item.Done() {
						notifyStream.Stop()
						Complete[[]T]().Send(subscriber)
						break outerLoop
					}

					mu.Lock()
					buffer = append(buffer, item.Value())
					mu.Unlock()
				}
			}

			wg.Wait()

			if err := errOnce.Load(); err != nil {
				Error[[]T](*err).Send(subscriber)
				return
			}
		})
	}
}

// Buffers the source Observable values until the size hits the maximum bufferSize given.
func BufferCount[T any](bufferSize uint, startBufferEvery ...uint) OperatorFunc[T, []T] {
	offset := bufferSize
	if len(startBufferEvery) > 0 {
		offset = startBufferEvery[0]
	}
	return func(source Observable[T]) Observable[[]T] {
		return newObservable(func(subscriber Subscriber[[]T]) {
			var (
				wg     = new(sync.WaitGroup)
				buffer = make([]T, 0, bufferSize)
			)

			wg.Add(1)

			var (
				upStream = source.SubscribeOn(wg.Done)
			)

			unshiftBuffer := func() {
				if offset > uint(len(buffer)) {
					buffer = make([]T, 0, bufferSize)
					return
				}
				buffer = buffer[offset:]
			}

		outerLoop:
			for {
				select {
				case <-subscriber.Closed():
					upStream.Stop()
					break outerLoop

				case item, ok := <-upStream.ForEach():
					if !ok {
						break outerLoop
					}

					if err := item.Err(); err != nil {
						Error[[]T](err).Send(subscriber)
						break outerLoop
					}

					if item.Done() {
						// flush remaining buffer
						for len(buffer) > 0 {
							Next(buffer).Send(subscriber)
							unshiftBuffer()
						}
						Complete[[]T]().Send(subscriber)
						break outerLoop
					}

					buffer = append(buffer, item.Value())
					if uint(len(buffer)) >= bufferSize {
						Next(buffer).Send(subscriber)
						unshiftBuffer()
					}
				}
			}

			wg.Wait()
		})
	}
}

// Buffers the source Observable values for a specific time period.
func BufferTime[T any](bufferTimeSpan time.Duration) OperatorFunc[T, []T] {
	return func(source Observable[T]) Observable[[]T] {
		return newObservable(func(subscriber Subscriber[[]T]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(1)

			var (
				buffer   []T
				upStream = source.SubscribeOn(wg.Done)
				timer    *time.Timer
			)

			stopTimer := func() {
				if timer != nil {
					timer.Stop()
				}
			}

			setValues := func() {
				buffer = make([]T, 0)
				stopTimer()
				timer = time.NewTimer(bufferTimeSpan)
			}

			setValues()

		observe:
			for {
				select {
				case <-subscriber.Closed():
					stopTimer()
					upStream.Stop()
					break observe

				case <-timer.C:
					Next(buffer).Send(subscriber)
					setValues()

				case item, ok := <-upStream.ForEach():
					if !ok {
						break observe
					}

					if err := item.Err(); err != nil {
						Error[[]T](err).Send(subscriber)
						break observe
					}

					if item.Done() {
						Next(buffer).Send(subscriber)
						Complete[[]T]().Send(subscriber)
						break observe
					}

					buffer = append(buffer, item.Value())
				}
			}

			wg.Wait()

			// prevent memory leak
			upStream.Stop()
			stopTimer()
		})
	}
}

// Buffers the source Observable values starting from an emission from openings and ending when the output of closingSelector emits.
func BufferToggle[T any, O any](openings Observable[O], closingSelector func(value O) Observable[O]) OperatorFunc[T, []T] {
	return func(source Observable[T]) Observable[[]T] {
		return newObservable(func(subscriber Subscriber[[]T]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(2)

			var (
				allowed     bool
				buffer      []T
				startStream = openings.SubscribeOn(wg.Done)
				upStream    = source.SubscribeOn(wg.Done)
				emitStream  Subscriber[O]
				stopCh      <-chan Notification[O]
			)

			setupValues := func() {
				allowed = false
				buffer = make([]T, 0)
				stopCh = make(<-chan Notification[O])
			}

			unsubscribeAll := func() {
				startStream.Stop()
				if emitStream != nil {
					emitStream.Stop()
				}
			}

			setupValues()

			// buffers values from the source by opening the buffer via signals from an Observable provided to openings, and closing and sending the buffers when a Subscribable or Promise returned by the closingSelector function emits.
		observe:
			for {
				select {
				case <-subscriber.Closed():
					break observe

				case item, ok := <-startStream.ForEach():
					if !ok {
						break observe
					}

					allowed = true
					if emitStream != nil {
						// unsubscribe the previous one
						emitStream.Stop()
					}
					wg.Add(1)
					emitStream = closingSelector(item.Value()).SubscribeOn(wg.Done)
					stopCh = emitStream.ForEach()

				case <-stopCh:
					Next(buffer).Send(subscriber)
					if emitStream != nil {
						emitStream.Stop()
					}
					setupValues()

				case item, ok := <-upStream.ForEach():
					if !ok {
						break observe
					}

					if err := item.Err(); err != nil {
						unsubscribeAll()
						Error[[]T](err).Send(subscriber)
						break observe
					}

					if item.Done() {
						unsubscribeAll()
						Complete[[]T]().Send(subscriber)
						break observe
					}

					if allowed {
						buffer = append(buffer, item.Value())
					}
				}
			}

			wg.Wait()
		})
	}
}

// Buffers the source Observable values, using a factory function of closing Observables to determine when to close, emit, and reset the buffer.
func BufferWhen[T any, R any](closingSelector func() Observable[R]) OperatorFunc[T, []T] {
	return func(source Observable[T]) Observable[[]T] {
		return newObservable(func(subscriber Subscriber[[]T]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(2)

			var (
				index         uint
				buffer        = make([]T, 0)
				upStream      = source.SubscribeOn(wg.Done)
				closingStream = closingSelector().SubscribeOn(wg.Done)
			)

			unsubscribeAll := func() {
				upStream.Stop()
				closingStream.Stop()
			}

			onError := func(err error) {
				unsubscribeAll()
				Error[[]T](err).Send(subscriber)
			}

			onComplete := func() {
				unsubscribeAll()
				if len(buffer) > 0 {
					Next(buffer).Send(subscriber)
				}
				Complete[[]T]().Send(subscriber)
			}

		observe:
			for {
				select {
				case <-subscriber.Closed():
					unsubscribeAll()
					break observe

				case item, ok := <-upStream.ForEach():
					// if the upstream closed, we break
					if !ok {
						break observe
					}

					if err := item.Err(); err != nil {
						onError(err)
						break observe
					}

					if item.Done() {
						onComplete()
						break observe
					}

					buffer = append(buffer, item.Value())
					index++

				case item, ok := <-closingStream.ForEach():
					if !ok {
						break observe
					}

					if err := item.Err(); err != nil {
						onError(err)
						break observe
					}

					if item.Done() {
						onComplete()
						break observe
					}

					Next(buffer).Send(subscriber)
					// reset buffer values after sent
					buffer = make([]T, 0)
				}
			}

			wg.Wait()

			unsubscribeAll()
		})
	}
}

// Projects each source value to an Observable which is merged in the output Observable, in a serialized fashion waiting for each one to complete before merging the next.
func ConcatMap[T any, R any](project ProjectionFunc[T, R]) OperatorFunc[T, R] {
	return func(source Observable[T]) Observable[R] {
		return newObservable(func(subscriber Subscriber[R]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(1)

			var (
				index       uint
				upStream    = source.SubscribeOn(wg.Done)
				innerStream Subscriber[R]
			)

		outerLoop:
			for {
				select {
				case <-subscriber.Closed():
					upStream.Stop()
					break outerLoop

				case item, ok := <-upStream.ForEach():
					// if the upstream closed, we break
					if !ok {
						break outerLoop
					}

					if err := item.Err(); err != nil {
						Error[R](err).Send(subscriber)
						break outerLoop
					}

					if item.Done() {
						Complete[R]().Send(subscriber)
						break outerLoop
					}

					wg.Add(1)
					// we should wait the projection to complete
					innerStream = project(item.Value(), index).SubscribeOn(wg.Done)

				innerLoop:
					for {
						select {
						case <-subscriber.Closed():
							upStream.Stop()
							innerStream.Stop()
							break outerLoop

						case item, ok := <-innerStream.ForEach():
							if !ok {
								upStream.Stop()
								break innerLoop
							}

							if err := item.Err(); err != nil {
								upStream.Stop()
								innerStream.Stop()
								item.Send(subscriber)
								break outerLoop
							}

							if item.Done() {
								innerStream.Stop()
								break innerLoop
							}

							item.Send(subscriber)
						}
					}

					innerStream.Stop()
					index++
				}
			}

			wg.Wait()
		})
	}
}

// Projects each source value to an Observable which is merged in the output Observable only if the previous projected Observable has completed.
func ExhaustMap[T any, R any](project ProjectionFunc[T, R]) OperatorFunc[T, R] {
	if project == nil {
		panic(`rxgo: "ExhaustMap" expected project func`)
	}
	return func(source Observable[T]) Observable[R] {
		return newObservable(func(subscriber Subscriber[R]) {
			var (
				index    uint
				err      error
				allowed  = new(atomic.Pointer[bool])
				upStream = source.SubscribeOn()
				g, ctx   = errgroup.WithContext(context.TODO())
			)

			flag := true
			allowed.Store(&flag)

			observeStream := func(ctx context.Context, index uint, value T) func() error {
				return func() error {
					var (
						stream = project(value, index).SubscribeOn()
					)

				innerLoop:
					for {
						select {
						case <-ctx.Done():
							return nil

						case <-subscriber.Closed():
							stream.Stop()
							return nil

						case item, ok := <-stream.ForEach():
							if !ok {
								break innerLoop
							}

							if err := item.Err(); err != nil {
								return err
							}

							if item.Done() {
								flag := true
								allowed.Store(&flag)
								return nil
							}

							item.Send(subscriber)
						}
					}

					return nil
				}
			}

		outerLoop:
			for {
				select {
				case <-subscriber.Closed():
					upStream.Stop()
					break outerLoop

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

					if *allowed.Load() {
						flag := false
						allowed.Store(&flag)
						g.Go(observeStream(ctx, index, item.Value()))
						index++
					}
				}
			}

			if err != nil {
				Error[R](err).Send(subscriber)
				return
			}

			if err := g.Wait(); err != nil {
				Error[R](err).Send(subscriber)
				return
			}

			Complete[R]().Send(subscriber)
		})
	}
}

// Recursively projects each source value to an Observable which is merged in the output Observable.
func Expand[T any, R any](project ProjectionFunc[T, R]) OperatorFunc[T, Either[T, R]] {
	return func(source Observable[T]) Observable[Either[T, R]] {
		return newObservable(func(subscriber Subscriber[Either[T, R]]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(1)

			var (
				index    uint
				streams  []Observable[R]
				upStream = source.SubscribeOn(wg.Done)
			)

			reset := func() {
				streams = make([]Observable[R], 0)
			}

			reset()

		innerLoop:
			for {
				select {
				case <-subscriber.Closed():
					upStream.Stop()
					reset()
					break innerLoop

				case item, ok := <-upStream.ForEach():
					if !ok {
						break innerLoop
					}

					if err := item.Err(); err != nil {
						Error[Either[T, R]](err).Send(subscriber)
						reset()
						break innerLoop
					}

					if item.Done() {
						break innerLoop
					}

					Next(Left[T, R](item.Value())).Send(subscriber)
					streams = append(streams, project(item.Value(), index))
					index++
				}
			}

			for len(streams) > 0 {
				log.Println(streams[0])
				wg.Add(1)
				stream := streams[0].SubscribeOn(wg.Done)

			outerLoop:
				for {
					select {
					case <-subscriber.Closed():
						break outerLoop

					case item, ok := <-stream.ForEach():
						if !ok {
							break outerLoop
						}

						log.Println(item)
					}
				}

				streams = streams[1:]
			}

			log.Println(streams)

			wg.Wait()
		})
	}
}

// Groups the items emitted by an Observable according to a specified criterion, and emits these grouped items as GroupedObservables, one GroupedObservable per group.
// FIXME: maybe we should have a buffer channel
func GroupBy[T any, K comparable](keySelector func(value T) K) OperatorFunc[T, GroupedObservable[K, T]] {
	if keySelector == nil {
		panic(`rxgo: "GroupBy" expected keySelector func`)
	}
	return func(source Observable[T]) Observable[GroupedObservable[K, T]] {
		return newObservable(func(subscriber Subscriber[GroupedObservable[K, T]]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(1)

			var (
				key      K
				upStream = source.SubscribeOn(wg.Done)
				keySet   = make(map[K]Subject[T])
			)

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
						break loop
					}

					if item.Done() {
						for k, kv := range keySet {
							Next(NewGroupedObservable(k, func() Subject[T] {
								return kv
							})).Send(subscriber)
						}
						Complete[GroupedObservable[K, T]]().Send(subscriber)
						break loop
					}

					key = keySelector(item.Value())
					if _, exists := keySet[key]; !exists {
						keySet[key] = NewSubscriber[T](10000)
					}

					item.Send(keySet[key])
				}
			}

			wg.Wait()
		})
	}
}

// Map transforms the items emitted by an Observable by applying a function to each item.
func Map[T any, R any](mapper func(T, uint) (R, error)) OperatorFunc[T, R] {
	if mapper == nil {
		panic(`rxgo: "Map" expected mapper func`)
	}
	return func(source Observable[T]) Observable[R] {
		var (
			index uint
		)
		return createOperatorFunc(
			source,
			func(obs Observer[R], v T) {
				output, err := mapper(v, index)
				index++
				if err != nil {
					obs.Error(err)
					return
				}
				obs.Next(output)
			},
			func(obs Observer[R], err error) {
				obs.Error(err)
			},
			func(obs Observer[R]) {
				obs.Complete()
			},
		)
	}
}

// Projects each source value to an Observable which is merged in the output Observable.
func MergeMap[T any, R any](project ProjectionFunc[T, R], concurrent ...uint) OperatorFunc[T, R] {
	// TODO: support concurrent
	return func(source Observable[T]) Observable[R] {
		return newObservable(func(subscriber Subscriber[R]) {
			var (
				errOnce     = new(atomic.Pointer[error])
				wg          = new(sync.WaitGroup)
				ctx, cancel = context.WithCancel(context.TODO())
			)

			wg.Add(1)

			var (
				index    uint
				upStream = source.SubscribeOn(wg.Done)
			)

			onError := func(err error) {
				errOnce.CompareAndSwap(nil, &err)
				cancel()
			}

			observeStream := func(stream Subscriber[R]) {
			innerLoop:
				for {
					select {
					case <-ctx.Done():
						stream.Stop()
						break innerLoop

					case <-stream.Closed():
						stream.Stop()
						break innerLoop

					case item, ok := <-stream.ForEach():
						if !ok {
							break innerLoop
						}

						if err := item.Err(); err != nil {
							onError(err)
							break innerLoop
						}

						if item.Done() {
							break innerLoop
						}

						item.Send(subscriber)
					}
				}
			}

		outerLoop:
			for {
				select {
				case <-subscriber.Closed():
					upStream.Stop()
					cancel()
					break outerLoop

				case item, ok := <-upStream.ForEach():
					if !ok {
						break outerLoop
					}

					if err := item.Err(); err != nil {
						onError(err)
						break outerLoop
					}

					if item.Done() {
						// even upstream is done, we need to wait downstream done as well
						break outerLoop
					}

					wg.Add(1)
					subscription := project(item.Value(), index).SubscribeOn(wg.Done)
					go observeStream(subscription)
					index++
				}
			}

			wg.Wait()

			if err := errOnce.Load(); err != nil {
				Error[R](*err).Send(subscriber)
				return
			}

			Complete[R]().Send(subscriber)
		})
	}
}

// Applies an accumulator function over the source Observable where the accumulator function itself returns an Observable, then each intermediate Observable returned is merged into the output Observable.
func MergeScan[V any, A any](accumulator func(acc A, value V, index uint) Observable[A], seed A, concurrent ...uint) OperatorFunc[V, A] {
	return func(source Observable[V]) Observable[A] {
		return newObservable(func(subscriber Subscriber[A]) {
			var (
				exception   error
				errOnce     = new(sync.Once)
				wg          = new(sync.WaitGroup)
				ctx, cancel = context.WithCancel(context.TODO())
			)

			wg.Add(1)

			var (
				index      uint
				finalValue = new(atomic.Pointer[A])
				upStream   = source.SubscribeOn(wg.Done)
			)

			onError := func(err error) {
				errOnce.Do(func() {
					exception = err
					cancel()
				})
			}

			finalValue.Store(&seed)

			// MergeScan internally keeps the value of the acc parameter: as long as the source Observable emits without inner Observable emitting, the acc will be set to seed.

			observeStream := func(stream Subscriber[A]) {
				var (
					value A
				)

			innerLoop:
				for {
					select {
					case <-ctx.Done():
						stream.Stop()
						cancel()
						break innerLoop

					case <-subscriber.Closed():
						stream.Stop()
						cancel()
						break innerLoop

					case item, ok := <-stream.ForEach():
						if !ok {
							break innerLoop
						}

						if err := item.Err(); err != nil {
							onError(err)
							break innerLoop
						}

						if item.Done() {
							break innerLoop
						}

						value = item.Value()
						finalValue.Store(&value)
						item.Send(subscriber)
					}
				}
			}

		outerLoop:
			for {
				select {
				case <-ctx.Done():
					upStream.Stop()
					cancel()
					break outerLoop

				case <-subscriber.Closed():
					upStream.Stop()
					cancel()
					break outerLoop

				case item, ok := <-upStream.ForEach():
					if !ok {
						break outerLoop
					}

					if err := item.Err(); err != nil {
						onError(err)
						break outerLoop
					}

					if item.Done() {
						break outerLoop
					}

					wg.Add(1)
					stream := accumulator(*finalValue.Load(), item.Value(), index).SubscribeOn(wg.Done)
					go observeStream(stream)
					index++
				}
			}

			wg.Wait()

			if exception != nil {
				Error[A](exception).Send(subscriber)
				return
			}

			Complete[A]().Send(subscriber)
		})
	}
}

// Groups pairs of consecutive emissions together and emits them as an array of two values.
func PairWise[T any]() OperatorFunc[T, Tuple[T, T]] {
	return func(source Observable[T]) Observable[Tuple[T, T]] {
		var (
			result     = make([]T, 0, 2)
			noOfRecord int
		)
		return createOperatorFunc(
			source,
			func(obs Observer[Tuple[T, T]], v T) {
				result = append(result, v)
				noOfRecord = len(result)
				if noOfRecord >= 2 {
					obs.Next(NewTuple(result[0], result[1]))
					result = result[1:]
				}
			},
			func(obs Observer[Tuple[T, T]], err error) {
				obs.Error(err)
			},
			func(obs Observer[Tuple[T, T]]) {
				obs.Complete()
			},
		)
	}
}

// Useful for encapsulating and managing state. Applies an accumulator (or "reducer function") to each value from the source after an initial state is established -- either via a seed value (second argument), or from the first value from the source.
func Scan[V any, A any](accumulator AccumulatorFunc[A, V], seed A) OperatorFunc[V, A] {
	if accumulator == nil {
		panic(`rxgo: "Scan" expected accumulator func`)
	}
	return func(source Observable[V]) Observable[A] {
		var (
			index  uint
			result = seed
			err    error
		)
		return createOperatorFunc(
			source,
			func(obs Observer[A], v V) {
				result, err = accumulator(result, v, index)
				if err != nil {
					obs.Error(err)
					return
				}
				obs.Next(result)
				index++
			},
			func(obs Observer[A], err error) {
				obs.Error(err)
			},
			func(obs Observer[A]) {
				obs.Complete()
			},
		)
	}
}

// Applies an accumulator function over the source Observable where the accumulator function itself returns an Observable, emitting values only from the most recently returned Observable.
func SwitchScan() {

}

// Projects each source value to an Observable which is merged in the output Observable, emitting values only from the most recently projected Observable.
func SwitchMap[T any, R any](project func(value T, index uint) Observable[R]) OperatorFunc[T, R] {
	return func(source Observable[T]) Observable[R] {
		return newObservable(func(subscriber Subscriber[R]) {
			var (
				wg           = new(sync.WaitGroup)
				errOnce      = new(atomic.Pointer[error])
				completeOnce = new(sync.Once)
				ctx, cancel  = context.WithCancel(context.TODO())
				index        uint
			)

			wg.Add(1)

			var (
				upStream   = source.SubscribeOn(wg.Done)
				downStream Subscriber[R]
			)

			onError := func(err error) {
				errOnce.CompareAndSwap(nil, &err)
				cancel()
			}

			onComplete := func() {
				completeOnce.Do(func() {
					Complete[R]().Send(subscriber)
				})
			}

			observeStream := func(stream Subscriber[R]) {
			innerLoop:
				for {
					select {
					case <-ctx.Done():
						stream.Stop()
						break innerLoop

					case <-subscriber.Closed():
						stream.Stop()
						break innerLoop

					case <-stream.Closed():
						stream.Stop()
						break innerLoop

					case item, ok := <-stream.ForEach():
						if !ok {
							break innerLoop
						}

						if err := item.Err(); err != nil {
							onError(err)
							break innerLoop
						}

						if item.Done() {
							break innerLoop
						}

						item.Send(subscriber)
					}
				}
			}

		outerLoop:
			for {
				select {
				case <-ctx.Done():
					upStream.Stop()
					break outerLoop

				case <-subscriber.Closed():
					upStream.Stop()
					cancel()
					break outerLoop

				case item, ok := <-upStream.ForEach():
					if !ok {
						break outerLoop
					}

					// if the previous stream are still being processed while a new change is already made, it will cancel the previous subscription and start a new subscription on the latest change.

					if err := item.Err(); err != nil {
						onError(err)
						break outerLoop
					}

					if item.Done() {
						if downStream == nil {
							onComplete()
						}
						break outerLoop
					}

					// stop the previous stream
					if downStream != nil {
						downStream.Stop()
					}

					wg.Add(1)
					downStream = project(item.Value(), index).SubscribeOn(wg.Done)
					go observeStream(downStream)
					index++
				}
			}

			wg.Wait()

			if err := errOnce.Load(); err != nil {
				Error[R](*err).Send(subscriber)
				return
			}

			onComplete()
		})
	}
}
