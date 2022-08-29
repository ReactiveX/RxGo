package rxgo

import (
	"log"
	"testing"
	"time"
)

func TestObservable(t *testing.T) {
	newObservable(func(obs Subscriber[int]) {
		for i := 1; i <= 5; i++ {
			time.Sleep(time.Second)
			obs.Next(i)
		}
		obs.Complete()
	}).SubscribeSync(func(i int) {
		log.Println(i)
	}, func(err error) {}, func() {})
}
