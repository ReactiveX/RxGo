package grx

// Stream is an interface which Observable implements.
type Stream interface {
	Subscribe(*Observer) (*Subscription, error)
	SubscribeWith(NextFunc, ErrFunc, DoneFunc) (*Subscription, error)
	SubscribeHandler(Handler) (*Subscription, error)
}
