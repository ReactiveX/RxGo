package rx

type Disposable interface {
	Dispose()
	IsDisposed() bool
}
