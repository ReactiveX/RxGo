package rxgo

// Option handles configurable options.
type Option interface {
	apply(*funcOption)
	toBeBuffered() (bool, int)
}

type funcOption struct {
	f        func(*funcOption)
	toBuffer bool
	buffer   int
}

func (fdo *funcOption) toBeBuffered() (bool, int) {
	return fdo.toBuffer, fdo.buffer
}

func (fdo *funcOption) apply(do *funcOption) {
	fdo.f(do)
}

func newFuncOption(f func(*funcOption)) *funcOption {
	return &funcOption{
		f: f,
	}
}

func parseOptions(opts ...Option) Option {
	o := new(funcOption)
	for _, opt := range opts {
		opt.apply(o)
	}
	return o
}

// WithBufferedChannel allows to configure the capacity of a buffered channel.
func WithBufferedChannel(capacity int) Option {
	return newFuncOption(func(options *funcOption) {
		options.toBuffer = true
		options.buffer = capacity
	})
}
