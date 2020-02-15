package rxgo

// Iterable is the interface returning an iterable channel.
type Iterable interface {
	Observe(opts ...Option) <-chan Item
}
