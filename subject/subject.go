package subject

import (
	"github.com/jochasinga/observable"
	"github.com/jochasinga/observer"
)

type Subject struct {
	*Observable
	*Observer
}

func New(fs ...func(*Subject)) *Subject {
	s := &Subject{
		Observable: observable.New(),
		Observer:   observer.New(),
	}
	if len(fs) > 0 {
		for _, f := range fs {
			f(s)
		}
	}
	return s
}
