package rx

// Option is the configuration of an observable
type Option interface {
	apply(*options)
}

// options configurable for an observable
type options struct {
	parallelism           int
	channelBufferCapacity int
}

// funcOption wraps a function that modifies options into an
// implementation of the Option interface.
type funcOption struct {
	f func(*options)
}

func (fdo *funcOption) apply(do *options) {
	fdo.f(do)
}

func newFuncOption(f func(*options)) *funcOption {
	return &funcOption{
		f: f,
	}
}

// WithParallelism allows to configure the level of parallelism
func WithParallelism(parallelism int) Option {
	return newFuncOption(func(options *options) {
		options.parallelism = parallelism
	})
}

// parseOptions parse the given options and mutate the options
// structure
func (o *options) parseOptions(opts ...Option) {
	for _, opt := range opts {
		opt.apply(o)
	}
}

// WithBufferedChannel allows to configure the capacity of a buffered channel
func WithBufferedChannel(capacity int) Option {
	return newFuncOption(func(options *options) {
		options.channelBufferCapacity = capacity
	})
}
