package rxgo

import (
	"fmt"
)

func skipPredicate[T any](T, uint) bool {
	return true
}

func minimum[T any](a T, b T) int8 {
	positive, negative := int8(-1), int8(1)
	switch any(a).(type) {
	case string:
		if any(a).(string) < any(b).(string) {
			return positive
		}
		return negative
	case []byte:
		if string(any(a).([]byte)) < string(any(b).([]byte)) {
			return positive
		}
		return negative
	case int8:
		if any(a).(int8) < any(b).(int8) {
			return positive
		}
		return negative
	case int:
		if any(a).(int) < any(b).(int) {
			return positive
		}
		return negative
	case int32:
		if any(a).(int32) < any(b).(int32) {
			return positive
		}
		return negative
	case int64:
		if any(a).(int64) < any(b).(int64) {
			return positive
		}
		return negative
	case uint8:
		if any(a).(uint8) < any(b).(uint8) {
			return positive
		}
		return negative
	case uint:
		if any(a).(uint) < any(b).(uint) {
			return positive
		}
		return negative
	case uint32:
		if any(a).(uint32) < any(b).(uint32) {
			return positive
		}
		return negative
	case uint64:
		if any(a).(uint64) < any(b).(uint64) {
			return positive
		}
		return negative
	case float32:
		if any(a).(float32) < any(b).(float32) {
			return positive
		}
		return negative
	case float64:
		if any(a).(float64) < any(b).(float64) {
			return positive
		}
		return negative

	default:
		if fmt.Sprintf("%v", a) < fmt.Sprintf("%v", b) {
			return positive
		}
		return negative
	}
}

func maximum[T any](a T, b T) int8 {
	positive, negative := int8(1), int8(-1)
	switch any(a).(type) {
	case string:
		if any(a).(string) > any(b).(string) {
			return positive
		}
		return negative
	case []byte:
		if string(any(a).([]byte)) > string(any(b).([]byte)) {
			return positive
		}
		return negative
	case int8:
		if any(a).(int8) > any(b).(int8) {
			return positive
		}
		return negative
	case int:
		if any(a).(int) > any(b).(int) {
			return positive
		}
		return negative
	case int32:
		if any(a).(int32) > any(b).(int32) {
			return positive
		}
		return negative
	case int64:
		if any(a).(int64) > any(b).(int64) {
			return positive
		}
		return negative
	case uint8:
		if any(a).(uint8) > any(b).(uint8) {
			return positive
		}
		return negative
	case uint:
		if any(a).(uint) > any(b).(uint) {
			return positive
		}
		return negative
	case uint32:
		if any(a).(uint32) > any(b).(uint32) {
			return positive
		}
		return negative
	case uint64:
		if any(a).(uint64) > any(b).(uint64) {
			return positive
		}
		return negative
	case float32:
		if any(a).(float32) > any(b).(float32) {
			return positive
		}
		return negative
	case float64:
		if any(a).(float64) > any(b).(float64) {
			return positive
		}
		return negative
	default:
		if fmt.Sprintf("%v", a) > fmt.Sprintf("%v", b) {
			return positive
		}
		return negative
	}
}

func createOperatorFunc[T any, R any](
	source Observable[T],
	onNext func(Observer[R], T),
	onError func(Observer[R], error),
	onComplete func(Observer[R]),
) Observable[R] {
	return newObservable(func(subscriber Subscriber[R]) {
		var (
			stop     bool
			upStream = source.SubscribeOn()
		)

		obs := &consumerObserver[R]{
			onNext: func(v R) {
				if !Next(v).Send(subscriber) {
					stop = true
				}
			},
			onError: func(err error) {
				upStream.Stop()
				stop = true
				select {
				case <-subscriber.Closed():
					return
				case subscriber.Send() <- Error[R](err):
				}
			},
			onComplete: func() {
				// Inform the up stream to stop emit value
				upStream.Stop()
				stop = true
				Complete[R]().Send(subscriber)
			},
		}

		for !stop {
			select {
			// If only the stream terminated, break it
			case <-subscriber.Closed():
				stop = true
				upStream.Stop()
				return

			case item, ok := <-upStream.ForEach():
				if !ok {
					// If only the data stream closed, break it
					stop = true
					return
				}

				if err := item.Err(); err != nil {
					onError(obs, err)
					return
				}

				if item.Done() {
					onComplete(obs)
					return
				}

				onNext(obs, item.Value())
			}
		}
	})
}
