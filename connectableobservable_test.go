package rxgo

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/options"
	"testing"
)

func TestConnectableObservable(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Connectable Observable suite")
}

var _ = Describe("Connectable Observable", func() {
	BeforeEach(func() {
		InitTests()
	})

	Context("when creating two subscriptions to a connectable observable", func() {
		in := make(chan interface{}, 2)
		connectableObs := FromChannel(in).Publish()
		connectableObs.Subscribe(handlers.NextFunc(Next1), options.WithBufferBackpressureStrategy(2))
		connectableObs.Subscribe(handlers.NextFunc(Next2), options.WithBufferBackpressureStrategy(2))
		in <- 1
		in <- 2

		It("should not trigger the next handlers", func() {
			Eventually(Got1, Timeout, PollingInterval).Should(HaveLen(0))
			Eventually(Got2, Timeout, PollingInterval).Should(HaveLen(0))
		})

		Context("when connect is called", func() {
			It("should trigger the next handlers", func() {
				connectableObs.Connect()
				Eventually(Got1, Timeout, PollingInterval).Should(Equal([]interface{}{1, 2}))
				Eventually(Got2, Timeout, PollingInterval).Should(Equal([]interface{}{1, 2}))
			})
		})
	})

	Context("when creating a subscription to a connectable observable", func() {
		in := make(chan interface{}, 2)
		connectableObs := FromChannel(in).Publish()
		connectableObs.Subscribe(handlers.NextFunc(Next1), options.WithBufferBackpressureStrategy(2))

		Context("when connect is called", func() {
			It("should not be blocking", func() {
				connectableObs.Connect()
				in <- 1
				in <- 2
				Eventually(Got1, Timeout, PollingInterval).Should(Equal([]interface{}{1, 2}))
			})
		})
	})
})
