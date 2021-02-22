package wall_test

import (
	"testing"

	"github.com/misterikkit/automata/wall"
)

func TestWall(t *testing.T) {
	m := wall.NewMaze(10, 10)
	m.Open(0, 0, wall.South)
	// m.Open(0, 0, wall.East)
	m.Open(1, 0, wall.East)
	m.Open(0, 1, wall.South)
	m.Open(0, 1, wall.East)
	m.Open(0, 2, wall.South)
	// TODO: more thorough testing.
	t.Logf("Maze:\n%v\n", m)
}
