package main

import (
	"context"
	"log"
	"math/rand"
	"time"
)

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

	wall1 := New("wall1", Wall())

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

	log.Printf("%-14v %-14q -> %-14v (%v)", "sender", "event", "recipient", "param")

	// make sure all _wire events are done
	time.AfterFunc(500*time.Millisecond, func() { cell1.Send(init, "visit", init) })
	// start the engine
	RunAll(ctx,
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
