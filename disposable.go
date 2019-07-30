package rxgo

type Disposable interface {
	Dispose()
	IsDisposed() bool
	Notify(chan<- struct{})
}
