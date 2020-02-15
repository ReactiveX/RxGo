package rxgo

import "context"

// Option handles configurable options.
type Option interface {
	apply(*funcOption)
	withBuffer() (bool, int)
	withContext() (bool, context.Context)
}

type funcOption struct {
	f        func(*funcOption)
	toBuffer bool
	buffer   int
	ctx      context.Context
}

func (fdo *funcOption) withBuffer() (bool, int) {
	return fdo.toBuffer, fdo.buffer
}

func (fdo *funcOption) withContext() (bool, context.Context) {
	return fdo.ctx != nil, fdo.ctx
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

// WithContext allows to pass a context.
func WithContext(ctx context.Context) Option {
	return newFuncOption(func(options *funcOption) {
		options.ctx = ctx
	})
}
