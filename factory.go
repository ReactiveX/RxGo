package rxgo

func FromChannel(next <-chan Item) Observable {
	return &observable{
		iterable: newIterable(next),
	}
}
