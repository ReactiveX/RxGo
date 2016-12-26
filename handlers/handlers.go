package handlers

import (
	"github.com/jochasinga/grx/bases"
)

type (
	NextFunc func(bases.Item)
	ErrFunc  func(error)
	DoneFunc func()
)

func (handle NextFunc) Handle(e bases.Emitter) {
	if item, err := e.Emit(); err == nil {
		if item != nil {
			handle(item)
		}
	}
}

func (handle ErrFunc) Handle(e bases.Emitter) {
	if item, err := e.Emit(); err != nil {
		if item == nil {
			handle(err)
		}
	}
}

func (handle DoneFunc) Handle(e bases.Emitter) {
	handle()
}
