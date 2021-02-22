package horizon

import (
	"fmt"
)

// Object represents a Horizon object, which can have a script attached or be
// referenced from other scripts.
type Object interface {
	Send(dst Object, eventName string, param interface{})
	// Set the object references for this object.
	Wire(Wiring)
	// Return the wired object
	Get(name string) (Object, bool)
	// TODO: add listen and connect
}

// NewObject creates an object for the given script and connects it to the event loop.
func NewObject(id string, script Script, el EventLoop) Object {
	return &object{
		id:        id,
		script:    script,
		eventLoop: el.(*eventLoop),
	}
}

// Script represents a behavior attached to a Horizon object, which is a
// collection of event handlers.
type Script func(self Object, e Event)

// Wiring represents a set of named Object pointers so objects can connect to each other.
type Wiring map[string]Object

// object represents any object in Horizon with a script attached.
type object struct {
	id        string
	script    Script
	wires     Wiring
	eventLoop *eventLoop
}

func (o *object) Send(dst Object, eventName string, param interface{}) {
	o.eventLoop.events <- Event{
		src:  o,
		dst:  dst.(*object),
		Name: eventName,
		Arg:  param,
	}
}

// Set the object references for this object.
func (o *object) Wire(w Wiring) { o.wires = w }

// Return the wired object. Returns false if object is not found.
func (o *object) Get(name string) (Object, bool) { v, ok := o.wires[name]; return v, ok }

// String returns the id of an object.
func (o *object) String() string { return fmt.Sprintf("{%s}", o.id) }
