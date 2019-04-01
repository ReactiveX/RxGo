package rxgo

import (
	"fmt"
	"github.com/reactivex/rxgo/errors"
	"github.com/reactivex/rxgo/optional"
	"sync"
	"time"
)

func model(iterable Iterable) ChanProducer {
	return func(out chan interface{}) {
	}
}

func all(iterable Iterable, predicate Predicate) ChanProducer {
	return func(out chan interface{}) {
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if !predicate(item) {
					out <- false
					close(out)
					return
				}
			} else {
				break
			}
		}
		out <- true
		close(out)
	}
}

func averageFloat32(iterable Iterable) ChanProducer {
	return func(out chan interface{}) {
		var sum float32 = 0
		var count float32 = 0
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if v, ok := item.(float32); ok {
					sum = sum + v
					count = count + 1
				} else {
					out <- errors.New(errors.IllegalInputError, fmt.Sprintf("type: %t", item))
					close(out)
					return
				}
			} else {
				break
			}
		}
		if count == 0 {
			out <- 0
		} else {
			out <- sum / count
		}
		close(out)
	}
}

func averageFloat64(iterable Iterable) ChanProducer {
	return func(out chan interface{}) {
		var sum float64 = 0
		var count float64 = 0
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if v, ok := item.(float64); ok {
					sum = sum + v
					count = count + 1
				} else {
					out <- errors.New(errors.IllegalInputError, fmt.Sprintf("type: %t", item))
					close(out)
					return
				}
			} else {
				break
			}
		}
		if count == 0 {
			out <- 0
		} else {
			out <- sum / count
		}
		close(out)
	}
}

func averageInt(iterable Iterable) ChanProducer {
	return func(out chan interface{}) {
		sum := 0
		count := 0
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if v, ok := item.(int); ok {
					sum = sum + v
					count = count + 1
				} else {
					out <- errors.New(errors.IllegalInputError, fmt.Sprintf("type: %t", item))
					close(out)
					return
				}
			} else {
				break
			}
		}
		if count == 0 {
			out <- 0
		} else {
			out <- sum / count
		}
		close(out)
	}
}

func averageInt8(iterable Iterable) ChanProducer {
	return func(out chan interface{}) {
		var sum int8 = 0
		var count int8 = 0
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if v, ok := item.(int8); ok {
					sum = sum + v
					count = count + 1
				} else {
					out <- errors.New(errors.IllegalInputError, fmt.Sprintf("type: %t", item))
					close(out)
					return
				}
			} else {
				break
			}
		}
		if count == 0 {
			out <- 0
		} else {
			out <- sum / count
		}
		close(out)
	}
}

func averageInt16(iterable Iterable) ChanProducer {
	return func(out chan interface{}) {
		var sum int16 = 0
		var count int16 = 0
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if v, ok := item.(int16); ok {
					sum = sum + v
					count = count + 1
				} else {
					out <- errors.New(errors.IllegalInputError, fmt.Sprintf("type: %t", item))
					close(out)
					return
				}
			} else {
				break
			}
		}
		if count == 0 {
			out <- 0
		} else {
			out <- sum / count
		}
		close(out)
	}
}

func averageInt32(iterable Iterable) ChanProducer {
	return func(out chan interface{}) {
		var sum int32 = 0
		var count int32 = 0
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if v, ok := item.(int32); ok {
					sum = sum + v
					count = count + 1
				} else {
					out <- errors.New(errors.IllegalInputError, fmt.Sprintf("type: %t", item))
					close(out)
					return
				}
			} else {
				break
			}
		}
		if count == 0 {
			out <- 0
		} else {
			out <- sum / count
		}
		close(out)
	}
}

func averageInt64(iterable Iterable) ChanProducer {
	return func(out chan interface{}) {
		var sum int64 = 0
		var count int64 = 0
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if v, ok := item.(int64); ok {
					sum = sum + v
					count = count + 1
				} else {
					out <- errors.New(errors.IllegalInputError, fmt.Sprintf("type: %t", item))
					close(out)
					return
				}
			} else {
				break
			}
		}
		if count == 0 {
			out <- 0
		} else {
			out <- sum / count
		}
		close(out)
	}
}

