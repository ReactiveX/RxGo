package rxgo

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/reactivex/rxgo/options"
)

var _ = Describe("Connectable Observable", func() {
	Context("when creating two subscriptions to a connectable observable", func() {
		in := make(chan interface{}, 2)
		out1 := make(chan interface{}, 2)
		out2 := make(chan interface{}, 2)
		connectableObs := FromChannel(in).Publish()
		connectableObs.Subscribe(next(out1), options.WithBufferBackpressureStrategy(2))
		connectableObs.Subscribe(next(out2), options.WithBufferBackpressureStrategy(2))
		in <- 1
		in <- 2

		It("should not trigger the next handlers", func() {
			Expect(get(out1, timeout)).Should(Equal(noData))
			Expect(get(out2, timeout)).Should(Equal(noData))
		})

		Context("when connect is called", func() {
			It("should trigger the next handlers", func() {
				connectableObs.Connect()
				Expect(get(out1, timeout)).Should(Equal(1))
				Expect(get(out1, timeout)).Should(Equal(2))
				Expect(get(out2, timeout)).Should(Equal(1))
				Expect(get(out2, timeout)).Should(Equal(2))
			})
		})
	})

	Context("when creating a subscription to a connectable observable", func() {
		in := make(chan interface{}, 2)
		out := make(chan interface{}, 2)
		connectableObs := FromChannel(in).Publish()
		connectableObs.Subscribe(next(out), options.WithBufferBackpressureStrategy(2))

		Context("when connect is called", func() {
			It("should not be blocking", func() {
				connectableObs.Connect()
				in <- 1
				in <- 2
				Expect(get(out, timeout)).Should(Equal(1))
				Expect(get(out, timeout)).Should(Equal(2))
			})
		})
	})
})
