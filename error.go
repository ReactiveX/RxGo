package rxgo

import (
	"errors"
	"sync"
	"time"

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

// Catches errors on the observable to be handled by returning a new observable or throwing an error.
func Catch[T any](catch func(err error, caught Observable[T]) Observable[T]) OperatorFunc[T, T] {
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

						item.Send(subscriber)
						if item.IsEnd() {
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
	Delay          time.Duration
	ResetOnSuccess bool
}

type retryConfig interface {
	constraints.Unsigned | RetryConfig
}

// Returns an Observable that mirrors the source Observable with the exception of an error.
func Retry[T any, C retryConfig](config ...C) OperatorFunc[T, T] {
	var (
		maxRetryCount  = int64(-1)
		delay          = time.Duration(0)
		resetOnSuccess bool
	)
	if len(config) > 0 {
		switch v := any(config[0]).(type) {
		case RetryConfig:
			if v.Count > 0 {
				maxRetryCount = int64(v.Count)
			}
			if v.Delay > 0 {
				delay = v.Delay
			}
			resetOnSuccess = v.ResetOnSuccess
		case uint8:
			maxRetryCount = int64(v)
		case uint16:
			maxRetryCount = int64(v)
		case uint32:
			maxRetryCount = int64(v)
		case uint64:
			maxRetryCount = int64(v)
		case uint:
			maxRetryCount = int64(v)
		}
	}
	return func(source Observable[T]) Observable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				wg       = new(sync.WaitGroup)
				errCount = int64(0)
				upStream Subscriber[T]
				forEach  <-chan Notification[T]
			)

			setupStream := func(first bool) {
				wg.Add(1)
				if delay > 0 && !first {
					time.Sleep(delay)
				}
				upStream = source.SubscribeOn(wg.Done)
				forEach = upStream.ForEach()
			}

			setupStream(true)

		observe:
			//  If count is omitted, retry will try to resubscribe on errors infinite number of times.
			for {
				select {
				case <-subscriber.Closed():
					upStream.Stop()
					break observe

				case item, ok := <-forEach:
					if !ok {
						break observe
					}

					if err := item.Err(); err != nil {
						errCount++
						if errCount > maxRetryCount {
							item.Send(subscriber)
							break observe
						}

						setupStream(false)
						continue
					}

					if errCount > 0 && resetOnSuccess {
						errCount = 0
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
