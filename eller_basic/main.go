package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/misterikkit/automata/wall"
)

// Props to [1] for helping me understand Eller's algorithm!
// [1]: https://weblog.jamisbuck.org/2010/12/29/maze-generation-eller-s-algorithm

// state represents one row of the maze, which is all that the Eller algorithm
// needs in memory at any time.
type state struct {
	// list of group IDs indexed by cell position
	groupIDs []int
	// list of cell positions indexed by group ID
	groups map[int][]int
	// Whether the east/south wall is open for the cell at that position
	openEast  []bool
	openSouth []bool
}

// new instantiates a fresh state.
func new(cols int) *state {
	s := &state{
		groupIDs:  make([]int, cols),
		groups:    make(map[int][]int),
		openEast:  make([]bool, cols),
		openSouth: make([]bool, cols),
	}
	s.nextRow()
	return s
}

// compute randomly removes walls between cells, causing groups to merge, then
// randomly selects 1 or more cell from each group to advance to the next row
// (by removing its south wall).
func (s *state) compute(lastRow bool) {
	for i := 0; i < len(s.groupIDs)-1; i++ {
		if s.groupIDs[i] == s.groupIDs[i+1] {
			continue
		}
		// Buck used 50% chance of joining adjacent, nonmatching neighbors.
		// On the last row, we connect all isolated subsections of the maze.
		if lastRow || p(0.5) {
			// log.Printf("Merge %v and %v", i, i+1)
			s.openEast[i] = true
			s.replace(s.groupIDs[i+1], s.groupIDs[i])
		}
	}
	if lastRow {
		return
	}
	for _, group := range s.groups {
		// Buck chose a uniformly random number of cells from each set to propagate
		// down, with minimum 1 and maximum all.
		propagate := 1 + rand.Intn(len(group))
		rand.Shuffle(len(group), func(i, j int) { group[i], group[j] = group[j], group[i] })
		for _, pos := range group[:propagate] {
			// log.Printf("Propagate %v", pos)
			s.openSouth[pos] = true
		}
	}
}

// replace merges two groups, replacing all references to the old group with the
// new group.
func (s *state) replace(old, new int) {
	s.groups[new] = append(s.groups[new], s.groups[old]...)
	delete(s.groups, old)
	for i := range s.groupIDs {
		if s.groupIDs[i] == old {
			s.groupIDs[i] = new
		}
	}
}

// nextRow resets the state in preparation for computing the next row, by
// - creating new group IDs for cells that don't advance to the next row
// - resetting all walls
func (s *state) nextRow() {
	for i := range s.groupIDs {
		if !s.openSouth[i] {
			s.removeOne(i)
		}
		s.openEast[i] = false
		s.openSouth[i] = false
	}
	nextID := max(s.groupIDs) + 1
	for i := range s.groupIDs {
		if s.groupIDs[i] != 0 {
			continue
		}
		s.groupIDs[i] = nextID
		s.groups[nextID] = []int{i}
		nextID++
	}
}

// removeOne removes a single cell position from a group, and sets that cell's
// group to 0 (invalid). This is common when some members of a group do not
// advance to the next row. If a group becomes empty, it is completely deleted.
func (s *state) removeOne(pos int) {
	groupID := s.groupIDs[pos]
	newGroup := []int{}
	for _, p := range s.groups[groupID] {
		if p == pos {
			continue
		}
		newGroup = append(newGroup, p)
	}
	if len(newGroup) == 0 {
		delete(s.groups, groupID)
	} else {
		s.groups[groupID] = newGroup
	}
	s.groupIDs[pos] = 0
}

func main() {
	h := flag.Int("h", 10, "height")
	w := flag.Int("w", 10, "width")
	flag.Parse()
	rand.Seed(time.Now().Unix())

	s := new(*w)
	// log.Println(s)
	// wall.Maze is a utility for pretty-printing mazes.
	maze := wall.NewMaze(*h, *w)

	for r := 0; r < *h; r++ {
		lastRow := r+1 == *h
		s.compute(lastRow)
		// log.Println(s)
		// Copy computed row into the wall.Maze
		for c := 0; c < *w; c++ {
			// maze.Set(r, c, asID(s.groupIDs[c]))
			if s.openEast[c] {
				maze.Open(r, c, wall.East)
			}
			if s.openSouth[c] {
				maze.Open(r, c, wall.South)
			}
		}
		s.nextRow()
	}

	fmt.Println(maze)
}

func max(vs []int) int {
	if len(vs) == 0 {
		return 0
	}
	val := vs[0]
	for _, v := range vs {
		if v > val {
			val = v
		}
	}
	return val
}

// convenience function to create single-character IDs for debugging the mazes
func asID(i int) string {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	return string(chars[i%len(chars)])
}

// p returns true with probability equal to p
func p(p float64) bool {
	return rand.Float64() < p
}
