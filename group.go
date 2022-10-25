package rxgo

type groupedObservable[K comparable, T any] struct {
	key K
	observableWrapper[T]
}

var (
	_ GroupedObservable[string, any] = (*groupedObservable[string, any])(nil)
)

func NewGroupedObservable[K comparable, T any](key K, connector func() Subject[T]) GroupedObservable[K, T] {
	obs := &groupedObservable[K, T]{}
	obs.key = key
	obs.connector = connector
	return obs
}

// func newGroupedObservable[K comparable, T any](key K) GroupedObservable[K, T] {
// 	obs := &groupedObservable[K, T]{}
// 	obs.key = key
// 	obs.connector = func() Subject[T] {
// 		return NewSubscriber[T]()
// 	}
// 	return obs
// }

func (g *groupedObservable[K, T]) Key() K {
	return g.key
}
