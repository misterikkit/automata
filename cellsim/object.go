package main

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// Icky global is a quick hack to block until all objects are wired.
var wireWG sync.WaitGroup

// Event mimics the data in a Horizon event
type Event struct {
	caller *Object // for logging
	Name   string
	Arg    interface{}
}

// Script represents a behavior attached to a Horizon object, which is a
// collection of event handlers.
type Script func(e Event, self *Object)

// Object represents any object in Horizon with a script attached.
type Object struct {
	id     string
	events chan Event
	script Script
}

// New creates an Object with the given id and script.
func New(id string, script Script) *Object {
	return &Object{
		id:     id,
		events: make(chan Event, 10),
		script: script,
	}
}

// String returns the id of an object.
func (o *Object) String() string { return fmt.Sprintf("{%s}", o.id) }

// convenience type for wiring "hack"
type wiring = map[string]*Object

// Wire is a shortcut for populating object relationships that would otherwise
// be wired or captured by triggers in Horizon. It works by sending a special
// event which must be interpreted by the script.
func (o *Object) Wire(m wiring) { wireWG.Add(1); o.Send(o, "_wire", m) }

// Run is an event loop for each object to process events.
func (o *Object) Run(ctx context.Context) {
	for {
		select {
		case e := <-o.events:
			if e.Name == "_wire" {
				wireWG.Done()
			}
			log.Printf("%-20v %-14q -> %-20v (%v)", e.caller, e.Name, o, e.Arg)
			o.script(e, o)
		case <-ctx.Done():
			return
		}
	}
}

// Send adds an event to this object's event queue.
func (o *Object) Send(caller *Object, name string, arg interface{}) {
	o.events <- Event{caller, name, arg}
}

// RunAll starts each object's event loop and blocks until the context is closed.
func RunAll(ctx context.Context, objs ...*Object) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var wg sync.WaitGroup
	for _, obj := range objs {
		wg.Add(1)
		go func(o *Object) {
			defer wg.Done()
			o.Run(ctx)
		}(obj)
	}
	wg.Wait()
}
