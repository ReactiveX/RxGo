package rxgo

// Observer represents a group of EventHandlers.
type Observer interface {
	EventHandler
	Disposable

	OnNext(item interface{})
	OnError(err error)
	OnDone()

	Block()
	setItemChannel(chan interface{})
	getItemChannel() chan interface{}
}

type observer struct {
	// itemChannel is the internal channel used to receive items from the parent observable
	itemChannel chan interface{}
	// nextHandler is the handler for the next items
	nextHandler NextFunc
	// errHandler is the error handler
	errHandler ErrFunc
	// doneHandler is the handler once an observable is done
	doneHandler DoneFunc
	// disposedChannel is the notification channel used when an observer is disposed
	disposedChannel chan struct{}
}

func (o *observer) setItemChannel(ch chan interface{}) {
	o.itemChannel = ch
}

func (o *observer) getItemChannel() chan interface{} {
	return o.itemChannel
}

// NewObserver constructs a new Observer instance with default Observer and accept
// any number of EventHandler
func NewObserver(eventHandlers ...EventHandler) Observer {
	ob := observer{
		disposedChannel: make(chan struct{}),
	}

	if len(eventHandlers) > 0 {
		for _, handler := range eventHandlers {
			switch handler := handler.(type) {
			case NextFunc:
				ob.nextHandler = handler
			case ErrFunc:
				ob.errHandler = handler
			case DoneFunc:
				ob.doneHandler = handler
			case *observer:
				ob = *handler
			}
		}
	}

	if ob.nextHandler == nil {
		ob.nextHandler = func(interface{}) {}
	}
	if ob.errHandler == nil {
		ob.errHandler = func(err error) {}
	}
	if ob.doneHandler == nil {
		ob.doneHandler = func() {}
	}

	return &ob
}

// Handle registers Observer to EventHandler.
func (o *observer) Handle(item interface{}) {
	switch item := item.(type) {
	case error:
		o.errHandler(item)
		return
	default:
		o.nextHandler(item)
	}
}

func (o *observer) Dispose() {
	close(o.disposedChannel)
}

func (o *observer) IsDisposed() bool {
	select {
	case <-o.disposedChannel:
		return true
	default:
		return false
	}
}

// OnNext applies Observer's NextHandler to an Item
func (o *observer) OnNext(item interface{}) {
	if !o.IsDisposed() {
		switch item := item.(type) {
		case error:
			return
		default:
			if o.nextHandler != nil {
				o.nextHandler(item)
			}
		}
	} else {
		// TODO
	}
}

// OnError applies Observer's ErrHandler to an error
func (o *observer) OnError(err error) {
	if !o.IsDisposed() {
		if o.errHandler != nil {
			o.errHandler(err)
			o.Dispose()
		}
	} else {
		// TODO
	}
}

// OnDone terminates the Observer's internal Observable
func (o *observer) OnDone() {
	if !o.IsDisposed() {
		if o.doneHandler != nil {
			o.doneHandler()
			o.Dispose()
		}
	} else {
		// TODO
	}
}

// OnDone terminates the Observer's internal Observable
func (o *observer) Block() {
	select {
	case <-o.disposedChannel:
		return
	}
}
