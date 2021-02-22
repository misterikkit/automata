package horizon

import (
	"context"
	"fmt"
	"log"
)

// EventLoop is responsible for brokering events between all objects, using a
// shared event queue.
type EventLoop interface {
	// Run runs the main event loop.
	Run(context.Context)
}

// NewEventLoop returns an initialized EventLoop.
func NewEventLoop() EventLoop {
	return &eventLoop{
		events: make(chan Event, 50),
	}
}

type eventLoop struct {
	events chan Event
}

func (el *eventLoop) Run(ctx context.Context) {
	for {
		select {
		case e := <-el.events:
			log.Println(e)
			e.dst.script(e.dst, e)
		case <-ctx.Done():
			return
		}
	}
}

// Event mimics the data in a Horizon event
type Event struct {
	src, dst *object // src is for logging

	Name string
	Arg  interface{}
}

// String returns a debug representation of the event.
func (e Event) String() string {
	return fmt.Sprintf("%-20v %-14q -> %-20v (%v)", e.src, e.Name, e.dst, e.Arg)
}
