package rxgo

import (
	"sync"
	"time"
)

const Timeout = 500 * time.Millisecond
const PollingInterval = 20 * time.Millisecond

var s1 []interface{}
var s2 []interface{}
var m1 sync.Mutex
var m2 sync.Mutex

func InitTests() {
	s1 = make([]interface{}, 0)
	s2 = make([]interface{}, 0)
}

func Got1() []interface{} {
	m1.Lock()
	defer m1.Unlock()
	return s1
}

func Got2() []interface{} {
	m2.Lock()
	defer m2.Unlock()
	return s2
}

func Next1(i interface{}) {
	m1.Lock()
	defer m1.Unlock()
	s1 = append(s1, i)
}

func Next2(i interface{}) {
	m2.Lock()
	defer m2.Unlock()
	s2 = append(s2, i)
}