func bufferWithCount(iterable Iterable, count, skip int) ChanProducer {
	return func(out chan interface{}) {
		if count <= 0 {
			out <- errors.New(errors.IllegalInputError, "count must be positive")
			close(out)
			return
		}

		if skip <= 0 {
			out <- errors.New(errors.IllegalInputError, "skip must be positive")
			close(out)
			return
		}

		buffer := make([]interface{}, count, count)
		iCount := 0
		iSkip := 0
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				switch item := item.(type) {
				case error:
					if iCount != 0 {
						out <- buffer[:iCount]
					}
					out <- item
					close(out)
					return
				default:
					if iCount >= count { // Skip
						iSkip++
					} else { // Add to buffer
						buffer[iCount] = item
						iCount++
						iSkip++
					}

					if iSkip == skip { // Send current buffer
						out <- buffer
						buffer = make([]interface{}, count, count)
						iCount = 0
						iSkip = 0
					}
				}
			} else {
				break
			}
		}
		if iCount != 0 {
			out <- buffer[:iCount]
		}

		close(out)
	}
}

func bufferWithTime(iterable Iterable, timespan, timeshift Duration) ChanProducer {
	return func(out chan interface{}) {
		if timespan == nil || timespan.duration() == 0 {
			out <- errors.New(errors.IllegalInputError, "timespan must not be nil")
			close(out)
			return
		}

		if timeshift == nil {
			timeshift = WithDuration(0)
		}

		var mux sync.Mutex
		var listenMutex sync.Mutex
		buffer := make([]interface{}, 0)
		stop := false
		listen := true

		// First goroutine in charge to check the timespan
		go func() {
			for {
				time.Sleep(timespan.duration())
				mux.Lock()
				if !stop {
					out <- buffer
					buffer = make([]interface{}, 0)
					mux.Unlock()

					if timeshift.duration() != 0 {
						listenMutex.Lock()
						listen = false
						listenMutex.Unlock()
						time.Sleep(timeshift.duration())
						listenMutex.Lock()
						listen = true
						listenMutex.Unlock()
					}
				} else {
					mux.Unlock()
					return
				}
			}
		}()

		// Second goroutine in charge to retrieve the items from the source observable
		go func() {
			it := iterable.Iterator()
			for {
				if item, err := it.Next(); err == nil {
					switch item := item.(type) {
					case error:
						mux.Lock()
						if len(buffer) > 0 {
							out <- buffer
						}
						out <- item
						close(out)
						stop = true
						mux.Unlock()
						return
					default:
						listenMutex.Lock()
						l := listen
						listenMutex.Unlock()

						mux.Lock()
						if l {
							buffer = append(buffer, item)
						}
						mux.Unlock()
					}
				} else {
					break
				}
			}
			mux.Lock()
			if len(buffer) > 0 {
				out <- buffer
			}
			close(out)
			stop = true
			mux.Unlock()
		}()
	}
}

func bufferWithTimeOrCount(iterable Iterable, timespan Duration, count int) ChanProducer {
	return func(out chan interface{}) {
		if timespan == nil || timespan.duration() == 0 {
			out <- errors.New(errors.IllegalInputError, "timespan must not be nil")
			close(out)
			return
		}

		if count <= 0 {
			out <- errors.New(errors.IllegalInputError, "count must be positive")
			close(out)
			return
		}

		sendCh := make(chan []interface{})
		errCh := make(chan error)
		buffer := make([]interface{}, 0)
		var bufferMutex sync.Mutex

		// First sender goroutine
		go func() {
			for {
				select {
				case currentBuffer := <-sendCh:
					out <- currentBuffer
				case error := <-errCh:
					if len(buffer) > 0 {
						out <- buffer
					}
					if error != nil {
						out <- error
					}
					close(out)
					return
				case <-time.After(timespan.duration()): // Send on timer
					bufferMutex.Lock()
					b := make([]interface{}, len(buffer))
					copy(b, buffer)
					buffer = make([]interface{}, 0)
					bufferMutex.Unlock()

					out <- b
				}
			}
		}()

		// Second goroutine in charge to retrieve the items from the source observable
		go func() {
			it := iterable.Iterator()
			for {
				if item, err := it.Next(); err == nil {
					switch item := item.(type) {
					case error:
						errCh <- item
						return
					default:
						bufferMutex.Lock()
						buffer = append(buffer, item)
						if len(buffer) >= count {
							b := make([]interface{}, len(buffer))
							copy(b, buffer)
							buffer = make([]interface{}, 0)
							bufferMutex.Unlock()

							sendCh <- b
						} else {
							bufferMutex.Unlock()
						}
					}
				} else {
					break
				}
			}
			errCh <- nil
		}()
	}
}

