package rxgo

// Iterable is the basic type that can be observed.
type Iterable interface {
	Observe(opts ...Option) <-chan Item
}
