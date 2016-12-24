package handlers

import (
	"github.com/jochasinga/bases"
	"github.com/jochasinga/grx"
)

type EventHandler interface {
	Apply(grx.Emitter)
}

type (
	NextFunc func(bases.Item)
	ErrFunc  func(error)
	DoneFunc func()
)

func (handle NextFunc) Apply(e grx.Emitter) {
	if item, err := e.Emit(); err == nil {
		if item != nil {
			handle(item)
		}
	} else {
		panic("Both Item and error are empty (Expecting an Item)")
	}
}

func (handle ErrFunc) Apply(e grx.Emitter) {
	if item, err := e.Emit(); err != nil {
		if item == nil {
			handle(err)
		}
	} else {
		panic("Both Item and error are empty (Expecting an error)")
	}
}

func (handle DoneFunc) Apply(e grx.Emitter) {
	handle()
}
