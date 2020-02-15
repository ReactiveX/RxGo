package rxgo

// Iterable is the interface returning an iterable channel.
type Iterable interface {
	Observe() <-chan Item
}
