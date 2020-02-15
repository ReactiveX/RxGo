package rxgo

type BlockingObserver interface {
	Block()
}

type observer struct {
	disposedChannel chan struct{}
}

func NewObserver(nextFunc func(interface{}), errFunc func(error), doneFunc func()) BlockingObserver {
	return &observer{
		disposedChannel: make(chan struct{}),
	}
}

func (o *observer) Block() {
	<-o.disposedChannel
}
