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
		connectableObs.Subscribe(nextHandler(out1), options.WithBufferBackpressureStrategy(2))
		connectableObs.Subscribe(nextHandler(out2), options.WithBufferBackpressureStrategy(2))
		in <- 1
		in <- 2

		It("should not trigger the next handlers", func() {
			Expect(pollItem(out1, timeout)).Should(Equal(noData))
			Expect(pollItem(out2, timeout)).Should(Equal(noData))
		})

		Context("when connect is called", func() {
			It("should trigger the next handlers", func() {
				connectableObs.Connect()
				Expect(pollItem(out1, timeout)).Should(Equal(1))
				Expect(pollItem(out1, timeout)).Should(Equal(2))
				Expect(pollItem(out2, timeout)).Should(Equal(1))
				Expect(pollItem(out2, timeout)).Should(Equal(2))
			})
		})
	})

	Context("when creating a subscription to a connectable observable", func() {
		in := make(chan interface{}, 2)
		out := make(chan interface{}, 2)
		connectableObs := FromChannel(in).Publish()
		connectableObs.Subscribe(nextHandler(out), options.WithBufferBackpressureStrategy(2))

		Context("when connect is called", func() {
			It("should not be blocking", func() {
				connectableObs.Connect()
				in <- 1
				in <- 2
				Expect(pollItem(out, timeout)).Should(Equal(1))
				Expect(pollItem(out, timeout)).Should(Equal(2))
			})
		})
	})

	Context("when creating a subscription to a connectable observable", func() {
		Context("with back pressure strategy", func() {
			in := make(chan interface{}, 2)
			out := make(chan interface{}, 2)
			connectableObs := FromChannel(in).Publish()
			connectableObs.Subscribe(nextHandler(out), options.WithBufferBackpressureStrategy(3))
			connectableObs.Connect()
			in <- 1
			in <- 2
			in <- 3

			It("should buffer items", func() {
				Expect(pollItem(out, timeout)).Should(Equal(1))
				Expect(pollItem(out, timeout)).Should(Equal(2))
				Expect(pollItem(out, timeout)).Should(Equal(3))
				Expect(pollItem(out, timeout)).Should(Equal(noData))
			})
		})
	})

	Context("when creating a connectable observable", func() {
		connectableObservable := FromSlice([]interface{}{1, 2, 3, 5}).Publish()
		Context("when calling Map operator", func() {
			observable := connectableObservable.Map(func(i interface{}) interface{} {
				return i.(int) + 1
			})
			Context("when creating a subscription", func() {
				outNext, _, _ := subscribe(observable)
				Context("when calling Connect operator", func() {
					connectableObservable.Connect()
					It("the observer should received the produced items", func() {
						Expect(pollItems(outNext, timeout)).Should(Equal([]interface{}{
							2, 3, 4, 6,
						}))
					})
				})
			})
		})
	})
})
