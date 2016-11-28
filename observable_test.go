package grx

import (
	"errors"
	"fmt"
        "testing"
        "time"
        "net/http"

        "github.com/stretchr/testify/assert"
)

// TestCreateObservableWithConstructor tests if the constructor method returns an Observable
func TestCreateObservableWithConstructor(t *testing.T) {
        source := NewObservable()
        assert.IsType(t, (*Observable)(nil), source)
}

// TestCreateObservableWithEmpty tests if Empty creates an empty Observable that terminates right away.
func TestCreateObservableWithEmpty(t *testing.T) {
        msg := "Sumpin's"
        source := Empty()
        watcher := &Observer{
                DoneHandler: DoneFunc(func() {
                        msg += " brewin'"
                }),
        }
        source.Subscribe(watcher)

        // Block until side-effect is made
        <-time.After(100 * time.Millisecond)
        assert.Equal(t, "Sumpin's brewin'", msg)
}

// TestCreateObservableWithJust tests if Just method returns an Observable
func TestCreateObservableWithJustOneArg(t *testing.T) {
        url := "http://api.com/api/v1.0/user"
        source := Just(url)

        // Observable must have Observable type.
        assert.IsType(t, (*Observable)(nil), source)
}

func TestCreateObservableWithJustManyArgs(t *testing.T) {
	source := Just('R', 'x', 'G', 'o')
	assert.IsType(t, (*Observable)(nil), source)
}

// TestCreateObservableWithFrom tests if From methods returns an <*Observable>
func TestCreateObservableWithFrom(t *testing.T) {
        // Provided a slice of URL strings
        urls := []interface{}{
                "http://api.com/api/v1.0/user",
                "https://dropbox.com/api/v2.1/get/storage",
                "http://googleapi.com/map",
        }

        // and an Observable created with From method
        source := From(urls)

        // Observable must have Observable type
        assert.IsType(t, (*Observable)(nil), source)
}

func TestCreateObservableWithStart(t *testing.T) {
        // Provided a "directive" function which return an event that emits
        // a slice of integer numbers from 0 to 20.
        directive := func() interface{} {
                return 333
        }

        source := Start(directive)

        // Make sure it's the right type
        assert.IsType(t, (*Observable)(nil), source)

        nums := []int{}
        watcher := &Observer{
                NextHandler: NextFunc(func(v interface{}) {
                        if num, ok := v.(int); ok {
				nums = append(nums , num * 2)
                        }
                }),
                DoneHandler: DoneFunc(func() {
                        nums = append(nums, 666)
                }),
        }

        source.Subscribe(watcher)

        // Block until side-effect is made
        <-time.After(100 * time.Millisecond)
        assert.Exactly(t, []int{666, 666}, nums)
}

func TestSubscribeToJustObservableWithOneArg(t *testing.T) {
        urlWithUserID := ""

        // Provided an Observable created with Just method
        url := "http://api.com/api/v1.0/user"
        expected := url + "?id=999"
        source := Just(url)

        watcher := &Observer{
                NextHandler: NextFunc(func(v interface{}) {
                        if u, ok := v.(string); ok {
                                urlWithUserID = u
                        }
                }),
                DoneHandler: DoneFunc(func() {
                        urlWithUserID = urlWithUserID + "?id=999"
                }),
        }
        source.Subscribe(watcher)

        // Block until side-effect is made
        <-time.After(100 * time.Millisecond)
        assert.Exactly(t, expected, urlWithUserID)
}

func TestSubscribeToJustObservableWithMultipleArgs(t *testing.T) {
        nums := []int{}
	words := []string{}
	chars := []rune{}

        // Create an Observable with an integer.
	source := Just(1, "Hello", 'ห')

        // Create an Observer object.
        watcher := &Observer{

                // While there is more events down the stream (in this case, just one), add 1 to the value emitted.
                NextHandler: NextFunc(func(v interface{}) {
			switch item := v.(type) {
			case int:
				num := item + 1
				nums = append(nums, num)
			case string:
				word := item + " world"
				words = append(words, word)
			case rune:
				chars = append(chars, item)
			}
                        //num := v.(int) + 1
                        //nums = append(nums, num)
                }),

                // If an error is encountered at any time, panic.
                ErrHandler: ErrFunc(func(err error) {
                        panic(err.Error)
                }),

                // When the stream comes to an end, append 0 to the slice.
                DoneHandler: DoneFunc(func() {
                        nums = append(nums, 0)
                }),
        }

        // Start listening to stream
        source.Subscribe(watcher)

        // Block until side-effect is made
        <-time.After(100 * time.Millisecond)
	
        assert.Exactly(t, []int{2, 0}, nums)
	assert.Exactly(t, []string{"Hello world"}, words)
	assert.Exactly(t, []rune{'ห'}, chars)
}

func TestSubscribeToFromObservable(t *testing.T) {
        nums := []interface{}{1, 2, 3, 4, 5, 6}
        numCopy := []int{}

        // Create an Observable with an integer.
        source := From(nums)
        watcher := &Observer{
                // While there is more events down the stream (in this case, just one), add 1 to the value emitted.
                NextHandler: NextFunc(func(v interface{}) {
                        numCopy = append(numCopy, v.(int) + 1)
                }),

                // When the stream comes to an end, append 0 to the slice.
                DoneHandler: DoneFunc(func() {
                        numCopy = append(numCopy, 0)
                }),
        }

        // Start listening to stream
        source.Subscribe(watcher)

        // Block until side-effect is made
        <-time.After(100 * time.Millisecond)
        assert.Exactly(t, []int{2, 3, 4, 5, 6, 7, 0}, numCopy)
}