func contains(iterable Iterable, equal Predicate) ChanProducer {
	return func(out chan interface{}) {
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if equal(item) {
					out <- true
					close(out)
					return
				}
			} else {
				break
			}
		}
		out <- false
		close(out)
	}
}

func count(iterable Iterable) ChanProducer {
	return func(out chan interface{}) {
		var count int64
		it := iterable.Iterator()
		for {
			if _, err := it.Next(); err == nil {
				count++
			} else {
				break
			}
		}
		out <- count
		close(out)
	}
}

func defaultIfEmpty(iterable Iterable, defaultValue interface{}) ChanProducer {
	return func(out chan interface{}) {
		empty := true
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				empty = false
				out <- item
			} else {
				break
			}
		}
		if empty {
			out <- defaultValue
		}
		close(out)
	}
}

func distinct(iterable Iterable, apply Function) ChanProducer {
	return func(out chan interface{}) {
		keysets := make(map[interface{}]struct{})
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				key := apply(item)
				_, ok := keysets[key]
				if !ok {
					out <- item
				}
				keysets[key] = struct{}{}
			} else {
				break
			}
		}
		close(out)
	}
}

func distinctUntilChanged(iterable Iterable, apply Function) ChanProducer {
	return func(out chan interface{}) {
		var current interface{}
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				key := apply(item)
				if current != key {
					out <- item
					current = key
				}
			} else {
				break
			}
		}
		close(out)
	}
}

func doOnEach(iterable Iterable, onNotification Consumer) ChanProducer {
	return func(out chan interface{}) {
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				out <- item
				onNotification(item)
			} else {
				break
			}
		}
		close(out)
	}
}

func elementAt(iterable Iterable, index uint) ChanProducer {
	return func(out chan interface{}) {
		takeCount := 0

		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if takeCount == int(index) {
					out <- item
					close(out)
					return
				}
				takeCount += 1
			} else {
				break
			}
		}
		out <- errors.New(errors.ElementAtError)
		close(out)
	}
}

func filter(iterable Iterable, apply Predicate) ChanProducer {
	return func(out chan interface{}) {
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if apply(item) {
					out <- item
				}
			} else {
				break
			}
		}
		close(out)
	}
}

func first(iterable Iterable) ChanProducer {
	return func(out chan interface{}) {
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				out <- item
				break
			} else {
				break
			}
		}
		close(out)
	}
}

func firstOrDefault(iterable Iterable, defaultValue interface{}) ChanProducer {
	return func(out chan interface{}) {
		first := defaultValue
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				first = item
				break
			} else {
				break
			}
		}
		out <- first
		close(out)
	}
}

func last(iterable Iterable) ChanProducer {
	return func(out chan interface{}) {
		var last interface{}
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				last = item
			} else {
				break
			}
		}
		out <- last
		close(out)
	}
}

func lastOrDefault(iterable Iterable, defaultValue interface{}) ChanProducer {
	return func(out chan interface{}) {
		last := defaultValue
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				last = item
			} else {
				break
			}
		}
		out <- last
		close(out)
	}
}

func mapFromFunction(iterable Iterable, apply Function) ChanProducer {
	return func(out chan interface{}) {
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				out <- apply(item)
			} else {
				break
			}
		}
		close(out)
	}
}

