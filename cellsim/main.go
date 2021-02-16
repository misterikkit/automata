package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Event struct {
	caller *Object // for logging
	Name   string
	Arg    interface{}
}

type Script func(e Event, self *Object)
type Object struct {
	id string
	// listeners []*Object
	events chan Event
	script Script
}

func New(id string, script Script) *Object {
	return &Object{
		id:     id,
		events: make(chan Event, 10),
		script: script,
	}
}

func (o *Object) String() string { return fmt.Sprintf("{%s}", o.id) }

type wiring = map[string]*Object

func (o *Object) Wire(m wiring) { o.Send(nil, "_wire", m) }

func (o *Object) Run(ctx context.Context) {
	for {
		select {
		case e := <-o.events:
			log.Printf("%-14v %-14q -> %-14v (%v)", e.caller, e.Name, o, e.Arg)
			o.script(e, o) // TODO: per-object lock?
			// for _, l := range o.listeners {
			// 	l.events <- e
			// }
		case <-ctx.Done():
			return
		}
	}
}

// func (o *Object) Listen(target *Object) { o.listeners = append(o.listeners, target) }

func (o *Object) Send(caller *Object, name string, arg interface{}) {
	o.events <- Event{caller, name, arg}
}

func main() {
	rand.Seed(time.Now().Unix())
	cell1 := New("cell1", Cell())
	c1Probes := []*Object{
		New("cell1-probe1", Probe()),
		New("cell1-probe2", Probe()),
		New("cell1-probe3", Probe()),
		New("cell1-probe4", Probe()),
	}

	cell2 := New("cell2", Cell())
	c2Probes := []*Object{
		New("cell2-probe1", Probe()),
		New("cell2-probe2", Probe()),
		New("cell2-probe3", Probe()),
		New("cell2-probe4", Probe()),
	}

	wall1 := New("wall[1,2]", Wall())

	border := New("border", Terminator())

	cell1.Wire(wiring{"probe": c1Probes[0]})
	cell2.Wire(wiring{"probe": c2Probes[0]})

	c1Probes[0].Wire(wiring{"cell": cell1, "next": c1Probes[1], "wall": wall1})
	c1Probes[1].Wire(wiring{"cell": cell1, "next": c1Probes[2], "wall": border})
	c1Probes[2].Wire(wiring{"cell": cell1, "next": c1Probes[3], "wall": border})
	c1Probes[3].Wire(wiring{"cell": cell1, "next": c1Probes[0], "wall": border})

	c2Probes[0].Wire(wiring{"cell": cell2, "next": c2Probes[1], "wall": wall1})
	c2Probes[1].Wire(wiring{"cell": cell2, "next": c2Probes[2], "wall": border})
	c2Probes[2].Wire(wiring{"cell": cell2, "next": c2Probes[3], "wall": border})
	c2Probes[3].Wire(wiring{"cell": cell2, "next": c2Probes[0], "wall": border})

	wall1.Wire(wiring{"probe1": c1Probes[0], "probe2": c2Probes[0]})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	init := New("init", func(Event, *Object) { cancel() })

	// make sure all _wire events are done
	time.AfterFunc(500*time.Millisecond, func() { cell1.Send(init, "visit", init) })
	// start the engine
	runAll(ctx,
		cell1,
		c1Probes[0],
		c1Probes[1],
		c1Probes[2],
		c1Probes[3],
		cell2,
		c2Probes[0],
		c2Probes[1],
		c2Probes[2],
		c2Probes[3],
		wall1,
		init,
		border,
	)
}

func runAll(ctx context.Context, objs ...*Object) {
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

func Cell() Script {
	// wired
	var (
		probe *Object
	)
	// variables
	var (
		visited = false
		back    *Object
	)
	return func(e Event, self *Object) {
		switch e.Name {
		case "_wire":
			w := e.Arg.(wiring)
			probe = w["probe"]
		case "visit":
			obj := e.Arg.(*Object)
			back = obj
			visited = true
			probe.Send(self, "visitRand", rand.Int()%4)
		case "check":
			obj := e.Arg.(*Object)
			obj.Send(self, "checkResult", visited)
		case "deadEnd":
			// pop stack
			back.Send(self, "backTrack", nil)
		case "backTrack":
			// Check for other paths before popping stack again
			self.Send(self, "visit", back)
		}
	}
}

func Probe() Script {
	// wired
	var (
		next *Object
		cell *Object
		wall *Object
	)
	// variables
	var (
		terminator *Object
	)
	return func(e Event, self *Object) {
		switch e.Name {
		case "_wire":
			w := e.Arg.(wiring)
			next = w["next"]
			cell = w["cell"]
			wall = w["wall"]
		case "visitRand":
			n := e.Arg.(int)
			if n <= 0 {
				next.Send(self, "tryVisit", self)
			}
			if n > 0 {
				next.Send(self, "visitRand", n-1)
			}
		case "tryVisit":
			obj := e.Arg.(*Object)
			terminator = obj
			wall.Send(self, "check", self)
		case "check":
			cell.Send(self, e.Name, e.Arg)
		case "checkResult":
			visited := e.Arg.(bool)
			if !visited {
				wall.Send(self, "visit", self)
			}
			if visited {
				if self == terminator {
					cell.Send(self, "deadEnd", nil)
				}
				if self != terminator {
					next.Send(self, "tryVisit", terminator)
				}
			}
		case "visit":
			cell.Send(self, e.Name, e.Arg)
		case "backTrack":
			// TODO
		}
	}
}

func Wall() Script {
	// wired
	var (
		probe1 *Object
		probe2 *Object
	)
	return func(e Event, self *Object) {
		switch e.Name {
		case "_wire":
			w := e.Arg.(wiring)
			probe1 = w["probe1"]
			probe2 = w["probe2"]
		case "visit":
			// TODO: Open wall!
			fallthrough // just to avoid copy-pasta
		case "check":
			if e.Arg == probe1 {
				probe2.Send(self, e.Name, e.Arg)
			}
			if e.Arg == probe2 {
				probe1.Send(self, e.Name, e.Arg)
			}
		}
	}
}

// Terminator returns a dummy "visited" script so that we don't traverse it.
func Terminator() Script {
	return func(e Event, self *Object) {
		switch e.Name {
		case "check":
			obj := e.Arg.(*Object)
			obj.Send(self, "checkResult", true)
		}
	}
}
