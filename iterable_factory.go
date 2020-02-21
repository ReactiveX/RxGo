package rxgo

type factoryIterable struct {
	factory func(opts ...Option) <-chan Item
}

func newFactoryIterable(factory func(opts ...Option) <-chan Item) Iterable {
	return &factoryIterable{factory: factory}
}

func (i *factoryIterable) Observe(opts ...Option) <-chan Item {
	return i.factory(opts...)
}
