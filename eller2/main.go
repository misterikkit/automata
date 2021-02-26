package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/misterikkit/automata/wall"
)

type state struct {
	// which set does cell i belong to
	groupIDs []int
	// which indexes are in group g
	groups    map[int][]int
	openEast  []bool
	openSouth []bool
}

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

func (s *state) compute(lastRow bool) {
	for i := 0; i < len(s.groupIDs)-1; i++ {
		if s.groupIDs[i] == s.groupIDs[i+1] {
			continue
		}
		if lastRow || rand.Intn(2) == 0 {
			// log.Printf("Merge %v and %v", i, i+1)
			s.openEast[i] = true
			s.replace(s.groupIDs[i+1], s.groupIDs[i])
		}
	}
	if lastRow {
		return
	}
	for _, group := range s.groups {
		propagate := 1 + rand.Intn(len(group))
		rand.Shuffle(len(group), func(i, j int) { group[i], group[j] = group[j], group[i] })
		for _, pos := range group[:propagate] {
			// log.Printf("Propagate %v", pos)
			s.openSouth[pos] = true
		}
	}
}

func (s *state) replace(old, new int) {
	s.groups[new] = append(s.groups[new], s.groups[old]...)
	delete(s.groups, old)
	for i := range s.groupIDs {
		if s.groupIDs[i] == old {
			s.groupIDs[i] = new
		}
	}
}

func (s *state) nextRow() {
	for i := range s.groupIDs {
		if !s.openSouth[i] {
			s.removeOne(i, s.groupIDs[i])
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

func (s *state) removeOne(pos, groupID int) {
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

func main() {
	h := flag.Int("h", 5, "height")
	w := flag.Int("w", 5, "width")
	flag.Parse()
	rand.Seed(time.Now().Unix())

	s := new(*w)
	// log.Println(s)
	maze := wall.NewMaze(*h, *w)

	for r := 0; r < *h; r++ {
		lastRow := r+1 == *h
		s.compute(lastRow)
		// log.Println(s)
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

func asID(i int) string {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	return string(chars[i%len(chars)])
}
