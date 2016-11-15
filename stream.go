package grx

// Stream is an interface which Observable implements.
type Stream interface {
	Subscribe(*Observer) (*Subscription, error)
	SubscribeFunc(func(interface{}), func(error), func()) (*Subscription, error)
	SubscribeHandler(Handler, ...Handler) (*Subscription, error)
}
