package rxgo

import (
	"sync"
)

// NotificationKind
type NotificationKind int

const (
	NextKind NotificationKind = iota
	ErrorKind
	CompleteKind
)

type ObservableNotification[T any] interface {
	Kind() NotificationKind
	Value() T // returns the underlying value if it's a "Next" notification
	Err() error
	IsEnd() bool
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

func (d notification[T]) IsEnd() bool {
	return d.err != nil || d.done
}

func (d *notification[T]) Send(sub Subscriber[T]) bool {
	select {
	case <-sub.Closed():
		return false
	case sub.Send() <- d:
		return true
	}
}

func Next[T any](v T) Notification[T] {
	return &notification[T]{kind: NextKind, v: v}
}

func Error[T any](err error) Notification[T] {
	return &notification[T]{kind: ErrorKind, err: err}
}

func Complete[T any]() Notification[T] {
	return &notification[T]{kind: CompleteKind, done: true}
}

// Converts an Observable of ObservableNotification objects into the emissions that they represent.
func Dematerialize[T any]() OperatorFunc[ObservableNotification[T], T] {
	return func(source Observable[ObservableNotification[T]]) Observable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(1)

			var (
				upStream = source.SubscribeOn(wg.Done)
				msg      ObservableNotification[T]
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

					if item.Done() {
						break observe
					}

					msg = item.Value()

					switch msg.Kind() {
					case NextKind:
						Next(msg.Value()).Send(subscriber)

					case ErrorKind:
						Error[T](msg.Err()).Send(subscriber)
						break observe

					case CompleteKind:
						Complete[T]().Send(subscriber)
						break observe
					}
				}
			}

			upStream.Stop()

			wg.Wait()
		})
	}
}

// Represents all of the notifications from the source Observable as next emissions marked with their original types within Notification objects.
func Materialize[T any]() OperatorFunc[T, ObservableNotification[T]] {
	return func(source Observable[T]) Observable[ObservableNotification[T]] {
		return newObservable(func(subscriber Subscriber[ObservableNotification[T]]) {
			var (
				wg = new(sync.WaitGroup)
			)

			wg.Add(1)

			var (
				upStream  = source.SubscribeOn(wg.Done)
				completed bool
				msg       Notification[ObservableNotification[T]]
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

					// When the source Observable emits complete, the output Observable will emit next as a Notification of type "Complete", and then it will emit complete as well. When the source Observable emits error, the output will emit next as a Notification of type "Error", and then complete.
					completed = item.Err() != nil || item.Done()
					msg = Next(item.(ObservableNotification[T]))

					if !msg.Send(subscriber) {
						upStream.Stop()
						break observe
					}

					if completed {
						upStream.Stop()
						Complete[ObservableNotification[T]]().Send(subscriber)
						break observe
					}
				}
			}

			wg.Wait()
		})
	}
}
