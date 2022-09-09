package rxgo

import (
	"log"
	"sync"
)

// Projects each source value to an Observable which is merged in the output Observable,
// emitting values only from the most recently projected Observable.
func SwitchMap[T any, R any](project func(value T, index uint) Observable[R]) OperatorFunc[T, R] {
	return func(source Observable[T]) Observable[R] {
		return newObservable(func(subscriber Subscriber[R]) {
			var (
				index uint
				wg    = new(sync.WaitGroup)
				// mu    = new(sync.RWMutex)
				stop = make(chan struct{})
				// closing  = make(chan struct{})
				upStream = source.SubscribeOn(wg.Done)
				// stream   Subscriber[R]
			)

			wg.Add(1)

			closeStream := func() {
				log.Println("Closing ---->")
				close(stop)
				stop = make(chan struct{})
			}

			startStream := func(obs Observable[R]) {
				log.Println("startStream --->")
				defer wg.Done()
				stream := obs.SubscribeOn()
				defer stream.Stop()

			loop:
				for {
					select {
					case <-stop:
						log.Println("STOP NOW")
						break loop

					case <-subscriber.Closed():
						stream.Stop()
						return

					case item, ok := <-stream.ForEach():
						if !ok {
							break loop
						}

						// log.Println(item)
						subscriber.Send() <- item
					}
				}

				log.Println("ENDED")
			}

		observe:
			for {
				select {
				case <-subscriber.Closed():
					upStream.Stop()
					closeStream()

				case item := <-upStream.ForEach():
					if err := item.Err(); err != nil {
						break observe
					}

					if item.Done() {
						break observe
					}

					closeStream()
					wg.Add(1)
					go startStream(project(item.Value(), index))
					index++
				}
			}

			wg.Wait()
			log.Println("Wait ended")
		})
	}
}

// Projects each source value to an Observable which is merged in the output
// Observable only if the previous projected Observable has completed.
func ExhaustMap[T any, R any](project func(value T, index uint) Observable[R]) OperatorFunc[T, R] {
	return func(source Observable[T]) Observable[R] {
		return newObservable(func(subscriber Subscriber[R]) {
			// var (
			// 	index        uint
			// 	isComplete   bool
			// 	subscription Subscription
			// )
			// source.SubscribeSync(
			// 	func(v T) {
			// 		if subscription == nil {
			// 			wg := new(sync.WaitGroup)
			// 			subscription = project(v, index).Subscribe(
			// 				func(v R) {
			// 					subscriber.Next(v)
			// 				},
			// 				func(error) {},
			// 				func() {
			// 					defer wg.Done()
			// 					subscription.Unsubscribe()
			// 					subscription = nil
			// 					if isComplete {
			// 						subscriber.Complete()
			// 					}
			// 				},
			// 			)
			// 			wg.Wait()
			// 		}
			// 		index++
			// 	},
			// 	subscriber.Error,
			// 	func() {
			// 		isComplete = true
			// 		if subscription == nil {
			// 			subscriber.Complete()
			// 		}
			// 	},
			// )

			// after collect the source
		})
	}
}

// Merge the values from all observables to a single observable result.
func Merge[T any](input Observable[T]) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				activeSubscription = 2
				wg                 = new(sync.WaitGroup)
				p1                 = source.SubscribeOn(wg.Done)
				p2                 = input.SubscribeOn(wg.Done)
				err                error
			)

			wg.Add(2)

			stopAll := func() {
				p1.Stop()
				p2.Stop()
				activeSubscription = 0
			}

			onNext := func(v Notification[T]) {
				if v == nil {
					return
				}

				// When any source errors, the resulting observable will error
				if err = v.Err(); err != nil {
					stopAll()
					subscriber.Send() <- Error[T](err)
					return
				}

				if v.Done() {
					activeSubscription--
					return
				}

				subscriber.Send() <- v
			}

			for activeSubscription > 0 {
				select {
				case <-subscriber.Closed():
					stopAll()
				case v1 := <-p1.ForEach():
					onNext(v1)
				case v2 := <-p2.ForEach():
					onNext(v2)
				}
			}

			// Wait for all input streams to unsubscribe
			wg.Wait()

			if err != nil {
				subscriber.Send() <- Error[T](err)
			} else {
				subscriber.Send() <- Complete[T]()
			}
		})
	}
}

// Creates an Observable that mirrors the first source Observable to emit a
// next, error or complete notification from the combination of the Observable
// to which the operator is applied and supplied Observables.
func RaceWith[T any](input Observable[T], inputs ...Observable[T]) OperatorFunc[T, T] {
	return func(source Observable[T]) Observable[T] {
		inputs = append([]Observable[T]{source, input}, inputs...)
		return newObservable(func(subscriber Subscriber[T]) {
			var (
				noOfInputs = len(inputs)
				// wg                  = new(sync.WaitGroup)
				fastestCh           = make(chan int, 1)
				activeSubscriptions = make([]Subscriber[T], noOfInputs)
				mu                  = new(sync.RWMutex)
				// unsubscribed        bool
			)
			// wg.Add(noOfInputs * 2)

			// unsubscribeAll := func(index int) {

			// 	var subscription Subscriber[T]

			// 	activeSubscriptions = []Subscriber[T]{subscription}
			// }

			// emit := func(index int, v Notification[T]) {
			// 	mu.RLock()
			// 	if unsubscribed {
			// 		mu.RUnlock()
			// 		return
			// 	}

			// 	log.Println("isThis", index)

			// 	mu.RUnlock()
			// 	mu.Lock()
			// 	unsubscribed = true
			// 	mu.Unlock()
			// 	// 	unsubscribeAll(index)

			// 	subscriber.Send() <- v
			// }

			for i, v := range inputs {
				activeSubscriptions[i] = v.SubscribeOn(func() {
					log.Println("DONE here")
					// wg.Done()
				})
				go func(idx int, obs Subscriber[T]) {
					// defer wg.Done()
					defer log.Println("closing routine", idx)

					for {
						select {
						case <-subscriber.Closed():
							log.Println("downstream closing ", idx)
							return
						case <-obs.Closed():
							log.Println("upstream closing ", idx)
							return
						case item := <-obs.ForEach():
							// mu.Lock()
							// defer mu.Unlock()
							// for _, sub := range activeSubscriptions {
							// 	sub.Stop()
							// }
							// activeSubscriptions = []Subscriber[T]{}
							log.Println("ForEach ah", idx, item)
							fastestCh <- idx
							// obs.Stop()
						}
					}
				}(i, activeSubscriptions[i])
			}

			log.Println("Fastest", <-fastestCh)
			mu.Lock()
			for _, v := range activeSubscriptions {
				v.Stop()
				log.Println(v)
			}
			mu.Unlock()
			// wg.Wait()

			log.Println("END")

			subscriber.Send() <- Complete[T]()
		})
	}
}
