package options

type BackpressureStrategy uint32
type Scheduler uint32

const (
	Drop BackpressureStrategy = iota
	Buffer
)

const (
	Blocking Scheduler = iota
	NonBlocking
)

// Option is the configuration of an observable
type Option interface {
	apply(*funcOption)
	Parallelism() int
	Buffer() int
	BackpressureStrategy() BackpressureStrategy
}

// funcOption wraps a function that modifies options into an
// implementation of the Option interface.
type funcOption struct {
	f           func(*funcOption)
	parallelism int
	buffer      int
	bpStrategy  BackpressureStrategy
}

func (fdo *funcOption) Parallelism() int {
	return fdo.parallelism
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

// WithParallelism allows to configure the level of parallelism
func WithParallelism(parallelism int) Option {
	return newFuncOption(func(options *funcOption) {
		options.parallelism = parallelism
	})
}

// WithBufferedChannel allows to configure the capacity of a buffered channel
func WithBufferedChannel(capacity int) Option {
	return newFuncOption(func(options *funcOption) {
		options.buffer = capacity
	})
}

func WithDropBackpressureStrategy() Option {
	return newFuncOption(func(options *funcOption) {
		options.bpStrategy = Drop
	})
}

func WithBufferBackpressureStrategy(buffer int) Option {
	return newFuncOption(func(options *funcOption) {
		options.bpStrategy = Buffer
		options.buffer = buffer
	})
}
