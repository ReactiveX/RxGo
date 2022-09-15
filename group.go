package rxgo

type groupedObservable[K comparable, T any] struct {
	key K
	observableWrapper[T]
}

var (
	_ GroupedObservable[string, any] = (*groupedObservable[string, any])(nil)
)

// func newGroupedObservable[K comparable, T any]() *groupedObservable[K, T] {
// 	obs := &groupedObservable[K, T]{}
// 	obs.connector = func() Subject[T] {
// 		return NewSubscriber[T]()
// 	}
// 	return obs
// }

func (g *groupedObservable[K, T]) Key() K {
	return g.key
}
