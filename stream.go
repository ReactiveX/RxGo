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