func max(iterable Iterable, comparator Comparator) ChanProducer {
	return func(out chan interface{}) {
		empty := true
		var max interface{} = nil
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				empty = false

				if max == nil {
					max = item
				} else {
					if comparator(max, item) == Smaller {
						max = item
					}
				}
			} else {
				break
			}
		}
		if empty {
			out <- optional.Empty()
		} else {
			out <- optional.Of(max)
		}
		close(out)
	}
}

func min(iterable Iterable, comparator Comparator) ChanProducer {
	return func(out chan interface{}) {
		empty := true
		var min interface{} = nil
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				empty = false

				if min == nil {
					min = item
				} else {
					if comparator(min, item) == Greater {
						min = item
					}
				}
			} else {
				break
			}
		}
		if empty {
			out <- optional.Empty()
		} else {
			out <- optional.Of(min)
		}
		close(out)
	}
}

func reduce(iterable Iterable, apply Function2) ChanProducer {
	return func(out chan interface{}) {
		var acc interface{}
		empty := true
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				empty = false
				acc = apply(acc, item)
			} else {
				break
			}
		}
		if empty {
			out <- optional.Empty()
		} else {
			out <- optional.Of(acc)
		}
		close(out)
	}
}

func repeat(iterable Iterable, count int64, frequency Duration) ChanProducer {
	return func(out chan interface{}) {
		if count != Indefinitely {
			if count < 0 {
				count = 0
			}
		}

		persist := make([]interface{}, 0)
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				out <- item
				persist = append(persist, item)
			} else {
				break
			}
		}
		for {
			if count != Indefinitely {
				if count == 0 {
					break
				}
			}

			if frequency != nil {
				time.Sleep(frequency.duration())
			}

			for _, v := range persist {
				out <- v
			}

			count = count - 1
		}
		close(out)
	}
}

func scan(iterable Iterable, apply Function2) ChanProducer {
	return func(out chan interface{}) {
		var current interface{}
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				tmp := apply(current, item)
				out <- tmp
				current = tmp
			} else {
				break
			}
		}
		close(out)
	}
}

func skip(iterable Iterable, nth uint) ChanProducer {
	return func(out chan interface{}) {
		skipCount := 0
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if skipCount < int(nth) {
					skipCount += 1
					continue
				}
				out <- item
			} else {
				break
			}
		}
		close(out)
	}
}

func skipLast(iterable Iterable, nth uint) ChanProducer {
	return func(out chan interface{}) {
		buf := make(chan interface{}, nth)
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				select {
				case buf <- item:
				default:
					out <- (<-buf)
					buf <- item
				}
			} else {
				break
			}
		}
		close(buf)
		close(out)
	}
}

func skipWhile(iterable Iterable, apply Predicate) ChanProducer {
	return func(out chan interface{}) {
		skip := true
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if !skip {
					out <- item
				} else {
					if !apply(item) {
						out <- item
						skip = false
					}
				}
			} else {
				break
			}
		}
		close(out)
	}
}

func startWithItems(iterable Iterable, items ...interface{}) ChanProducer {
	return func(out chan interface{}) {
		for _, item := range items {
			out <- item
		}

		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				out <- item
			} else {
				break
			}
		}

		close(out)
	}
}

func startWithIterable(iterable Iterable, startIterable Iterable) ChanProducer {
	return func(out chan interface{}) {
		it := startIterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				out <- item
			} else {
				break
			}
		}

		it = iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				out <- item
			} else {
				break
			}
		}

		close(out)
	}
}

func startWithObservable(iterable Iterable, obs Observable) ChanProducer {
	return func(out chan interface{}) {
		it := obs.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				out <- item
			} else {
				break
			}
		}

		it = iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				out <- item
			} else {
				break
			}
		}

		close(out)
	}
}

