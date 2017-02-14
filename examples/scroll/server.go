/*
 * Copyright (c) 2016-2017 Joe Chasinga
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/observable"
	"github.com/reactivex/rxgo/observer"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	lastPos := 0
	pchan := make(chan interface{})

	go func() {
		for {
			select {
			case <-time.After(2 * time.Second):
				log.Printf("User has stopped @%dpx\n", lastPos)
			case pos := <-pchan:
				log.Printf("User @%dpx\n", pos.(int))
			}
		}
	}()

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			return
		}

		source := observable.Just(p)

		onNext := handlers.NextFunc(func(item interface{}) {
			if item, ok := item.([]byte); ok {
				offset, _ := strconv.Atoi(string(item))
				pchan <- offset
				lastPos = offset
			}
		})

		_ = source.Subscribe(observer.New(onNext))
	}

}

func main() {
	http.HandleFunc("/scroll", indexHandler)
	http.Handle("/", http.FileServer(http.Dir(".")))

	log.Println("Server listening on port 4000")
	log.Fatal(http.ListenAndServe(":4000", nil))
}
