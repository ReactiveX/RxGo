package tooth

import "sync"

// Tooth is an entity for publishing and receiving messages.
type Tooth struct {
	Name       string
	Dupls      map[string]chan string
	Subscribed []*Tooth
}

// New is a constructor that returns a new Tooth.
func New(name string) *Tooth {
	return &Tooth{
		Name:  name,
		Dupls: make(map[string]chan string),
	}
}

// Publish publishes a message to all of the subscribers.
func (th *Tooth) Publish(msg string) {
	var wg sync.WaitGroup
	output := func(msg string, ch chan string) {
		ch <- msg
		wg.Done()
	}
	wg.Add(len(th.Dupls))
	for _, ch := range th.Dupls {
		go output(msg, ch)
	}
	go func() {
		wg.Wait()
		for _, ch := range th.Dupls {
			close(ch)
		}
	}()
}

// Subscribe subscribes to another Tooth's published message(s).
func (th *Tooth) Subscribe(to *Tooth) {
	if to.Dupls[th.Name] == nil {
		to.Dupls[th.Name] = make(chan string)
	}
	th.Subscribed = append(th.Subscribed, to)
}

// FetchAll returns a slice of string containing all messages from subscribed Tooths.
func (th *Tooth) FetchAll() []string {
	temp := make(chan string)
	messages := []string{}
	go func(c chan string) {
		for _, sub := range th.Subscribed {
			for msg := range sub.Dupls[th.Name] {
				c <- msg
			}
		}
		close(c)
	}(temp)
	for msg := range temp {
		messages = append(messages, msg)
	}
	return messages
}

// Fetch returns a message from a specific Tooth.
func (th *Tooth) Fetch(from *Tooth) string {
	temp := make(chan string)
	go func(c chan string) {
		for msg := range from.Dupls[th.Name] {
			c <- msg
		}
		close(c)
	}(temp)
	return <-temp
}
