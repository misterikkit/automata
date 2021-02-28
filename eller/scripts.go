package main

import (
	"math/rand"

	"github.com/misterikkit/automata/horizon"
)

func Controller() horizon.Script {
	return func(self horizon.Object, e horizon.Event) {
		head := self.Wires()["head"]
		switch e.Name {
		case "computeEWBegin":
			self.Send(head, "computeEW", nil)
		case "computeEW":
			// computeEW is done now
		}
	}
}
func Cell() horizon.Script {
	var (
		// open funcs take the place of object pointer + moveTo
		openEast  func()
		openSouth func()
		// doubly linked list (cycle) of group members
		groupNext horizon.Object
		groupPrev horizon.Object
	)
	return func(self horizon.Object, e horizon.Event) {
		// nextCell := self.Wires()["nextCell"]
		switch e.Name {
		case "triggerEast":
			openEast = e.Arg.(func())
		case "triggerSouth":
			openSouth = e.Arg.(func())

		case "computeEW":
			// check if nextCell is in the same group
			// randomly decide to merge
			// invoke openEast
			// update groupNext & groupPrev of self and nextCell
			// send "computeEW" to nextCell
		}
	}
}

// p returns true with probability equal to p.
func p(p float32) bool {
	return rand.Float32() < p
}
