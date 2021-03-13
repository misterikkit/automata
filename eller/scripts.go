package main

import (
	"math/rand"

	"github.com/misterikkit/automata/horizon"
)

func Controller(rows int, moveToNext func(), done func()) horizon.Script {
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
		case "computeNS":
			// computeNS is done now
			moveToNext()
			rows--
			if rows > 0 {
				self.Send(self, "computeEWBegin", nil)
			}
			if rows <= 0 {
				done()
			}
		}
	}
}

func Cell(last bool) horizon.Script {
	var (
		lastCell = last
		// open funcs each take the place of a trigger + wall object
		openEast  func()
		openSouth func()
		// doubly linked list (cycle) of group members
		groupNext horizon.Object
		groupPrev horizon.Object
		// temporary variable
		swapBuddy horizon.Object

		groupHead bool // true if we should terminate a group walk here. (e.g. group search or group count)
		rowDone   bool // tracks when both east and south decisions have been made for this cell
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
			rowDone = false
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

		case "setGroupPrev":
			obj := e.Arg.(horizon.Object)
			groupPrev = obj

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
			if rowDone {
				self.Send(nextCell, "computeNS", nil)
			}
			if !rowDone {
				groupHead = true
				self.Send(groupNext, "groupCount", 1)
			}

		case "groupCount":
			count := e.Arg.(int)
			if !groupHead {
				self.Send(groupNext, "groupCount", count+1)
			}
			if groupHead {
				numToOpen := 1 + rand.Intn(count) // TODO: inline this ):
				self.Send(groupNext, "openSouthMaybe", vector{x: float32(numToOpen), y: float32(count)})
			}

		case "openSouthMaybe":
			v := e.Arg.(vector)
			// In this group, open v.x of the remaining v.y cells. In otherwords, open
			// this cell with probability v.x/v.y
			rowDone = true
			if p(v.x / v.y) {
				openSouth()
				if !groupHead {
					self.Send(groupNext, "openSouthMaybe", vector{x: v.x - 1, y: v.y - 1})
				}
			} else {
				// TODO: no else ):

				// Remove self from the group
				self.Send(groupNext, "setGroupPrev", groupPrev)
				self.Send(groupPrev, "setGroupNext", groupNext)
				if !groupHead {
					// ordering here is weird since we need a ref to groupNext to send this
					// message, but we want to update groupNext.groupPrev before it potentially
					// updates its groupPrev.groupNext.
					self.Send(groupNext, "openSouthMaybe", vector{x: v.x, y: v.y - 1})
				}
				groupNext, groupPrev = self, self
			}
			if groupHead {
				groupHead = false
				self.Send(nextCell, "computeNS", nil)
			}
		}

	}
}

type vector struct{ x, y, z float32 }

// p returns true with probability equal to p.
func p(p float32) bool {
	return rand.Float32() < p
}
