package rxgo

// BackpressureStrategy is the backpressure strategy type
type BackpressureStrategy uint32

const (
	// None means backpressure management
	None BackpressureStrategy = iota
	// Drop least recent items
	Drop
	// Buffer items
	Buffer
)

// Option is the configuration of an observable
type Option interface {
	apply(*funcOption)
	Buffer() int
	BackpressureStrategy() BackpressureStrategy
}

// funcOption wraps a function that modifies options into an
// implementation of the Option interface.
type funcOption struct {
	f          func(*funcOption)
	buffer     int
	bpStrategy BackpressureStrategy
}

func (fdo *funcOption) Buffer() int {
	return fdo.buffer
}

func (fdo *funcOption) BackpressureStrategy() BackpressureStrategy {
	return fdo.bpStrategy
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

// WithBufferedChannel allows to configure the capacity of a buffered channel
func WithBufferedChannel(capacity int) Option {
	return newFuncOption(func(options *funcOption) {
		options.buffer = capacity
	})
}

// WithoutBackpressureStrategy indicates to apply the None backpressure strategy
func WithoutBackpressureStrategy() Option {
	return newFuncOption(func(options *funcOption) {
		options.bpStrategy = None
	})
}

// WithDropBackpressureStrategy indicates to apply the Drop backpressure strategy
func WithDropBackpressureStrategy() Option {
	return newFuncOption(func(options *funcOption) {
		options.bpStrategy = Drop
	})
}

// WithBufferBackpressureStrategy indicates to apply the Drop backpressure strategy
func WithBufferBackpressureStrategy(buffer int) Option {
	return newFuncOption(func(options *funcOption) {
		options.bpStrategy = Buffer
		options.buffer = buffer
	})
}
