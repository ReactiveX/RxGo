package rxgo

import (
	"context"
	"fmt"
	"io"
)

type InOutBoundObservable interface {
	Observable
	MapInBound(marshaller InBoundMarshaller, opts ...Option) Observable
	MapOutBound(marshaller OutBoundMarshaller, opts ...Option) Observable
	DoOnNextInBound(nextFunc NextInBoundFunc, opts ...Option) Disposed
	DoOnNextOutBound(nextFunc NextOutBoundFunc, opts ...Option) Disposed
}

type ReadWriterSize interface {
	io.ReadWriter
	Size() int
}

type InBoundMarshaller func(context.Context, io.Reader) (io.Reader, error)
type OutBoundMarshaller func(context.Context, ReadWriterSize) (ReadWriterSize, error)
type NextInBoundFunc func(context.Context, io.Reader)
type NextOutBoundFunc func(context.Context, ReadWriterSize)

func (o *ObservableImpl) DoOnNextInBound(nextFunc NextInBoundFunc, opts ...Option) Disposed {
	dispose := make(chan struct{})
	handler := func(ctx context.Context, src <-chan Item) {
		defer close(dispose)
		for {
			select {
			case <-ctx.Done():
				return
			case i, ok := <-src:
				if !ok {
					return
				}
				if i.Error() {
					return
				}
				if reader, ok := i.V.(io.Reader); ok {
					nextFunc(ctx, reader)
					continue
				}
				return
			}
		}
	}

	option := parseOptions(opts...)
	ctx := option.buildContext()
	go handler(ctx, o.Observe(opts...))
	return dispose
}

func (o *ObservableImpl) DoOnNextOutBound(nextFunc NextOutBoundFunc, opts ...Option) Disposed {
	dispose := make(chan struct{})
	handler := func(ctx context.Context, src <-chan Item) {
		defer close(dispose)
		for {
			select {
			case <-ctx.Done():
				return
			case i, ok := <-src:
				if !ok {
					return
				}
				if i.Error() {
					return
				}
				if rw, ok := i.V.(ReadWriterSize); ok {
					nextFunc(ctx, rw)
					continue
				}
				return
			}
		}
	}

	option := parseOptions(opts...)
	ctx := option.buildContext()
	go handler(ctx, o.Observe(opts...))
	return dispose
}

func (o *ObservableImpl) MapInBound(marshaller InBoundMarshaller, opts ...Option) Observable {
	return o.Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		if reader, ok := i.(io.Reader); ok {
			return marshaller(ctx, reader)
		}
		return nil, fmt.Errorf("type not a io.Reader")
	}, opts...)
}
func (o *ObservableImpl) MapOutBound(marshaller OutBoundMarshaller, opts ...Option) Observable {
	return o.Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		if rw, ok := i.(ReadWriterSize); ok {
			return marshaller(ctx, rw)
		}
		return nil, fmt.Errorf("type not a io.Reader")
	}, opts...)
}
