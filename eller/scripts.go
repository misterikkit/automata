package main

import (
	"log"
	"math/rand"

	"github.com/misterikkit/automata/horizon"
)

func Controller() horizon.Script {
	return func(self horizon.Object, e horizon.Event) {
		head := self.Wires()["head"]
		switch e.Name {
		case "worldStart":
			self.Send(self, "computeEWBegin", nil)
		case "computeEWBegin":
			self.Send(head, "computeEW", nil)
		case "computeEW":
			// computeEW is done now
			self.Send(head, "computeNS", nil)
		}
	}
}

func Cell(last bool) horizon.Script {
	var (
		lastCell = last
		// open funcs take the place of object pointer + moveTo
		openEast  func()
		openSouth func()
		// doubly linked list (cycle) of group members
		groupNext horizon.Object
		groupPrev horizon.Object
		// temporary variable
		swapBuddy horizon.Object

		groupHead bool // true if we should terminate a group walk here. (e.g. group search or group count)
	)
	return func(self horizon.Object, e horizon.Event) {
		nextCell := self.Wires()["nextCell"]
		switch e.Name {
		case "worldStart":
			groupNext, groupPrev = self, self
		case "triggerEast":
			openEast = e.Arg.(func())
		case "triggerSouth":
			openSouth = e.Arg.(func())

		case "computeEW":
			if lastCell {
				// Skip last cell of each row so we don't open an outer wall.
				self.Send(nextCell, "computeEW", nil)
			}
			if !lastCell {
				// check if nextCell is in the same group. Response is handled in
				// "groupSearch" and "groupFound" events.
				groupHead = true
				self.Send(groupNext, "groupSearch", nextCell)
			}

		case "groupSearch":
			obj := e.Arg.(horizon.Object)
			if !groupHead {
				if self == obj {
					self.Send(groupNext, "groupFound", obj)
				}
				if self != obj {
					self.Send(groupNext, "groupSearch", obj)
				}
			}
			if groupHead {
				groupHead = false
				// nextCell is not in our group!
				// randomly decide to merge
				if p(0.5) {
					// invoke openEast
					openEast()
					// Swap the groupNext of self and nextCell, and the groupPrev of
					// self.groupNext and nextCell.groupNext.
					self.Send(nextCell, "getGroupNext", self)
				} else {
					// TODO: else is not supported in Horizon
					// send "computeEW" to nextCell
					self.Send(nextCell, "computeEW", nil)
				}
			}

		case "getGroupNext":
			obj := e.Arg.(horizon.Object)
			self.Send(obj, "getGroupNextResp", groupNext)

		case "getGroupNextResp":
			obj := e.Arg.(horizon.Object)
			// This is fire-and-forget, but should execute before the swapGroupPrev maneuver.
			self.Send(nextCell, "setGroupNext", groupNext)
			self.Send(groupNext, "swapGroupPrev", obj)
			groupNext = obj

		case "setGroupNext":
			obj := e.Arg.(horizon.Object)
			groupNext = obj

		case "swapGroupPrev":
			obj := e.Arg.(horizon.Object)
			// Trade groupPrev values with obj
			swapBuddy = obj
			self.Send(swapBuddy, "getGroupPrev", self)

		case "getGroupPrev":
			obj := e.Arg.(horizon.Object)
			self.Send(obj, "getGroupPrevResp", groupPrev)

		case "getGroupPrevResp":
			obj := e.Arg.(horizon.Object)
			// One of these groupPrev values is the one which initiated the
			// swapGroupPrev. Signal it that swap is complete.
			self.Send(swapBuddy, "setGroupPrev", groupPrev)
			self.Send(groupPrev, "swapComplete", nil)
			groupPrev = obj

		case "swapComplete":
			self.Send(nextCell, "computeEW", nil)

		case "groupFound":
			if !groupHead {
				self.Send(groupNext, "groupFound", e.Arg)
			}
			if groupHead {
				// neighbor is in our group. Move on
				groupHead = false
				self.Send(nextCell, "computeEW", nil)
			}

		case "computeNS":
			groupHead = true
			self.Send(groupNext, "groupCount", 1)

		case "groupCount":
			count := e.Arg.(int)
			if !groupHead {
				self.Send(groupNext, "groupCount", count+1)
			}
			if groupHead {
				log.Printf("group with %v has size %d", self, count)
			}
		}

	}
}

// p returns true with probability equal to p.
func p(p float32) bool {
	return rand.Float32() < p
}
