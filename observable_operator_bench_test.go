package rxgo

import (
	"testing"
)

const (
	benchChannelCap       = 1000
	benchNumberOfElements = 1000000
)

func Benchmark_Sequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		disposed := Range(0, benchNumberOfElements, WithBufferedChannel(benchChannelCap)).
			Map(func(i interface{}) (interface{}, error) {
				return i, nil
			}).Run()
		b.StartTimer()
		<-disposed
	}
}

func Benchmark_Serialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		disposed := Range(0, benchNumberOfElements, WithBufferedChannel(benchChannelCap)).
			Map(func(i interface{}) (interface{}, error) {
				return i, nil
			}, WithCPUPool(), WithBufferedChannel(benchChannelCap)).Run()
		b.StartTimer()
		<-disposed
	}
}
