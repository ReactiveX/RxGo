package observable

// Option is the configuration of an observable
type Option interface {
	apply(*options)
}

// options configurable for an observable
type options struct {
	parallelism int
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
