// Package onecontext provides a mechanism to merge multiple existing contexts.
package onecontext

import (
	"context"
	"sync"
	"time"
)

// Canceled is the error returned when the CancelFunc returned by Merge is called
type Canceled struct {
}

func (c *Canceled) Error() string {
	return "canceled context"
}

type onecontext struct {
	ctx        context.Context
	ctxs       []context.Context
	done       chan struct{}
	err        error
	errMutex   sync.Mutex
	cancelFunc context.CancelFunc
	cancelCtx  context.Context
}

// Merge allows to merge multiple contexts.
// It returns the merged context and a CancelFunc to cancel it.
func Merge(ctx context.Context, ctxs ...context.Context) (context.Context, context.CancelFunc) {
	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	o := &onecontext{
		done:       make(chan struct{}),
		ctx:        ctx,
		ctxs:       ctxs,
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
	}
	go o.run()
	return o, cancelFunc
}

func (o *onecontext) Deadline() (time.Time, bool) {
	min := time.Time{}

	if deadline, ok := o.ctx.Deadline(); ok {
		min = deadline
	}

	for _, ctx := range o.ctxs {
		if deadline, ok := ctx.Deadline(); ok {
			if min.IsZero() || deadline.Before(min) {
				min = deadline
			}
		}
	}

	return min, !min.IsZero()
}

func (o *onecontext) Done() <-chan struct{} {
	return o.done
}

func (o *onecontext) Err() error {
	o.errMutex.Lock()
	defer o.errMutex.Unlock()
	return o.err
}

func (o *onecontext) Value(key interface{}) interface{} {
	if value := o.ctx.Value(key); value != nil {
		return value
	}

	for _, ctx := range o.ctxs {
		if value := ctx.Value(key); value != nil {
			return value
		}
	}

	return nil
}

func (o *onecontext) run() {
	once := sync.Once{}

	if len(o.ctxs) == 1 {
		o.runTwoContexts(o.ctx, o.ctxs[0])
		return
	}

	o.runMultipleContexts(o.ctx, &once)
	for _, ctx := range o.ctxs {
		o.runMultipleContexts(ctx, &once)
	}
}

func (o *onecontext) cancel(err error) {
	o.cancelFunc()
	o.errMutex.Lock()
	o.err = err
	o.errMutex.Unlock()
	close(o.done)
}

func (o *onecontext) runTwoContexts(ctx1, ctx2 context.Context) {
	go func() {
		select {
		case <-o.cancelCtx.Done():
			o.cancel(&Canceled{})
		case <-ctx1.Done():
			o.cancel(ctx1.Err())
		case <-ctx2.Done():
			o.cancel(ctx2.Err())
		}
	}()
}

func (o *onecontext) runMultipleContexts(ctx context.Context, once *sync.Once) {
	go func() {
		select {
		case <-o.cancelCtx.Done():
			once.Do(func() {
				o.cancel(&Canceled{})
			})
		case <-ctx.Done():
			once.Do(func() {
				o.cancel(ctx.Err())
			})
		}
	}()
}
