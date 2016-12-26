package subject

import (
	"github.com/jochasinga/grx/observable"
	"github.com/jochasinga/grx/observer"
)

type Subject struct {
	*observable.Observable
	*observer.Observer
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
