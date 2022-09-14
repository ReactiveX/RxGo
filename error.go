package rxgo

import (
	"errors"
	"sync"

	"golang.org/x/exp/constraints"
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

type RetryConfig struct {
	Count          uint
	ResetOnSuccess bool
}

type retryConfig interface {
	constraints.Unsigned | RetryConfig
}

// Returns an Observable that mirrors the source Observable with the exception of an error.
func Retry[T any, C retryConfig](config ...C) OperatorFunc[T, T] {
	var maxRetryCount uint
	switch v := any(config[0]).(type) {
	case RetryConfig:
	case uint8:
		maxRetryCount = uint(v)
	case uint16:
		maxRetryCount = uint(v)
	case uint32:
		maxRetryCount = uint(v)
	case uint64:
		maxRetryCount = uint(v)
	case uint:
		maxRetryCount = v
	}
	return func(source Observable[T]) Observable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(1)

			var (
				upStream = source.SubscribeOn(wg.Done)
				errCount uint
			)

		observe:
			//  If count is omitted, retry will try to resubscribe on errors infinite number of times.
			for maxRetryCount == 0 || errCount <= maxRetryCount {
				select {
				case <-subscriber.Closed():
					upStream.Stop()
					break observe

				case item, ok := <-upStream.ForEach():
					if !ok {
						break observe
					}

					if err := item.Err(); err != nil {
						errCount++
						if errCount > maxRetryCount {
							item.Send(subscriber)
							break observe
						}

						upStream.Stop()
						upStream = source.SubscribeOn()
						continue
					}

					item.Send(subscriber)
					if item.Done() {
						break observe
					}
				}
			}

			wg.Wait()
		})
	}
}
