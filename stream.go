package grx

type Stream interface {
	Iterator

	//Subscribe(Observer) (Subscriptor, error)
	//SubscribeFunc(func(interface{}), func(error), func()) (Subscriptor, error)
	//SubscribeHandler(EventHandler, ...EventHandler) (Subscriptor, error)
	Subscribe(EventHandler) (Subscriptor, error)
}

// Stream is an interface which Observable implements.
type Observable interface {
	/*
		Subscribe(*Observer) (*Subscription, error)
		SubscribeFunc(func(interface{}), func(error), func()) (*Subscription, error)
		SubscribeHandler(Handler, ...Handler) (*Subscription, error)
	*/

	/*
		Subscribe(Observer) (Subscriptor, error)
		SubscribeFunc(func(interface{}), func(error), func()) (Subscriptor, error)
		SubscribeHandler(Handler, ...Handler) (Subscriptor, error)
	*/
	Stream

	/******** Unexported ********/
	getInnerObserver() Observer
	setInnerObserverTo(Observer) Observer
	addItem(interface{})
	//signalDone()
	terminate()
	unsubscribe()
	getC() chan interface{}

	//Iterator
}
