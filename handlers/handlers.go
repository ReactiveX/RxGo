package handlers

import (
	"github.com/jochasinga/grx/bases"
)

type (
	NextFunc func(bases.Item)
	ErrFunc  func(error)
	DoneFunc func()
)

func (handle NextFunc) Apply(e bases.Emitter) {
	if item, err := e.Emit(); err == nil {
		if item != nil {
			handle(item)
		}
	} else {
		panic("Both Item and error are empty (Expecting an Item)")
	}
}

func (handle ErrFunc) Apply(e bases.Emitter) {
	if item, err := e.Emit(); err != nil {
		if item == nil {
			handle(err)
		}
	} else {
		panic("Both Item and error are empty (Expecting an error)")
	}
}

func (handle DoneFunc) Apply(e bases.Emitter) {
	handle()
}
