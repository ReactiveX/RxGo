package rxgo

import (
	"context"
	"testing"
	"time"
)

const (
	benchChannelCap            = 1000
	benchNumberOfElementsSmall = 1000
	ioPool                     = 32
)

func Benchmark_Range_Sequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		obs := Range(0, benchNumberOfElementsSmall, WithBufferedChannel(benchChannelCap)).
			Map(func(_ context.Context, i interface{}) (interface{}, error) {
				// Simulate a blocking IO call
				time.Sleep(5 * time.Millisecond)
				return i, nil
			})
		b.StartTimer()
		<-obs.Run()
	}
}

func Benchmark_Range_Serialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		obs := Range(0, benchNumberOfElementsSmall, WithBufferedChannel(benchChannelCap)).
			Map(func(_ context.Context, i interface{}) (interface{}, error) {
				// Simulate a blocking IO call
				time.Sleep(5 * time.Millisecond)
				return i, nil
			}, WithCPUPool(), WithBufferedChannel(benchChannelCap)).
			Serialize(0, func(i interface{}) int {
				return i.(int)
			})
		b.StartTimer()
		<-obs.Run()
	}
}

func Benchmark_Range_OptionSerialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		obs := Range(0, benchNumberOfElementsSmall, WithBufferedChannel(benchChannelCap)).
			Map(func(_ context.Context, i interface{}) (interface{}, error) {
				// Simulate a blocking IO call
				time.Sleep(5 * time.Millisecond)
				return i, nil
			}, WithCPUPool(), WithBufferedChannel(benchChannelCap), Serialize(func(i interface{}) int {
				return i.(int)
			}))
		b.StartTimer()
		<-obs.Run()
	}
}

func Benchmark_Reduce_Sequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		obs := Range(0, benchNumberOfElementsSmall, WithBufferedChannel(benchChannelCap)).
			Reduce(func(_ context.Context, acc, elem interface{}) (interface{}, error) {
				// Simulate a blocking IO call
				time.Sleep(5 * time.Millisecond)
				if a, ok := acc.(int); ok {
					if b, ok := elem.(int); ok {
						return a + b, nil
					}
				} else {
					return elem.(int), nil
				}
				return 0, errFoo
			})
		b.StartTimer()
		<-obs.Run()
	}
}

func Benchmark_Reduce_Parallel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		obs := Range(0, benchNumberOfElementsSmall, WithBufferedChannel(benchChannelCap)).
			Reduce(func(_ context.Context, acc, elem interface{}) (interface{}, error) {
				// Simulate a blocking IO call
				time.Sleep(5 * time.Millisecond)
				if a, ok := acc.(int); ok {
					if b, ok := elem.(int); ok {
						return a + b, nil
					}
				} else {
					return elem.(int), nil
				}
				return 0, errFoo
			}, WithPool(ioPool))
		b.StartTimer()
		<-obs.Run()
	}
}

func Benchmark_Map_Sequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		obs := Range(0, benchNumberOfElementsSmall, WithBufferedChannel(benchChannelCap)).
			Map(func(_ context.Context, i interface{}) (interface{}, error) {
				// Simulate a blocking IO call
				time.Sleep(5 * time.Millisecond)
				return i, nil
			})
		b.StartTimer()
		<-obs.Run()
	}
}

func Benchmark_Map_Parallel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		obs := Range(0, benchNumberOfElementsSmall, WithBufferedChannel(benchChannelCap)).
			Map(func(_ context.Context, i interface{}) (interface{}, error) {
				// Simulate a blocking IO call
				time.Sleep(5 * time.Millisecond)
				return i, nil
			}, WithCPUPool())
		b.StartTimer()
		<-obs.Run()
	}
}
