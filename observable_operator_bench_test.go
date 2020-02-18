package rxgo

import (
	"testing"
)

const (
	benchChannelCap       = 1_000
	benchNumberOfElements = 1_000_000
)

func producer(ch chan Item, n int) {
	for i := 0; i < n; i++ {
		Of(i).SendBlocking(ch)
	}
	close(ch)
}

func Benchmark_Sequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ch := make(chan Item, benchChannelCap)
		go producer(ch, benchNumberOfElements)
		<-FromChannel(ch).Map(func(i interface{}) (interface{}, error) {
			return i, nil
		}).Run()
	}
}

func Benchmark_Serialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ch := make(chan Item, benchChannelCap)
		go producer(ch, benchNumberOfElements)
		<-FromChannel(ch).Map(func(i interface{}) (interface{}, error) {
			return i, nil
		}, WithCPUPool(), WithBufferedChannel(benchChannelCap)).Run()
	}
}
