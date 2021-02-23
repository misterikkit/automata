package wall

import "bytes"

// Direction is a compass direction
type Direction int

// The four directions. These can be masked together
const (
	North Direction = 1 << iota
	East
	South
	West
)

// Maze is a 2D, walled maze. Each cell has four walls around it which can be
// opened to create a maze.
type Maze struct {
	cells [][]cell
}

func NewMaze(rows, cols int) *Maze {
	m := &Maze{
		cells: make([][]cell, rows),
	}
	for r := range m.cells {
		m.cells[r] = make([]cell, cols)
	}
	return m
}

func (m *Maze) Open(row, col int, d Direction) {
	nextRow, nextCol := row, col
	var nextDirection Direction
	switch d {
	case North:
		nextRow--
		nextDirection = South
	case East:
		nextCol++
		nextDirection = West
	case South:
		nextRow++
		nextDirection = North
	case West:
		nextCol--
		nextDirection = East
	default:
		panic("one direction at a time, please")
	}
	if !m.valid(row, col) || !m.valid(nextRow, nextCol) {
		return
		// TODO: error here?
	}
	m.cells[row][col].openings |= d
	m.cells[nextRow][nextCol].openings |= nextDirection
}

// Set sets a one-rune value to print in the cell.
func (m *Maze) Set(row, col int, val string) {
	if !m.valid(row, col) {
		return
	}
	m.cells[row][col].value = val
}

func (m *Maze) valid(row, col int) bool {
	return row >= 0 && row < len(m.cells) && col >= 0 && col < len(m.cells[row])
}

type cell struct {
	// A bitmask of which walls are open
	openings Direction
	// optional display value
	value string
}

func (m *Maze) String() string {
	var b bytes.Buffer
	// north border
	d := South | East
	for c := range m.cells[0] {
		b.WriteString(m.cornerSegmentNW(0, c) + segment(East|West))
		d |= West
	}
	b.WriteString(m.cornerSegmentNW(0, len(m.cells[0])) + "\n")
	for r, row := range m.cells {
		// cell row
		b.WriteString(segment(North | South))
		for _, cell := range row {
			if len(cell.value) == 1 {
				b.WriteString(cell.value)
			} else {
				b.WriteString(" ")
			}
			if cell.openings&East > 0 {
				b.WriteString(" ")
			} else {
				b.WriteString("" + segment(North|South))
			}
		}
		b.WriteString("\n")
		// wall row
		b.WriteString(m.cornerSegmentNW(r+1, 0))
		for c, cell := range row {
			if cell.openings&South > 0 {
				b.WriteString(" " + m.cornerSegmentNW(r+1, c+1))
			} else {
				b.WriteString(segment(East|West) + m.cornerSegmentNW(r+1, c+1))
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (m *Maze) rows() int { return len(m.cells) }
func (m *Maze) cols() int { return len(m.cells[0]) }

func (m *Maze) cornerSegmentNW(row, col int) string {
	mask := North | East | South | West
	// assume 0 <= row and col <= len
	if row == 0 {
		mask &= ^North
	}
	if col == 0 {
		mask &= ^West
	}
	if row >= m.rows() {
		mask &= ^South
	}
	if col >= m.cols() {
		mask &= ^East
	}

	if row < m.rows() && col < m.cols() {
		open := m.cells[row][col].openings
		if open&North > 0 {
			mask &= ^East
		}
		if open&West > 0 {
			mask &= ^South
		}
	}
	if row-1 >= 0 && col < m.cols() && m.cells[row-1][col].openings&West > 0 {
		mask &= ^North
	}
	if col-1 >= 0 && row < m.rows() && m.cells[row][col-1].openings&North > 0 {
		mask &= ^West
	}
	return segment(mask)
}

func segment(d Direction) string {
	// This could be a compile-time lookup table,
	// but arranging that is somewhat tedious.
	switch d {
	case West:
		return "╴"
	case South:
		return "╵"
	case East:
		return "╶"
	case North:
		return "╷"
	case East | West:
		return "─"
	case North | South:
		return "│"
	case East | South:
		return "┌"
	case West | South:
		return "┐"
	case North | East:
		return "└"
	case North | West:
		return "┘"
	case North | East | South:
		return "├"
	case North | South | West:
		return "┤"
	case East | South | West:
		return "┬"
	case North | East | West:
		return "┴"
	case North | East | South | West:
		return "┼"
	}
	return "•"
}
