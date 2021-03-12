package horizon

import (
	"context"
	"fmt"
	"log"
	"strings"
)

// EventLoop is responsible for brokering events between all objects, using a
// shared event queue.
type EventLoop interface {
	// Run runs the main event loop.
	Run(context.Context)
	Diagram() string // TODO: this belongs elsewhere
}

// NewEventLoop returns an initialized EventLoop.
func NewEventLoop() EventLoop {
	return &eventLoop{
		// TODO: NewObject writes to this channel before Run() is called, so this can deadlock.
		events: make(chan Event, 50),
	}
}

type eventLoop struct {
	events chan Event
	log    []Event // for diagrams
}

func (el *eventLoop) Run(ctx context.Context) {
	for {
		select {
		case e := <-el.events:
			log.Println(e)
			el.log = append(el.log, e)
			e.dst.script(e.dst, e)
		case <-ctx.Done():
			return
		}
	}
}

func (el *eventLoop) Diagram() string {
	es := make([]string, 0, len(el.log))
	objFmt := func(v interface{}) string {
		if o, ok := v.(*object); ok {
			return o.id
		}
		return fmt.Sprint(v)
	}
	for _, e := range el.log {
		es = append(es, fmt.Sprintf("%s -> %s: %s(%v)", objFmt(e.src), objFmt(e.dst), e.Name, objFmt(e.Arg)))
	}
	return strings.Join(es, "\n")
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
