package rxgo

// Empty creates an Observable with no item and terminate immediately.
func Empty() Observable {
	next := make(chan Item)
	close(next)
	return &observable{
		iterable: newChannelIterable(next),
	}
}

// FromChannel creates a cold observable from a channel.
func FromChannel(next <-chan Item) Observable {
	return &observable{
		iterable: newChannelIterable(next),
	}
}

// Just creates an Observable with the provided items.
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
