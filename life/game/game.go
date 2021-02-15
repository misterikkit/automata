package game

import (
	"math/rand"
)

// Cell is a binary state in a game grid.
type Cell bool

// Game is a grid of cells
type Game [][]Cell

// New creates a game with the given size.
func New(rows, cols int) Game {
	w := make(Game, rows)
	for jj := range w {
		w[jj] = make([]Cell, cols)
	}
	return w
}

// Next returns the next generation, applying the given rule.
func (g Game) Next(rule Rule) Game {
	next := New(g.Rows(), g.Cols())
	for r := range g {
		for c := range g[r] {
			next[r][c] = rule(g, r, c)
		}
	}
	return next
}

// Get returns the value in the given cell.
func (g Game) Get(row, col int) Cell { return g[row][col] }
func (g Game) Rows() int             { return len(g) }
func (g Game) Cols() int             { return len(g[0]) }

// Rule defines how to advance a game to the next generation.
type Rule func(g Game, row, col int) Cell

type printer interface {
	Printf(string, ...interface{}) (int, error)
	Next() error
}

// Life implements Conway's game of life.
func Life(w Game, row, col int) Cell {
	n := countNeighbors(w, row, col)
	switch w[row][col] {
	case true:
		if n == 2 || n == 3 {
			return true
		}
	case false:
		if n == 3 {
			return true
		}
	}
	return false
}

// Random sets cells randomly.
func Random(_ Game, _, _ int) Cell {
	return rand.Int()%2 == 0
}

func countNeighbors(w Game, row, col int) int {
	total := 0
	for r := -1; r <= 1; r++ {
		if !(row+r > 0 && row+r < w.Rows()) {
			continue
		}

		for c := -1; c <= 1; c++ {
			if r == 0 && c == 0 {
				continue
			}
			if !(col+c > 0 && col+c < w.Cols()) {
				continue
			}

			if w[row+r][col+c] {
				total++
			}
		}

	}
	return total
}
