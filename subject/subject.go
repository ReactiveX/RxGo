package subject

import (
	"github.com/jochasinga/grx/bases"
)

type Subject struct {
	bases.Stream
	bases.Sentinel
}

var DefaultSubject = &Subject{}

func New(fs ...func(*Subject)) *Subject {
	s := DefaultSubject
	if len(fs) > 0 {
		for _, f := range fs {
			f(s)
		}
	}
	return s
}

/*
func (s *Subject) Done() {
	if o, ok := s.Stream.(*observable.Observable); ok {
		o.Done()
	}
}
*/