func sumFloat32(iterable Iterable) ChanProducer {
	return func(out chan interface{}) {
		var sum float32
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				switch item := item.(type) {
				case int:
					sum = sum + float32(item)
				case int8:
					sum = sum + float32(item)
				case int16:
					sum = sum + float32(item)
				case int32:
					sum = sum + float32(item)
				case int64:
					sum = sum + float32(item)
				case float32:
					sum = sum + item
				default:
					out <- errors.New(errors.IllegalInputError,
						fmt.Sprintf("expected type: float32, int, int8, int16, int32 or int64, got %t", item))
					close(out)
					return
				}
			} else {
				break
			}
		}
		out <- sum
		close(out)
	}
}

func sumFloat64(iterable Iterable) ChanProducer {
	return func(out chan interface{}) {
		var sum float64
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				switch item := item.(type) {
				case int:
					sum = sum + float64(item)
				case int8:
					sum = sum + float64(item)
				case int16:
					sum = sum + float64(item)
				case int32:
					sum = sum + float64(item)
				case int64:
					sum = sum + float64(item)
				case float32:
					sum = sum + float64(item)
				case float64:
					sum = sum + item
				default:
					out <- errors.New(errors.IllegalInputError,
						fmt.Sprintf("expected type: float32, float64, int, int8, int16, int32 or int64, got %t", item))
					close(out)
					return
				}
			} else {
				break
			}
		}
		out <- sum
		close(out)
	}
}

func sumInt64(iterable Iterable) ChanProducer {
	return func(out chan interface{}) {
		var sum int64
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				switch item := item.(type) {
				case int:
					sum = sum + int64(item)
				case int8:
					sum = sum + int64(item)
				case int16:
					sum = sum + int64(item)
				case int32:
					sum = sum + int64(item)
				case int64:
					sum = sum + item
				default:
					out <- errors.New(errors.IllegalInputError,
						fmt.Sprintf("expected type: int, int8, int16, int32 or int64, got %t", item))
					close(out)
					return
				}
			} else {
				break
			}
		}
		out <- sum
		close(out)
	}
}

func take(iterable Iterable, nth uint) ChanProducer {
	return func(out chan interface{}) {
		takeCount := 0
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if takeCount < int(nth) {
					takeCount += 1
					out <- item
					continue
				}
				break
			} else {
				break
			}
		}
		close(out)
	}
}

func takeLast(iterable Iterable, nth uint) ChanProducer {
	return func(out chan interface{}) {
		buf := make([]interface{}, nth)
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if len(buf) >= int(nth) {
					buf = buf[1:]
				}
				buf = append(buf, item)
			} else {
				break
			}
		}
		for _, takenItem := range buf {
			out <- takenItem
		}
		close(out)
	}
}

func takeWhile(iterable Iterable, apply Predicate) ChanProducer {
	return func(out chan interface{}) {
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				if apply(item) {
					out <- item
					continue
				}
				break
			} else {
				break
			}
		}
		close(out)
	}
}

func toList(iterable Iterable) ChanProducer {
	return func(out chan interface{}) {
		s := make([]interface{}, 0)
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				s = append(s, item)
			} else {
				break
			}
		}
		out <- s
		close(out)
	}
}

func toMap(iterable Iterable, keySelector Function) ChanProducer {
	return func(out chan interface{}) {
		m := make(map[interface{}]interface{})
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				m[keySelector(item)] = item
			} else {
				break
			}
		}
		out <- m
		close(out)
	}
}

func toMapWithValueSelector(iterable Iterable, keySelector Function, valueSelector Function) ChanProducer {
	return func(out chan interface{}) {
		m := make(map[interface{}]interface{})
		it := iterable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				m[keySelector(item)] = valueSelector(item)
			} else {
				break
			}
		}
		out <- m
		close(out)
	}
}

func zipFromObservable(iterable Iterable, publisher Observable, zipper Function2) ChanProducer {
	return func(out chan interface{}) {
		it := iterable.Iterator()
		it2 := publisher.Iterator()
	OuterLoop:
		for {
			if item1, err := it.Next(); err == nil {
				for {
					if item2, err := it2.Next(); err == nil {
						out <- zipper(item1, item2)
						continue OuterLoop
					} else {
						break
					}
				}
				break OuterLoop
			} else {
				break
			}
		}
		close(out)
	}
}
