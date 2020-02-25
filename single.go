package rxgo

import "context"

// Single is a observable with a single element.
type Single interface {
	Iterable
	Filter(apply Predicate, opts ...Option) OptionalSingle
	Get(opts ...Option) (Item, error)
	Map(apply Func, opts ...Option) Single
	Run(opts ...Option) Disposed
}

// SingleImpl implements Single.
type SingleImpl struct {
	iterable Iterable
}

// Filter emits only those items from an Observable that pass a predicate test.
func (s *SingleImpl) Filter(apply Predicate, opts ...Option) OptionalSingle {
	return optionalSingle(s, func() operator {
		return &filterOperatorSingle{apply: apply}
	}, true, true, opts...)
}

// Get returns the item. The error returned is if the context has been cancelled.
// This method is blocking.
func (s *SingleImpl) Get(opts ...Option) (Item, error) {
	option := parseOptions(opts...)
	ctx := option.buildContext()

	observe := s.Observe(opts...)
	for {
		select {
		case <-ctx.Done():
			return Item{}, ctx.Err()
		case v := <-observe:
			return v, nil
		}
	}
}

// Map transforms the items emitted by a Single by applying a function to each item.
func (s *SingleImpl) Map(apply Func, opts ...Option) Single {
	return single(s, func() operator {
		return &mapOperatorSingle{apply: apply}
	}, false, true, opts...)
}

type mapOperatorSingle struct {
	apply Func
}

func (op *mapOperatorSingle) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	res, err := op.apply(ctx, item.V)
	if err != nil {
		Error(err).SendContext(ctx, dst)
		operatorOptions.stop()
		return
	}
	Of(res).SendContext(ctx, dst)
}

func (op *mapOperatorSingle) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *mapOperatorSingle) end(_ context.Context, _ chan<- Item) {
}

func (op *mapOperatorSingle) gatherNext(ctx context.Context, item Item, dst chan<- Item, _ operatorOptions) {
	switch item.V.(type) {
	case *mapOperatorSingle:
		return
	}
	item.SendContext(ctx, dst)
}

// Observe observes a Single by returning its channel.
func (s *SingleImpl) Observe(opts ...Option) <-chan Item {
	return s.iterable.Observe(opts...)
}

type filterOperatorSingle struct {
	apply Predicate
}

func (op *filterOperatorSingle) next(ctx context.Context, item Item, dst chan<- Item, _ operatorOptions) {
	if op.apply(item.V) {
		item.SendContext(ctx, dst)
	}
}

func (op *filterOperatorSingle) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *filterOperatorSingle) end(_ context.Context, _ chan<- Item) {
}

func (op *filterOperatorSingle) gatherNext(_ context.Context, _ Item, _ chan<- Item, _ operatorOptions) {
}

// Run creates an observer without consuming the emitted items.
func (s *SingleImpl) Run(opts ...Option) Disposed {
	dispose := make(chan struct{})
	option := parseOptions(opts...)
	ctx := option.buildContext()

	go func() {
		defer close(dispose)
		observe := s.Observe(opts...)
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
