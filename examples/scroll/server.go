package main

import (
	"log"
	"strconv"
	"time"
	"net/http"

	"github.com/jochasinga/grx"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader {
	ReadBufferSize: 1024,
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

		source := grx.Just(p)
		source.Subscribe(&grx.Observer{
			NextHandler: grx.NextFunc(func(v interface{}) {
				offset, _ := strconv.Atoi(string(v.([]byte)))
				pchan <- offset
				lastPos = offset
			}),
		})
	}
	
}

func main() {
	http.HandleFunc("/scroll", indexHandler)
	http.Handle("/", http.FileServer(http.Dir(".")))
	
	log.Fatal(http.ListenAndServe(":4000", nil))
}
