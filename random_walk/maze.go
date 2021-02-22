package main

import (
	"bytes"
	"context"
	"fmt"

	"github.com/misterikkit/automata/horizon"
)

type CellGroup struct {
	cell                     horizon.Object
	north, south, east, west horizon.Object // probes
}

// CellPartial contains a cell, it's probes, and two walls (north & west) such
// that CellPartials can be tiled to create a maze. Southmost and Eastmost walls
// must be treated separately.
type CellPartial struct {
	cell                           horizon.Object
	probeN, probeE, probeS, probeW horizon.Object
	wallN, wallW                   horizon.Object
	openN, openW                   bool // whether those walls are open
	// TODO: corner?
}

type Maze struct {
	el     horizon.EventLoop
	cells  [][]CellPartial
	border horizon.Object
}

func NewMaze(rows, cols int) *Maze {
	m := &Maze{el: horizon.NewEventLoop()}
	m.cells = make([][]CellPartial, rows)
	for i := range m.cells {
		m.cells[i] = make([]CellPartial, cols)
	}
	for r := range m.cells {
		for c := range m.cells[r] {
			name := fmt.Sprintf("cell[%d,%d]", r, c)
			m.cells[r][c] = CellPartial{
				cell:   horizon.NewObject(name, Cell(), m.el),
				probeN: horizon.NewObject(fmt.Sprintf("%s-probe-N", name), Probe(), m.el),
				probeE: horizon.NewObject(fmt.Sprintf("%s-probe-E", name), Probe(), m.el),
				probeS: horizon.NewObject(fmt.Sprintf("%s-probe-S", name), Probe(), m.el),
				probeW: horizon.NewObject(fmt.Sprintf("%s-probe-W", name), Probe(), m.el),
			}
			// Capture bool address in local var for the closure
			openN := &m.cells[r][c].openN
			openW := &m.cells[r][c].openW
			m.cells[r][c].wallN = horizon.NewObject(fmt.Sprintf("%s-wall-N", name), Wall(func() { *openN = true }), m.el)
			m.cells[r][c].wallW = horizon.NewObject(fmt.Sprintf("%s-wall-W", name), Wall(func() { *openW = true }), m.el)
		}
	}
	// Time to wire them up!
	m.border = horizon.NewObject("border", Terminator(), m.el)

	for r := range m.cells {
		for c := range m.cells[r] {
			partial := m.cells[r][c]
			partial.cell.Wire(horizon.Wiring{"probe": partial.probeN})
			// Determine whether to use real wall, or border sentinel.
			wallN, wallE, wallS, wallW := m.border, m.border, m.border, m.border
			if r > 0 {
				partial.wallN.Wire(horizon.Wiring{"probe1": partial.probeN, "probe2": m.cells[r-1][c].probeS})
				wallN = partial.wallN
			}
			if c > 0 {
				partial.wallW.Wire(horizon.Wiring{"probe1": partial.probeW, "probe2": m.cells[r][c-1].probeE})
				wallW = partial.wallW
			}
			if r+1 < len(m.cells) {
				wallS = m.cells[r+1][c].wallN
			}
			if c+1 < len(m.cells[r]) {
				wallE = m.cells[r][c+1].wallW
			}
			partial.probeN.Wire(horizon.Wiring{"cell": partial.cell, "next": partial.probeE, "wall": wallN})
			partial.probeE.Wire(horizon.Wiring{"cell": partial.cell, "next": partial.probeS, "wall": wallE})
			partial.probeS.Wire(horizon.Wiring{"cell": partial.cell, "next": partial.probeW, "wall": wallS})
			partial.probeW.Wire(horizon.Wiring{"cell": partial.cell, "next": partial.probeN, "wall": wallW})
		}
	}
	return m
}

// Run runs the maze generation algorithm, returning upon completion.
func (m *Maze) Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	init := horizon.NewObject("init", func(horizon.Object, horizon.Event) { cancel() }, m.el)
	init.Send(m.cells[0][0].cell, "visit", init)
	m.el.Run(ctx)
}

func (m *Maze) String() string {
	var b bytes.Buffer
	// Each partial has two rows of text: north walls and east-west walls
	for r, row := range m.cells {
		for c, partial := range row {
			b.WriteString(m.cornerNW(r, c))
			if r == 0 || !partial.openN {
				b.WriteString("─")
			} else {
				b.WriteString(" ")
			}
		}
		// Eastern corner at end of row
		b.WriteString(m.cornerNW(r, len(row)))
		b.WriteString("\n")
		for c, partial := range row {
			if c == 0 || !partial.openW {
				b.WriteString("│ ")
			} else {
				b.WriteString("  ")
			}
		}
		b.WriteString("│\n")

	}
	for c := range m.cells[0] {
		b.WriteString(m.cornerNW(len(m.cells), c))
		b.WriteString("─")
	}
	b.WriteString(m.cornerNW(len(m.cells), len(m.cells[0])))
	b.WriteString("\n")
	return b.String()
}

// Returns an appropriate symbol for the northwest corner of the given partial
func (m *Maze) cornerNW(row, col int) string {
	const n, e, s, w = 1, 2, 4, 8
	var mask int
	rows, cols := len(m.cells), len(m.cells[0])
	valid := func(r, c int) bool { return r >= 0 && r < rows && c >= 0 && c < cols }
	if valid(row, col) && !m.cells[row][col].openN {
		mask |= e
	}
	if valid(row, col) && !m.cells[row][col].openW {
		mask |= s
	}
	if r := row - 1; valid(r, col) && !m.cells[r][col].openW {
		mask |= n
	}
	if c := col - 1; valid(row, c) && !m.cells[row][c].openN {
		mask |= w
	}
	// Special case "always closed" border walls on south & east
	if col >= cols {
		if row > 0 {
			mask |= n
		}
		if row < rows {
			mask |= s
		}
	}
	if row >= rows {
		if col > 0 {
			mask |= w
		}
		if col < cols {
			mask |= e
		}
	}
	switch mask {
	case w:
		return "╴"
	case s:
		return "╵"
	case e:
		return "╶"
	case n:
		return "╷"
	case e | w:
		return "─"
	case n | s:
		return "│"
	case e | s:
		return "┌"
	case w | s:
		return "┐"
	case n | e:
		return "└"
	case n | w:
		return "┘"
	case n | e | s:
		return "├"
	case n | s | w:
		return "┤"
	case e | s | w:
		return "┬"
	case n | e | w:
		return "┴"
	case n | e | s | w:
		return "┼"
	}
	return "•"
}

// ─
// │
// ┌
// ┐
// └
// ┘
// ├
// ┤
// ┬
// ┴
// ┼

// +-+-+-+-+
// |.|.|.|.|
// +-+-+-+-+
// |.|.|.|.|
// +-+-+-+-+
