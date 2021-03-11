package rxgo

import (
	"context"
	"fmt"
	"github.com/bhbosman/goprotoextra"
	"time"
)

type StreamDirection uint8

const (
	StreamDirectionInbound StreamDirection = iota
	StreamDirectionOutbound
)

type IPublishToConnectionManager interface {
	PublishStackData(index int, connectionId, name string, direction StreamDirection, msgValue, byteValue int)
}

type nullPublisher struct {
}

func (self nullPublisher) PublishStackData(index int, connectionId, name string, direction StreamDirection, msgValue, byteValue int) {
	// do nothing
}

var nullPublisherInstance = &nullPublisher{}

type InOutBoundObservable interface {
	Observable
	MapInOutBound(index int, connectionId string, stackName string, streamDirection StreamDirection, connectionManager IPublishToConnectionManager, marshaller InOutBoundMarshaller, opts ...Option) Observable
	DoOnNextInOutBound(index int, connectionId string, stackName string, streamDirection StreamDirection, connectionManager IPublishToConnectionManager, nextFunc NextInOutBoundFunc, opts ...Option) Disposed
	DoNextExternal(index int, connectionId string, stackName string, streamDirection StreamDirection, connectionManager IPublishToConnectionManager, nextFunc NextExternalFunc, opts ...Option) Disposed
}
type EmptyQueue struct {
	Count int
}
type InOutBoundMarshaller func(context.Context, goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error)
type NextInOutBoundFunc func(context.Context, goprotoextra.ReadWriterSize)
type NextExternalFunc func(bool, interface{})

type NextExternal struct {
	External bool
	Data     interface{}
}

func NewNextExternal(external bool, data interface{}) *NextExternal {
	return &NextExternal{External: external, Data: data}
}

func (o *ObservableImpl) DoNextExternal(index int, connectionId string, stackName string, streamDirection StreamDirection, connectionManager IPublishToConnectionManager, nextFunc NextExternalFunc, opts ...Option) Disposed {
	if connectionManager == nil {
		connectionManager = nullPublisherInstance
	}
	dispose := make(chan struct{})
	handler := func(ctx context.Context, src <-chan Item) {
		defer close(dispose)
		count := 0
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
				if nextExternal, ok := i.V.(*NextExternal); ok {
					nextFunc(nextExternal.External, nextExternal.Data)
					count++
					if len(src) == 0 {
						nextFunc(false, &EmptyQueue{
							Count: count,
						})
						count = 0
					}
					continue
				}
			}
		}
	}

	option := parseOptions(opts...)
	ctx := option.buildContext()
	go handler(ctx, o.Observe(opts...))
	return dispose
}

func (o *ObservableImpl) DoOnNextInOutBound(index int, connectionId string, stackName string, streamDirection StreamDirection, connectionManager IPublishToConnectionManager, nextFunc NextInOutBoundFunc, opts ...Option) Disposed {
	dispose := make(chan struct{})
	if connectionManager == nil {
		connectionManager = nullPublisherInstance
	}

	handler := func(ctx context.Context, src <-chan Item) {
		defer close(dispose)

		var lastUpdate time.Time
		messageCount := 0
		byteCount := 0
		for {
			if ctx.Err() == nil {
				now := time.Now()
				if now.Sub(lastUpdate) >= time.Second {
					lastUpdate = now
					connectionManager.PublishStackData(index, connectionId, stackName, streamDirection, messageCount, byteCount)
				}
			}

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
				if reader, ok := i.V.(goprotoextra.ReadWriterSize); ok {
					messageCount++
					byteCount += reader.Size()
					nextFunc(ctx, reader)
					continue
				}
			}
		}
	}

	option := parseOptions(opts...)
	ctx := option.buildContext()
	go handler(ctx, o.Observe(opts...))
	return dispose
}

type mapOperatorWithCounters struct {
	mapOperator
	streamDirection   StreamDirection
	index             int
	connectionId      string
	stackName         string
	connectionManager IPublishToConnectionManager
	messageCount      int
	byteCount         int
	lastUpdate        time.Time
}

func newMapOperatorWithCounters(
	apply Func,
	index int,
	connectionId string,
	streamDirection StreamDirection,
	name string,
	connectionManager IPublishToConnectionManager) *mapOperatorWithCounters {
	return &mapOperatorWithCounters{
		mapOperator: mapOperator{
			apply: apply,
		},
		streamDirection:   streamDirection,
		index:             index,
		connectionId:      connectionId,
		stackName:         name,
		connectionManager: connectionManager,
	}
}

func (self *mapOperatorWithCounters) end(ctx context.Context, dst chan<- Item) {
	self.mapOperator.end(ctx, dst)
}

func (self *mapOperatorWithCounters) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	if ctx.Err() == nil {
		now := time.Now()
		if now.Sub(self.lastUpdate) >= time.Second {
			self.lastUpdate = now
			self.connectionManager.PublishStackData(self.index, self.connectionId, self.stackName, self.streamDirection, self.messageCount, self.byteCount)
		}
	}

	self.messageCount++
	if rws, ok := item.V.(goprotoextra.ReadWriterSize); ok {
		self.byteCount += rws.Size()
	}
	self.mapOperator.next(ctx, item, dst, operatorOptions)
}

func (o *ObservableImpl) MapInOutBound(
	index int,
	connectionId string,
	name string,
	streamDirection StreamDirection,
	connectionManager IPublishToConnectionManager,
	marshaller InOutBoundMarshaller,
	opts ...Option) Observable {
	if connectionManager == nil {
		connectionManager = nullPublisherInstance
	}
	return observable(
		o,
		func() operator {
			return newMapOperatorWithCounters(
				func(ctx context.Context, i interface{}) (interface{}, error) {
					if reader, ok := i.(goprotoextra.ReadWriterSize); ok {
						return marshaller(ctx, reader)
					}
					return nil, fmt.Errorf("type not a io.Reader")
				},
				index,
				connectionId,
				streamDirection,
				name,
				connectionManager)
		},
		false,
		true,
		opts...)
}
