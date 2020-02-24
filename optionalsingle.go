package rxgo

import "context"

// OptionalSingleEmpty is the constant returned when an OptionalSingle is empty.
var OptionalSingleEmpty = Item{}

// OptionalSingle is an optional single.
type OptionalSingle interface {
	Iterable
	Get(opts ...Option) (Item, error)
	Map(apply Func, opts ...Option) OptionalSingle
	Run(opts ...Option) Disposed
}

// OptionalSingleImpl implements OptionalSingle.
type OptionalSingleImpl struct {
	iterable Iterable
}

// Get returns the item or rxgo.OptionalEmpty. The error returned is if the context has been cancelled.
// This method is blocking.
func (o *OptionalSingleImpl) Get(opts ...Option) (Item, error) {
	option := parseOptions(opts...)
	ctx := option.buildContext()

	observe := o.Observe(opts...)
	for {
		select {
		case <-ctx.Done():
			return Item{}, ctx.Err()
		case v, ok := <-observe:
			if !ok {
				return OptionalSingleEmpty, nil
			}
			return v, nil
		}
	}
}

// Map transforms the items emitted by an OptionalSingle by applying a function to each item.
func (o *OptionalSingleImpl) Map(apply Func, opts ...Option) OptionalSingle {
	return optionalSingle(o, func() operator {
		return &mapOperatorOptionalSingle{apply: apply}
	}, false, true, opts...)
}

// Observe observes an OptionalSingle by returning its channel.
func (o *OptionalSingleImpl) Observe(opts ...Option) <-chan Item {
	return o.iterable.Observe(opts...)
}

type mapOperatorOptionalSingle struct {
	apply Func
}

func (op *mapOperatorOptionalSingle) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	res, err := op.apply(ctx, item.V)
	if err != nil {
		dst <- Error(err)
		operatorOptions.stop()
		return
	}
	dst <- Of(res)
}

func (op *mapOperatorOptionalSingle) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *mapOperatorOptionalSingle) end(_ context.Context, _ chan<- Item) {
}

func (op *mapOperatorOptionalSingle) gatherNext(_ context.Context, item Item, dst chan<- Item, _ operatorOptions) {
	switch item.V.(type) {
	case *mapOperatorOptionalSingle:
		return
	}
	dst <- item
}

// Run creates an observer without consuming the emitted items.
func (o *OptionalSingleImpl) Run(opts ...Option) Disposed {
	dispose := make(chan struct{})
	option := parseOptions(opts...)
	ctx := option.buildContext()

	go func() {
		defer close(dispose)
		observe := o.Observe(opts...)
		for {
			select {
			case <-ctx.Done():
				return
			case _, ok := <-observe:
				if !ok {
					return
				}
			}
		}
	}()

	return dispose
}
