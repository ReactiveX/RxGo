package rxgo

type groupedObservable[K comparable, T any] struct {
	observableWrapper[T]
	key       K
	connector Subscriber[T]
}

var (
	_ GroupedObservable[string, any] = (*groupedObservable[string, any])(nil)
)

func newGroupedObservable[K comparable, T any]() *groupedObservable[K, T] {
	return &groupedObservable[K, T]{connector: NewSubscriber[T]()}
}

func (g *groupedObservable[K, T]) Key() K {
	return g.key
}
