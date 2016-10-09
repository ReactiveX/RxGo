package tooth

import "errors"

type Observer func(e Event, cxt Context) interface{}

type Context map[string]interface{}

type Event struct {
	Value interface{}
	Error error
	Done bool
}

type Observable struct {
	Name string
	Duplex map[string]chan Event
	Observer Observer
}

func (o *Observable) Subscribe(Observer) error {
	// Returns an error if there's anything wrong.
	return errors.New("Fake errror")
}

func NewObservable(name string) *Observable {
	return &Observable{
		Name:  name,
		Duplex: make(map[string]chan Event),
	}
}

