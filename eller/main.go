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
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/misterikkit/automata/horizon"
	"github.com/misterikkit/automata/wall"
)

func main() {
	rand.Seed(time.Now().Unix())
	maze := wall.NewMaze(10, 10)
	loop := horizon.NewEventLoop()

	cells := make([]horizon.Object, 10)
	for i := range cells {
		last := i == len(cells)-1
		cells[i] = horizon.NewObject(fmt.Sprintf("cell-%02d", i), Cell(last), loop)
	}
	// Workaround to simulate the moving and trigger detecting of wall objects
	row := 0
	ctrl := horizon.NewObject("controller", Controller(func() {
		row++
		updateTriggers(cells, maze, row)
	}), loop)
	updateTriggers(cells, maze, row)

	ctrl.Wire(horizon.Wiring{"head": cells[0]})
	for i := range cells {
		j := i + 1
		if j < len(cells) {
			cells[i].Wire(horizon.Wiring{"nextCell": cells[j]})
		} else {
			cells[i].Wire(horizon.Wiring{"nextCell": ctrl})
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	loop.Run(ctx)
	fmt.Println(maze)
}

func updateTriggers(cells []horizon.Object, maze *wall.Maze, row int) {
	for i := range cells {
		col := i
		cells[i].Send(cells[i], "triggerEast", func() {
			maze.Open(row, col, wall.East)
		})
		cells[i].Send(cells[i], "triggerSouth", func() {
			maze.Open(row, col, wall.South)
		})
	}
}
