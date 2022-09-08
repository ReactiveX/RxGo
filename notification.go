package rxgo

import (
	"sync"
)

// Represents all of the notifications from the source Observable as next emissions
// marked with their original types within Notification objects.
func Materialize[T any]() OperatorFunc[T, ObservableNotification[T]] {
	return func(source IObservable[T]) IObservable[ObservableNotification[T]] {
		return newObservable(func(subscriber Subscriber[ObservableNotification[T]]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(1)

			var (
				upStream  = source.SubscribeOn(wg.Done)
				completed bool
				notice    Notification[ObservableNotification[T]]
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

					// When the source Observable emits complete,
					// the output Observable will emit next as a Notification of type "complete",
					// and then it will emit complete as well.
					// When the source Observable emits error,
					// the output will emit next as a Notification of type "error", and then complete.
					completed = item.Err() != nil || item.Done()
					notice = NextNotification(item.(ObservableNotification[T]))

					if !notice.Send(subscriber) {
						upStream.Stop()
						break observe
					}

					if completed {
						CompleteNotification[ObservableNotification[T]]().Send(subscriber)
						upStream.Stop()
						break observe
					}
				}
			}

			wg.Wait()
		})
	}
}

// NotificationKind
type NotificationKind int

const (
	NextKind NotificationKind = iota
	ErrorKind
	CompleteKind
)

type ObservableNotification[T any] interface {
	Kind() NotificationKind
	Value() T
	Err() error
}

type Notification[T any] interface {
	ObservableNotification[T]
	Send(Subscriber[T]) bool
	Done() bool
}

type notification[T any] struct {
	kind NotificationKind
	v    T
	err  error
	done bool
}

var _ Notification[any] = (*notification[any])(nil)

func (d notification[T]) Kind() NotificationKind {
	return d.kind
}

func (d notification[T]) Value() T {
	return d.v
}

func (d notification[T]) Err() error {
	return d.err
}

func (d notification[T]) Done() bool {
	return d.done
}

func (d *notification[T]) Send(sub Subscriber[T]) bool {
	select {
	case <-sub.Closed():
		return false
	case sub.Send() <- d:
		return true
	}
}

func NextNotification[T any](v T) Notification[T] {
	return &notification[T]{kind: NextKind, v: v}
}

func ErrorNotification[T any](err error) Notification[T] {
	return &notification[T]{kind: ErrorKind, err: err}
}

func CompleteNotification[T any]() Notification[T] {
	return &notification[T]{kind: CompleteKind, done: true}
}
