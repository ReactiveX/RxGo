package rxgo

type factoryIterable struct {
	factory func() <-chan Item
}

func newColdIterable(factory func() <-chan Item) Iterable {
	return &factoryIterable{factory: factory}
}

func (i *factoryIterable) Observe() <-chan Item {
	return i.factory()
}
