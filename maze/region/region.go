package region

import (
	"fmt"
	"strings"

	"github.com/misterikkit/automata/maze/game"
)

const (
	wall = -1
)

type Mapped struct {
	game.Game
	tags [][]int
}

func (m Mapped) GetTag(row, col int) int { return m.tags[row][col] }

func Map(g game.Game) Mapped {
	m := Mapped{
		Game: g,
		tags: make([][]int, g.Rows()),
	}
	for i := range m.tags {
		m.tags[i] = make([]int, m.Cols())
	}
	// Pass 1: distinguish wall from empty space
	for r := range m.tags {
		for c := range m.tags[r] {
			if m.Get(r, c) {
				m.tags[r][c] = wall
			}
		}
	}
	// Pass 2: tag all cells and build equivalence map
	equiv := map[int]int{}
	nextTag := 1
	for r := range m.tags {
		for c := range m.tags[r] {
			if m.tags[r][c] == wall {
				continue
			}
			ns := neighborTags(m.tags, r, c)
			if len(ns) == 0 {
				// time for a new tag
				tag := nextTag
				nextTag++
				m.tags[r][c] = tag
				continue
			}
			// Congratulations! All these tags are equivalent
			tag := mergeTags(equiv, ns)
			m.tags[r][c] = tag
		}
	}

	// Pass 3: replace all tags with equivalence, and normalize tags at the same time
	norms := map[int]int{}
	nextNorm := 0
	for r := range m.tags {
		for c := range m.tags[r] {
			if m.tags[r][c] == wall {
				continue
			}
			tag := follow(equiv, m.tags[r][c])
			if _, ok := norms[tag]; !ok {
				norms[tag] = nextNorm
				nextNorm++
			}
			m.tags[r][c] = norms[tag]
		}
	}

	return m
}

// neighborTags returns a slice of tag values from a cell's four neighbors (no diagonals).
// This omits values <= 0 because those represent walls and untagged cells.
func neighborTags(tags [][]int, row, col int) []int {
	var ret []int
	if r := row - 1; r >= 0 {
		if t := tags[r][col]; t > 0 {
			ret = append(ret, t)
		}
	}
	if r := row + 1; r < len(tags) {
		if t := tags[r][col]; t > 0 {
			ret = append(ret, t)
		}
	}
	if c := col - 1; c >= 0 {
		if t := tags[row][c]; t > 0 {
			ret = append(ret, t)
		}
	}
	if c := col + 1; c < len(tags[0]) {
		if t := tags[row][c]; t > 0 {
			ret = append(ret, t)
		}
	}
	return ret
}

// Update the equivalence map so that all given tags are equivalent, then return
// the canonical version of this equivalence class.
func mergeTags(m map[int]int, tags []int) int {
	// Get the tails of each equivalence chain
	actual := make([]int, len(tags))
	for i := range tags {
		actual[i] = follow(m, tags[i])
	}
	// Pick the smallest one to be canonical
	tag := min(actual)
	for _, t := range actual {
		if t == tag {
			continue
		}
		m[t] = tag
	}
	return tag
}

// Follows the chain of equivalence to its terminus. For unknown values, returns the input. Warning: no cycles allowed
func follow(m map[int]int, k int) int {
	prev := k
	next, ok := m[prev]
	for ok {
		prev = next
		next, ok = m[prev]
	}
	return prev
}

// generics coming soon to a golang near you!
func min(vs []int) int {
	if len(vs) == 0 {
		panic("Don't call min with empty input") // lazy way
	}
	ret := vs[0]
	for _, v := range vs {
		if v < ret {
			ret = v
		}
	}
	return ret
}

func (m Mapped) String() string {
	var parts []string
	for r := range m.tags {
		for c := range m.tags[r] {
			parts = append(parts, fmt.Sprintf("%6d", m.tags[r][c]))
		}
		parts = append(parts, "\n")
	}
	return strings.Join(parts, "")
}
