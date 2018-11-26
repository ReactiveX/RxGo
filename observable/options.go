package observable

type Option interface {
	apply(*options)
}

type options struct {
	parallelism int
}

// funcOption wraps a function that modifies dialOptions into an
// implementation of the DialOption interface.
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

func WithParallelism(parallelism int) Option {
	return newFuncOption(func(options *options) {
		options.parallelism = parallelism
	})
}
