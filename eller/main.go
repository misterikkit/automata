package main

/*
# Eller's maze algorithm implemented using FB Horizon's object/event model

# High level overview
For each row,
1. Compute east-west openings
 1.1. randomly decide which ones (last row special case!)
 1.2. merge sets
 1.3. mark state as open
2. Compute north-south openings
 2.1. from each set
  2.1.1. pick one cell at random (skipping already picked ones)
  2.1.2. mark state as open
  2.1.3. decide whether to do that again.
3. Apply decisions to world
4. Advance to the next row
 4.1. any cell closed to the south is removed from current set and added to brand new set
 4.2. all walls are marked closed again

# Wiring

There are "cell" assemblies for each cell in one row of the maze. The overall row assembly will compute the state for one row, then move to the next row. Triggers will update references to the actual maze walls that need mutating.

A Cell assembly uses two triggers to capture pointers to the east and south wall of its cell.
┌──────┐  ┌───────┐
│ Cell ├──► East  │
└──┬───┘  │Trigger│
┌──▼────┐ └───────┘
│ South │
│Trigger│
└───────┘


The Row assembly is a singly-linked list that cycles back to the control node.
 ┌───────┐
 │Control◄───────────────────────────────────┐
 │ Node  │                                   │
 └───┬───┘┌──────┐     ┌──────┐     ┌──────┐ │
     └────► Cell ├─────► Cell ├─...─► Cell ├─┘
          └──────┘     └──────┘     └──────┘

TODO: Describe set membership wiring
*/

import (
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
}
