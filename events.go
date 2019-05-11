package pubsub

import (
	"sync"
	"time"
)

type Actor struct {
	ID         string
	Attributes map[string]string
}

type Msg struct {
	// Scope project company contract etc...
	Scope string `json:"scope,omitempty"`

	// Action create update delete
	Action string
	Actor  Actor

	Time     int64 `json:"time,omitempty"`
	TimeNano int64 `json:"time_nano,omitempty"`
}

type Events struct {
	mu     sync.Mutex
	events []Msg
	pub    *Publisher
}

const (
	eventsLimit = 256
	bufferSize  = 1024
)

// New returns new *Events instance
func NewEvents() *Events {
	return &Events{
		events: make([]Msg, 0, eventsLimit),
		pub:    NewPublisher(100*time.Millisecond, bufferSize),
	}
}

// Subscribe adds new listener to events,
func (e *Events) Subscribe() ([]Msg, chan interface{}) {
	e.mu.Lock()
	current := make([]Msg, len(e.events))
	copy(current, e.events)
	l := e.pub.Subscribe()
	e.mu.Unlock()
	return current, l
}

// SubscribeTopic adds new listener to events,
func (e *Events) SubscribeTopic(ef *Filter) ([]Msg, chan interface{}) {
	e.mu.Lock()

	var topic func(m interface{}) bool
	if ef != nil && ef.Len() > 0 {
		topic = func(m interface{}) bool { return ef.Include(m.(Msg)) }
	}

	var buffered []Msg
	for i := len(e.events) - 1; i >= 0; i-- {
		ev := e.events[i]
		if topic == nil || topic(ev) {
			buffered = append([]Msg{ev}, buffered...)
		}
	}

	var ch chan interface{}
	if topic != nil {
		ch = e.pub.SubscribeTopic(topic)
	} else {
		// Subscribe to all events if there are no filters
		ch = e.pub.Subscribe()
	}

	e.mu.Unlock()
	return buffered, ch
}

// PublishMessage broadcasts event to listeners. Each listener has 100 milliseconds to
// receive the event or it will be skipped.
func (e *Events) PublishMessage(jm Msg) {
	e.mu.Lock()
	if len(e.events) == cap(e.events) {
		// discard oldest event
		copy(e.events, e.events[1:])
		e.events[len(e.events)-1] = jm
	} else {
		e.events = append(e.events, jm)
	}
	e.mu.Unlock()
	e.pub.Publish(jm)
}

// Evict evicts listener from pubsub
func (e *Events) Evict(l chan interface{}) {
	e.pub.Evict(l)
}

// SubscribersCount returns number of event listeners
func (e *Events) SubscribersCount() int {
	return e.pub.Len()
}
