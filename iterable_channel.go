package rxgo

type channelIterable struct {
	next <-chan Item
}

func newChannelIterable(next <-chan Item) Iterable {
	return &channelIterable{next: next}
}

func (i *channelIterable) Observe(_ ...Option) <-chan Item {
	return i.next
}
