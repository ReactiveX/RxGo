package rxgo

// Represents all of the notifications from the source Observable as next emissions
// marked with their original types within Notification objects.
func Materialize[T any]() OperatorFunc[T, Notification[T]] {
	return func(source IObservable[T]) IObservable[Notification[T]] {
		return newObservable(func(subscriber Subscriber[Notification[T]]) {

		})
	}
}

type NotificationKind int

const (
	NextKind NotificationKind = iota
	ErrorKind
	CompleteKind
)

type Notification[T any] interface {
	Kind() NotificationKind
	Value() T
	Err() error
	Done() bool
	Send(Subscriber[T]) bool
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
