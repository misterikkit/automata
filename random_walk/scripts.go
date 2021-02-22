package main

import (
	"math/rand"

	"github.com/misterikkit/automata/horizon"
)

// These scripts implement a random walk maze generation algorithm. They assume
// the following wiring diagram.
//
//
//                +-------+                                                +-------+
//     +--------->+ PROBE +----------+                          +--------->+ PROBE +----------+
//     |          +--+--+-+          |                          |          +--+--+-+          |
//     |             ^  |            |                          |             ^  |            |
//     |             |  |            |          +----+          |             |  |            |
//     |             |  |            |          |    |          |             |  |            |
//     |             |  v            |          |    |          |             |  v            |
// +---+---+      +--+--+-+      +---v---+      |    |      +---+---+      +--+--+-+      +---v---+
// | PROBE +----->+ CELL  +<-----+ PROBE +----->+WALL+<-----+ PROBE +----->+ CELL  +<-----+ PROBE |
// +---^---+      +---+---+      +---+-^-+      |    |      +-^-^---+      +---+---+      +---+---+
//     |              ^              | |        |    |        | |              ^              |
//     |              |              | +--------+    +--------+ |              |              |
//     |              |              |          +----+          |              |              |
//     |              |              |                          |              |              |
//     |          +---+---+          |                          |          +---+---+          |
//     +----------+ PROBE +<---------+                          +----------+ PROBE +<---------+
//                +--+-^--+                                                +--+-^--+
//                   | |                                                      | |
//                   | |                                                      | |
//                   | |                                                      | |
//                   v |                                                      v |
//              +----+-+----+                                            +----+-+----+
//              |   WALL    |                                            |   WALL    |
//              +----+-+----+                                            +----+-+----+
//                   ^ |                                                      ^ |
//                   | |                                                      | |
//                   | |                                                      | |
//                   | |                                                      | |
//                +--+-v--+                                                +--+-v--+
//     +--------->+ PROBE +----------+                          +--------->+ PROBE +----------+
//     |          +--+--+-+          |                          |          +--+--+-+          |
//     |             ^  |            |                          |             ^  |            |
//     |             |  |            |          +----+          |             |  |            |
//     |             |  |            |          |    |          |             |  |            |
//     |             |  v            |          |    |          |             |  v            |
// +---+---+      +--+--+-+      +---v---+      |    |      +---+---+      +--+--+-+      +---v---+
// | PROBE +----->+ CELL  +<-----+ PROBE +----->+WALL+<-----+ PROBE +----->+ CELL  +<-----+ PROBE |
// +---^---+      +---+---+      +---+-^-+      |    |      +-^-^---+      +---+---+      +---+---+
//     |              ^              | |        |    |        | |              ^              |
//     |              |              | +--------+    +--------+ |              |              |
//     |              |              |          +----+          |              |              |
//     |              |              |                          |              |              |
//     |          +---+---+          |                          |          +---+---+          |
//     +----------+ PROBE +<---------+                          +----------+ PROBE +<---------+
//                +-------+                                                +-------+

// Cell implements the behavior for one empty space in the maze. It expects a
// Probe for each neighboring cell, arranged in a linked list cycle. Only one
// Probe needs to be wired into the Cell.
func Cell() horizon.Script {
	// variables
	var (
		visited = false
		back    horizon.Object
	)
	return func(self horizon.Object, e horizon.Event) {
		probe := self.Wires()["probe"]
		switch e.Name {
		case "visit":
			obj := e.Arg.(horizon.Object)
			back = obj
			visited = true
			// Select a random probe to initiate the next visit.
			self.Send(probe, "visitRand", rand.Int()%4)
		case "check":
			// Someone wants to know if this cell has been visited.
			obj := e.Arg.(horizon.Object)
			self.Send(obj, "checkResult", visited)
		case "deadEnd":
			// pop stack. the `back` cell will check for other unvisited neighbors before itself popping.
			self.Send(back, "backTrack", nil)
		case "backTrack":
			// Check for other paths before popping stack again. Simplest way is to
			// re-use the existing "visit" logic.
			self.Send(self, "visit", back)
		}
	}
}

// Probe does most of the work in this algorithm, and represents one neighbor of
// one cell. There is another Probe in the neighboring cell which represents the
// cell of this Probe. These two probes pass messages to each other through the
// Wall between them.
func Probe() horizon.Script {
	// variables
	var (
		terminator horizon.Object
	)
	return func(self horizon.Object, e horizon.Event) {
		next := self.Wires()["next"]
		cell := self.Wires()["cell"]
		wall := self.Wires()["wall"]
		switch e.Name {
		case "visitRand":
			// Walk along our linked cycle the chosen number of steps
			n := e.Arg.(int)
			if n <= 0 {
				self.Send(next, "tryVisit", self)
			}
			if n > 0 {
				self.Send(next, "visitRand", n-1)
			}
		case "tryVisit":
			// Check if this probe's neighbor can be visited.
			obj := e.Arg.(horizon.Object)
			terminator = obj // for linked list termination
			self.Send(wall, "check", self)
		case "checkResult":
			// Response to our check from tryVisit
			visited := e.Arg.(bool)
			if !visited {
				// Begin a visit in the next cell. Our work here is done.
				self.Send(wall, "visit", self)
			}
			if visited {
				// Can't visit this probe's neighbor, so try the next probe or signal deadEnd.
				if self == terminator {
					self.Send(cell, "deadEnd", nil)
				}
				if self != terminator {
					self.Send(next, "tryVisit", terminator)
				}
			}

		case "check":
			// These messages are simple pass-through back to the cell.
			fallthrough
		case "visit":
			fallthrough
		case "backTrack":
			self.Send(cell, e.Name, e.Arg)
		}
	}
}

// Wall represents a barrier between cells. It passes messages between two cells
// via their Probes, and is also the primary output of the algorithm!
func Wall(onOpen func()) horizon.Script {
	return func(self horizon.Object, e horizon.Event) {
		probe1 := self.Wires()["probe1"]
		probe2 := self.Wires()["probe2"]
		switch e.Name {
		case "visit":
			onOpen()
			fallthrough // just to avoid copy-pasta
		case "check":
			// Pass message through to the probe that did not send it.
			if e.Arg == probe1 {
				self.Send(probe2, e.Name, e.Arg)
			}
			if e.Arg == probe2 {
				self.Send(probe1, e.Name, e.Arg)
			}
		}
	}
}

// Terminator returns a dummy "visited" script so that we don't traverse it.
func Terminator() horizon.Script {
	return func(self horizon.Object, e horizon.Event) {
		switch e.Name {
		case "check":
			obj := e.Arg.(horizon.Object)
			self.Send(obj, "checkResult", true)
		}
	}
}
