package main

import "math/rand"

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
			// Select a random probe to initiate the next visit.
			probe.Send(self, "visitRand", rand.Int()%4)
		case "check":
			// Someone wants to know if this cell has been visited.
			obj := e.Arg.(*Object)
			obj.Send(self, "checkResult", visited)
		case "deadEnd":
			// pop stack. the `back` cell will check for other unvisited neighbors before itself popping.
			back.Send(self, "backTrack", nil)
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
			// Walk along our linked cycle the chosen number of steps
			n := e.Arg.(int)
			if n <= 0 {
				next.Send(self, "tryVisit", self)
			}
			if n > 0 {
				next.Send(self, "visitRand", n-1)
			}
		case "tryVisit":
			// Check if this probe's neighbor can be visited.
			obj := e.Arg.(*Object)
			terminator = obj // for linked list termination
			wall.Send(self, "check", self)
		case "checkResult":
			// Response to our check from tryVisit
			visited := e.Arg.(bool)
			if !visited {
				// Begin a visit in the next cell. Our work here is done.
				wall.Send(self, "visit", self)
			}
			if visited {
				// Can't visit this probe's neighbor, so try the next probe or signal deadEnd.
				if self == terminator {
					cell.Send(self, "deadEnd", nil)
				}
				if self != terminator {
					next.Send(self, "tryVisit", terminator)
				}
			}

		case "check":
			// These messages are simple pass-through back to the cell.
			fallthrough
		case "visit":
			fallthrough
		case "backTrack":
			cell.Send(self, e.Name, e.Arg)
		}
	}
}

// Wall represents a barrier between cells. It passes messages between two cells
// via their Probes, and is also the primary output of the algorithm!
func Wall(onOpen func()) Script {
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
			onOpen()
			fallthrough // just to avoid copy-pasta
		case "check":
			// Pass message through to the probe that did not send it.
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
