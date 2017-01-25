package handlers

type (
	NextFunc func(interface{})
	ErrFunc  func(error)
	DoneFunc func()
)

func (handle NextFunc) Handle(item interface{}) {
	switch item := item.(type) {
	case error:
		return
	default:
		handle(item)
	}
}

func (handle ErrFunc) Handle(item interface{}) {
	switch item := item.(type) {
	case error:
		handle(item)
	default:
		return
	}
}

func (handle DoneFunc) Handle(item interface{}) {
	handle()
}
