package channels

import (
	"container/list"
)

// Nonblocking makes input/output channels.
// If promise flag is false, output channel will close as soon as after closing input channel evenif unread messages exists.
// otherwise, all messages will send to output channel after closing input channel, but needs to read all messages for terminate inner goroutin.
func Nonblocking(promise bool) (<-chan interface{}, chan<- interface{}) {
	type closed struct{}

	output := make(chan interface{})
	input := make(chan interface{})

	go func(send chan<- interface{}, recv <-chan interface{}) {
		defer close(send)

		queue := list.New()
		elm := queue.Front()
		val := interface{}(nil)
		pipe := chan<- interface{}(nil)

		for {
			select {
			case msg, ok := <-recv:
				if !ok {
					if !promise {
						return
					}
					recv = nil
					queue.PushBack(closed{})
				} else {
					queue.PushBack(msg)
				}

			case pipe <- val:
				queue.Remove(elm)
			}

			if elm = queue.Front(); elm != nil {
				if _, ok := elm.Value.(closed); ok {
					return
				}
				pipe, val = send, elm.Value
			} else {
				pipe, val = nil, nil
			}
		}
	}(output, input)

	return output, input
}
