package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/misterikkit/automata/wall"
)

// Props to [1] for helping me understand Eller's algorithm!
// [1]: https://weblog.jamisbuck.org/2010/12/29/maze-generation-eller-s-algorithm
// Last row has a tendency to be highly connected

func main() {
	h := flag.Int("h", 5, "height")
	w := flag.Int("w", 5, "width")
	flag.Parse()
	rand.Seed(time.Now().Unix())

	maze := buildEllerMaze(*h, *w)
	fmt.Println(maze)
}

func buildEllerMaze(rows, cols int) *wall.Maze {
	maze := wall.NewMaze(rows, cols)

	state := newState(cols)
	prev := state
	log.Printf("init: %v", state)

	for r := 0; r < rows; r++ {
		lastRow := r == rows-1
		state = joinHorizontal(state, lastRow)
		log.Printf("row:  %v", state)
		maze.Set(r, 0, state.cells[0].String())
		for c := 1; c < cols; c++ {
			maze.Set(r, c, state.cells[c].String())
			// Only open walls if cells are newly joined
			if state.cells[c] == state.cells[c-1] && prev.cells[c] != prev.cells[c-1] {
				maze.Open(r, c, wall.West)
			}
		}
		if lastRow {
			break
		}
		next := joinVertical(state)
		log.Printf("col:  %v", state)
		for c := 0; c < cols; c++ {
			if state.cells[c] == next.cells[c] {
				maze.Open(r, c, wall.South)
			}
		}
		prev, state = state, next
	}
	return maze
}

type ellerState struct {
	// Set id of each cell
	cells []setID
	// map of set id to list of cell indexes
	sets map[setID]set
}

type cellPos int
type setID int

type set map[cellPos]struct{}

func newState(cols int) ellerState {
	state := ellerState{
		cells: make([]setID, cols),
		sets:  make(map[setID]set),
	}
	for i := range state.cells {
		state.cells[i] = setID(i)
		state.sets[setID(i)] = set{cellPos(i): struct{}{}}
	}
	return state
}

// Return a state where adjacent sets are randomly merged
func joinHorizontal(state ellerState, all bool) ellerState {
	next := state.Copy()
	for i := 0; i < len(next.cells)-1; i++ {
		if next.cells[i] == next.cells[i+1] {
			continue
		}
		// Buck used 50% chance of joining adjacent, nonmatching neighbors
		if all || rand.Int()%2 == 0 {
			log.Printf("Joining cols %d and %d", i, i+1)
			// Flood fill to the right with the set ID from the left
			new, old := next.cells[i], next.cells[i+1]
			for j := i + 1; j < len(next.cells) && next.cells[j] == old; j++ {
				next.cells[j] = new
				next.sets[new][cellPos(j)] = struct{}{}
				delete(next.sets[old], cellPos(j))
			}
			// Remove empty sets
			if len(next.sets[old]) == 0 {
				delete(next.sets, old)
			}
		}
	}
	return next
}

func joinVertical(state ellerState) ellerState {
	// Each row can add at most N-1 new set IDs, so any set ID + N is available for
	// the next row.
	nextID := max(state.cells) + setID(len(state.cells))
	// Start with a copy, which effectively propagates everything, then pare it back.
	next := state.Copy()
	for id, cells := range state.sets {

		// Copy the set into a slice for uniform shuffling (don't rely on golang arbitrary order)
		var dontPropagate []cellPos
		for cell := range cells {
			dontPropagate = append(dontPropagate, cell)
		}
		rand.Shuffle(len(dontPropagate), func(i, j int) {
			dontPropagate[i], dontPropagate[j] = dontPropagate[j], dontPropagate[i]
		})
		// Buck chose a uniformly random number of cells from each set to propagate
		// down, with minimum 1 and maximum all.
		// We want to keep at least 1, so we will "un-keep" at most n-1
		dropN := 0
		if len(dontPropagate) > 1 {
			dropN = rand.Intn(len(dontPropagate))
		}
		dontPropagate = dontPropagate[:dropN]
		for _, drop := range dontPropagate {
			log.Printf("Not propagating col %d", drop)
			// brand new set ID in position drop
			next.cells[drop] = nextID
			// add drop to the new set
			next.sets[nextID] = set{drop: struct{}{}}
			// remove drop from the old set
			delete(next.sets[id], drop)
			nextID++
		}
	}
	return next
}

func (e ellerState) Copy() ellerState {
	other := ellerState{
		cells: make([]setID, len(e.cells)),
		sets:  make(map[setID]set),
	}
	copy(other.cells, e.cells)
	for k, v := range e.sets {
		other.sets[k] = make(set)
		for c := range v {
			other.sets[k][c] = struct{}{}
		}
	}
	return other
}

func (i setID) String() string {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	return string(chars[int(i)%len(chars)])
}

func (state ellerState) String() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprint(state.cells))
	b.WriteString(" {")
	for id, cells := range state.sets {
		b.WriteString(fmt.Sprintf("%v:%v ", id, cells))
		// TODO: remove trailing space
	}
	b.WriteString("}")
	return b.String()
}

func (s set) String() string {
	keys := []cellPos{}
	for key := range s {
		keys = append(keys, key)
	}
	// TODO: sort
	return fmt.Sprint(keys)
}

func max(s []setID) setID {
	if len(s) == 0 {
		return -1
	}
	m := s[0]
	for _, v := range s {
		if v > m {
			m = v
		}
	}
	return m
}
