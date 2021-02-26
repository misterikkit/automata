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
	state.nextRow()
	log.Printf("init: %v", state)

	for r := 0; r < rows; r++ {
		lastRow := r == rows-1
		state.goRight(lastRow)
		for c := 0; c < cols; c++ {
			maze.Set(r, c, fmt.Sprint(state.cells[c]))
			if state.eastOpen[c] {
				maze.Open(r, c, wall.East)
			}
		}
		log.Printf("right: %v", state)
		if lastRow {
			break
		}
		state.goDown()
		for c := 0; c < cols; c++ {
			if state.southOpen[c] {
				maze.Open(r, c, wall.South)
			}
		}
		log.Printf("right: %v", state)
		state.nextRow()
	}

	return maze
}

type ellerState struct {
	// Set id of each cell
	cells               []setID
	eastOpen, southOpen []bool
	// map of set id to list of cell indexes
	sets map[setID]set
}

type cellPos int
type setID int

type set map[cellPos]struct{}

// Return n keys from the set, uniformly random
func (s set) rand(n int) []cellPos {
	ret := make([]cellPos, 0, len(s))
	for k := range s {
		ret = append(ret, k)
	}
	rand.Shuffle(len(ret), func(i, j int) { ret[i], ret[j] = ret[j], ret[i] })
	return ret[:n]
}

func newState(cols int) ellerState {
	state := ellerState{
		cells:     make([]setID, cols),
		eastOpen:  make([]bool, cols),
		southOpen: make([]bool, cols),
		sets:      make(map[setID]set),
	}
	for i := range state.cells {
		state.cells[i] = -1
		// state.sets[setID(i)] = set{cellPos(i): struct{}{}}
	}
	return state
}

func (s *ellerState) goRight(all bool) {
	for i := 0; i < len(s.cells)-1; i++ {
		if s.cells[i] == s.cells[i+1] {
			continue
		}
		if all || rand.Intn(2) == 0 {
			log.Printf("Joining %v and %v", i, i+1)
			new, old := s.cells[i], s.cells[i+1]
			// s.set(cellPos(i+1), s.cells[i])
			for pos := range s.sets[old] {
				s.set(pos, new)
			}
			s.eastOpen[i] = true
		}
	}
}

func (s *ellerState) goDown() {
	for _, set := range s.sets {
		propagate := 1
		if len(set) > 1 {
			propagate = 1 + rand.Intn(len(set)-1)
		}
		for _, pos := range set.rand(propagate) {
			log.Printf("Propagating %v", pos)
			s.southOpen[pos] = true
		}
	}
	for i := range s.cells {
		if !s.southOpen[i] {
			s.cells[i] = -1
		}
	}
}

func (s *ellerState) nextRow() {
	nextID := max(s.cells) + 1
	for i := range s.cells {
		if s.cells[i] != -1 {
			continue
		}
		s.set(cellPos(i), nextID)
		nextID++
	}
	for i := range s.eastOpen {
		s.eastOpen[i] = false
		s.southOpen[i] = false
	}
}

// set a cell to be in a set, and update internal state accordingly
func (s *ellerState) set(p cellPos, i setID) {
	old := s.cells[p]
	s.cells[p] = i
	if s.sets[i] == nil {
		s.sets[i] = make(set)
	}
	s.sets[i][p] = struct{}{}
	if oldSet, ok := s.sets[old]; ok {
		delete(oldSet, p) // oldSet can't be empty because we will delete it if so
		if len(oldSet) == 0 {
			delete(s.sets, old)
		}
	}
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
			next.eastOpen[i] = true
			next.set(cellPos(i+1), next.cells[i])
			// new, old := next.cells[i], next.cells[i+1]
			// next.cells[i+1] = new
			// next.sets[new][cellPos(i+1)] = struct{}{}
			// delete(next.sets[old], cellPos(i+1))
			// if len(next.sets[old]) == 0 {
			// 	delete(next.sets, old)
			// }
		}
	}
	return next
}

func joinVertical(state ellerState) ellerState {
	nextID := max(state.cells) + 1
	// Start with a copy, which effectively propagates everything, then pare it back.
	next := newState(len(state.cells))
	// Pick at least one cell from each set to propagate down.
	for id, cells := range state.sets {

		// Copy the set into a slice for uniform shuffling (don't rely on golang arbitrary order)
		var cellList []cellPos
		for cell := range cells {
			cellList = append(cellList, cell)
		}
		rand.Shuffle(len(cellList), func(i, j int) {
			cellList[i], cellList[j] = cellList[j], cellList[i]
		})
		// Buck chose a uniformly random number of cells from each set to propagate
		// down, with minimum 1 and maximum all.
		keep := 1
		if len(cellList) > 1 {
			keep = rand.Intn(len(cellList)-1) + 1
		}
		cellList = cellList[:keep]
		for _, froob := range cellList {
			log.Printf("Propagating col %d", froob)
			// brand new set ID in position drop
			next.cells[froob] = state.cells[froob]
			next.sets[state.cells[froob]][froob] = struct{}{}
			// add drop to the new set
			next.sets[nextID] = set{froob: struct{}{}}
			// remove drop from the old set
			delete(next.sets[id], froob)
			nextID++
		}
	}
	return next
}

func (e ellerState) Copy() ellerState {
	other := newState(len(e.cells))
	copy(other.cells, e.cells)
	copy(other.eastOpen, e.eastOpen)
	copy(other.southOpen, e.southOpen)
	// other.sets = make(map[setID]set)
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
