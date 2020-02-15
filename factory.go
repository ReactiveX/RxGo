package rxgo

func Empty() Observable {
	next := make(chan Item)
	close(next)
	return &observable{
		iterable: newIterable(next),
	}
}

func FromChannel(next <-chan Item) Observable {
	return &observable{
		iterable: newIterable(next),
	}
}
