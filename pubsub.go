package channels

import (
	"sync"
)

// Hub is publisher of PubSub.
type Hub chan<- interface{}

func makeonce(fn func()) func() {
	once := new(sync.Once)
	return func() {
		once.Do(fn)
	}
}

type subscribe struct {
	sub    chan<- interface{}
	closer func()
	done   chan func()
}

// Subscribe makes new subscriber of PubSub.
// subscriber channel is Nonblocking one, promise flag is same as Nonblocking parameter.
func (s Hub) Subscribe(promise bool) (<-chan interface{}, func()) {
	out, in := Nonblocking(promise)
	closer := makeonce(func() {
		close(in)
	})

	msg := subscribe{
		sub:    in,
		closer: closer,
		done:   make(chan func()),
	}
	s <- msg

	return out, <-msg.done
}

// PubSub makes Hub.
func PubSub() Hub {
	hub := make(chan interface{})

	go func(recv <-chan interface{}) {
		lock := make(chan struct{}, 1)

		id := 0
		subs := map[int]chan<- interface{}{}
		closers := map[int]func(){}

		for msg := range recv {
			if subscribe, ok := msg.(subscribe); ok {
				lock <- struct{}{}
				subid := id
				id++
				subs[subid] = subscribe.sub
				closers[subid] = subscribe.closer
				<-lock

				subscribe.done <- func() {
					lock <- struct{}{}
					delete(subs, subid)
					delete(closers, subid)
					<-lock
					subscribe.closer()
				}
				close(subscribe.done)

			} else {
				lock <- struct{}{}
				for _, sub := range subs {
					sub <- msg
				}
				<-lock
			}
		}

		lock <- struct{}{}
		for _, closer := range closers {
			closer()
		}
		<-lock
	}(hub)

	return hub
}
