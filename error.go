package rxgo

import (
	"errors"
	"sync"
)

var (
	// An error thrown when an Observable or a sequence was queried but has no elements.
	ErrEmpty = errors.New("rxgo: empty value")
	// An error thrown when a value or values are missing from an observable sequence.
	ErrNotFound           = errors.New("rxgo: no values match")
	ErrSequence           = errors.New("rxgo: too many values match")
	ErrArgumentOutOfRange = errors.New("rxgo: argument out of range")
	// An error thrown by the timeout operator.
	ErrTimeout = errors.New("rxgo: timeout")
)

// Catches errors on the observable to be handled by returning a new observable or
// throwing an error.
func CatchError[T any](catch func(err error, caught Observable[T]) Observable[T]) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(1)

			var (
				upStream    = source.SubscribeOn(wg.Done)
				catchStream Subscriber[T]
			)

		observe:
			for {
				select {
				case <-subscriber.Closed():
					upStream.Stop()
					break observe

				case item, ok := <-upStream.ForEach():
					if !ok {
						break observe
					}

					if err := item.Err(); err != nil {
						wg.Add(1)
						catchStream = catch(err, source).SubscribeOn(wg.Done)
						break observe
					}

					item.Send(subscriber)
					if item.Done() {
						break observe
					}
				}
			}

			if catchStream != nil {
			catchLoop:
				for {
					select {
					case <-subscriber.Closed():
						catchStream.Stop()
						break catchLoop

					case item, ok := <-catchStream.ForEach():
						if !ok {
							break catchLoop
						}

						ended := item.Err() != nil || item.Done()
						item.Send(subscriber)
						if ended {
							break catchLoop
						}
					}
				}
			}

			wg.Wait()
		})
	}
}
