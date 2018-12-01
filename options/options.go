package options

// Option is the configuration of an observable
type Option interface {
	apply(*funcOption)
	Parallelism() int
	BufferedChannelCapacity() int
}

// funcOption wraps a function that modifies options into an
// implementation of the Option interface.
type funcOption struct {
	f                       func(*funcOption)
	parallelism             int
	bufferedChannelCapacity int
}

func (fdo *funcOption) Parallelism() int {
	return fdo.parallelism
}

func (fdo *funcOption) BufferedChannelCapacity() int {
	return fdo.bufferedChannelCapacity
}

func (fdo *funcOption) apply(do *funcOption) {
	fdo.f(do)
}

func newFuncOption(f func(*funcOption)) *funcOption {
	return &funcOption{
		f: f,
	}
}

// ParseOptions parse the given options and mutate the options
// structure
func ParseOptions(opts ...Option) Option {
	o := new(funcOption)
	for _, opt := range opts {
		opt.apply(o)
	}
	return o
}

// WithParallelism allows to configure the level of parallelism
func WithParallelism(parallelism int) Option {
	return newFuncOption(func(options *funcOption) {
		options.parallelism = parallelism
	})
}

// WithBufferedChannel allows to configure the capacity of a buffered channel
func WithBufferedChannel(capacity int) Option {
	return newFuncOption(func(options *funcOption) {
		options.bufferedChannelCapacity = capacity
	})
}
