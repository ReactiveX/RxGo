package rxgo

type factoryIterable struct {
	factory func() <-chan Item
}

func newColdIterable(factory func() <-chan Item) Iterable {
	return &factoryIterable{factory: factory}
}

func (i *factoryIterable) Observe(_ ...Option) <-chan Item {
	return i.factory()
}
