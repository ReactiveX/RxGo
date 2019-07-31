package rxgo

import (
	"context"
	"runtime"
	"sync"
)

type workerPool struct {
	ctx   context.Context
	input chan<- task
}

type task struct {
	item   interface{}
	apply  Function
	output chan<- interface{}
}

var cpuPool workerPool = newWorkerPool(context.Background(), runtime.NumCPU())

func newWorkerPool(ctx context.Context, capacity int) workerPool {
	input := make(chan task, capacity)
	for i := 0; i < capacity; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case t := <-input:
					t.output <- t.apply(t.item)
				}
			}
		}()
	}
	return workerPool{
		ctx:   ctx,
		input: input,
	}
}

func (wp *workerPool) sendTask(item interface{}, apply Function, output chan<- interface{}, wg *sync.WaitGroup) {
	wg.Add(1)
	wp.input <- task{
		item:   item,
		apply:  apply,
		output: output,
	}
}

func (wp *workerPool) wait(f func(interface{}), output chan interface{}, wg *sync.WaitGroup) {
	go func() {
		for {
			select {
			case o := <-output:
				f(o)
				wg.Done()
			}
		}
	}()

	wg.Wait()
}
