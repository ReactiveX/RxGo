package bases

type Stream interface {
	Subscribe(EventHandler) (Subscriptor, error)
	Unsubscribe() Subscriptor
	Done()

	Iterator
	//Subscribe(Observer) (Subscriptor, error)
	//SubscribeFunc(func(interface{}), func(error), func()) (Subscriptor, error)
	//SubscribeHandler(EventHandler, ...EventHandler) (Subscriptor, error)
}
