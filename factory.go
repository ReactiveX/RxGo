package rxgo

import "context"

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

func FromEventSource(ctx context.Context, next <-chan Item, strategy BackpressureStrategy) Observable {
	return &observable{
		iterable: newEventSourceIterable(ctx, next, strategy),
	}
}

func FromItem(item Item) Single {
	return &single{
		iterable: newSliceIterable([]Item{item}),
	}
}

// FromItems creates an Observable with the provided items.
func FromItems(item Item, items ...Item) Observable {
	if len(items) > 0 {
		items = append([]Item{item}, items...)
	} else {
		items = []Item{item}
	}
	return &observable{
		iterable: newSliceIterable(items),
	}
}