func TestStartMethodWithFakeExternalCalls(t *testing.T) {
        fakeResponses := []*http.Response{}

        // Fake directives that returns an Event containing an HTTP response.
        directive1 := func() interface{} {
                res := &http.Response{
                        Status: "404 NOT FOUND",
                        StatusCode: 404,
                        Proto: "HTTP/1.0",
                        ProtoMajor: 1,
                }
                time.Sleep(time.Millisecond * 20)
                return res
        }

	directive2 := func() interface{} {
                res := &http.Response{
                        Status: "200 OK",
                        StatusCode: 200,
                        Proto: "HTTP/1.0",
                        ProtoMajor: 1,
                }
                time.Sleep(10 * time.Millisecond)
                return res
        }

        directive3 := func() interface{} {
                res := &http.Response{
                        Status: "500 SERVER ERROR",
                        StatusCode: 500,
                        Proto: "HTTP/1.0",
                        ProtoMajor: 1,
                }
                time.Sleep(30 * time.Millisecond)
                return res
        }

        watcher := &Observer{
                NextHandler: NextFunc(func(v interface{}) {
                        fakeResponses = append(fakeResponses, v.(*http.Response))
                }),
                DoneHandler: DoneFunc(func() {
                        fakeResponses = append(fakeResponses, &http.Response{
                                Status: "999 End",
                                StatusCode: 999,
                        })
                }),
        }

        source := Start(directive1, directive2, directive3)
        source.Subscribe(watcher)

        // Block until side-effect is made
        <-time.After(100 * time.Millisecond)

        assert := assert.New(t)
        assert.IsType((*Observable)(nil), source)
        assert.Equal(4, len(fakeResponses))
        assert.Equal(200, fakeResponses[0].StatusCode)
        assert.Equal(404, fakeResponses[1].StatusCode)
        assert.Equal(500, fakeResponses[2].StatusCode)
        assert.Equal(999, fakeResponses[3].StatusCode)
}

func TestCreateObservableWithInterval(t *testing.T) {
        numch := make(chan int, 1)
	source := Interval(1 * time.Millisecond)
        go func() {
                source.Subscribe(&Observer{
                        NextHandler: NextFunc(func(v interface{}) {
                                numch <- v.(int)
                        }),
                })
        }()

        for i := 0; i <= 10; i++ {
                <-time.After(1 * time.Millisecond)
                assert.Equal(t, i, <-numch)
        }
}

func TestCreateObservableWithRange(t *testing.T) {
        nums := []int{}
        watcher := &Observer{
                NextHandler: NextFunc(func(v interface{}) {
                        nums = append(nums, v.(int))
                }),
	}
        _, _ = Range(1, 10).Subscribe(watcher)
        <-time.After(100 * time.Millisecond)
        assert.Exactly(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, nums)
}

func TestSubscriptionDoesNotBlock(t *testing.T) {
        source := Just("Hello")
        watcher:= &Observer{
                NextHandler: NextFunc(func(v interface{}) {
                        time.Sleep(1 * time.Second)
                        return
                }),
                DoneHandler: DoneFunc(func() {
                        t.Log("DONE")
                }),
        }
        first := time.Now()
        source.Subscribe(watcher)
        elapsed := time.Since(first)

        comp := assert.Comparison(func() bool {
                return elapsed < 1 * time.Second
        })
        assert.Condition(t, comp)
}

func TestSubscribeFuncToJustObservable(t *testing.T) {
	source := Just(1, '2', "yes", 9.012, 'ง', errors.New("Damn"), struct{}{}, func(){})
	source.SubscribeFunc(
		func(v interface{}) { fmt.Println(v) },
		func(err error) { fmt.Println(err) },
		nil,
	)
	<-time.After(100 * time.Millisecond)
}

func TestObservableIsDone(t *testing.T) {
	
	//o := Just(1, 2, 3, "hi", "bye", []float64{5.0, 10.2})
	o := Range(1, 10)
	done := make(chan struct{}, 1)
	o.Subscribe(&Observer{
		NextHandler: NextFunc(func(v interface{}) {
			t.Logf("Test value: %v", v)
		}),
		DoneHandler: DoneFunc(func() {
			done <- struct{}{}
		}),
	})

	<-time.After(100 * time.Millisecond)
	assert.Equal(t, struct{}{}, <-done)
	assert.Equal(t, true, o.isDone())
}

func TestCreateObservableAndSubscribeFunc(t *testing.T) {
	o := CreateObservable(func(ob *Observer) {
		ob.OnNext(99)
		ob.OnError(errors.New("Yike"))
		ob.OnDone()
	})

	o.SubscribeFunc(
		func(v interface{}) { assert.Equal(t, 99, v.(int)) },
		func(err error) { assert.NotNil(t, err) },
		func() { assert.NotNil(t, o.done) },
	)

	<-time.After(100 * time.Millisecond)
}

func TestCloseSubscription(t *testing.T) {
	nums := []int{}
	source := Interval(100 * time.Millisecond)
	sub, _ := source.Subscribe(&Observer{
		NextHandler: NextFunc(func(v interface{}) {
			nums = append(nums, v.(int))
		}),
	})
	
	if sub != nil {
		<-time.After(300 * time.Millisecond)
		sub.Close()
	}
	
	assert.Equal(t, []int{0, 1, 2}, nums)
}


