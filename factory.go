package rxgo

func Empty() Observable {
	next := make(chan Item)
	close(next)
	return &observable{
		iterable: newChannelIterable(next),
	}
}

func FromChannel(next <-chan Item) Observable {
	return &observable{
		iterable: newChannelIterable(next),
	}
}

func Just(item Item, items ...Item) Observable {
	if len(items) > 0 {
		items = append([]Item{item}, items...)
	} else {
		items = []Item{item}
	}
	return &observable{
		iterable: newSliceIterable(items),
	}
}
