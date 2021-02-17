package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
)

type CellGroup struct {
	cell                     *Object
	north, south, east, west *Object // probes
}

// CellPartial contains a cell, it's probes, and two walls (north & west) such
// that CellPartials can be tiled to create a maze. Southmost and Eastmost walls
// must be treated separately.
type CellPartial struct {
	cell                           *Object
	probeN, probeE, probeS, probeW *Object
	wallN, wallW                   *Object
	openN, openW                   bool // whether those walls are open
	// TODO: corner?
}

type Maze struct {
	cells [][]CellPartial
}

func NewMaze(rows, cols int) *Maze {
	m := &Maze{}
	m.cells = make([][]CellPartial, rows)
	for i := range m.cells {
		m.cells[i] = make([]CellPartial, cols)
	}
	for r := range m.cells {
		for c := range m.cells[r] {
			name := fmt.Sprintf("cell[%d,%d]", r, c)
			m.cells[r][c] = CellPartial{
				cell:   New(name, Cell()),
				probeN: New(fmt.Sprintf("%s-probe-N", name), Probe()),
				probeE: New(fmt.Sprintf("%s-probe-E", name), Probe()),
				probeS: New(fmt.Sprintf("%s-probe-S", name), Probe()),
				probeW: New(fmt.Sprintf("%s-probe-W", name), Probe()),
				// Did I capture r and c correctly in these closures?
				wallN: New(fmt.Sprintf("%s-wall-N", name), Wall(func() { m.cells[r][c].openN = true; log.Printf("Open wall N of (%d, %d)", r, c) })),
				wallW: New(fmt.Sprintf("%s-wall-W", name), Wall(func() { m.cells[r][c].openW = true; log.Printf("Open wall W of (%d, %d)", r, c) })),
			}
		}
	}
	// Time to wire them up!
	border := New("border", Terminator())

	for r := range m.cells {
		for c := range m.cells[r] {
			partial := m.cells[r][c]
			partial.cell.Wire(wiring{"probe": partial.probeN})
			// Determine whether to use real wall, or border sentinel.
			wallN, wallE, wallS, wallW := border, border, border, border
			if r > 0 {
				partial.wallN.Wire(wiring{"probe1": partial.probeN, "probe2": m.cells[r-1][c].probeS})
				wallN = partial.wallN
			}
			if c > 0 {
				partial.wallW.Wire(wiring{"probe1": partial.probeW, "probe2": m.cells[r][c-1].probeE})
				wallW = partial.wallW
			}
			if r+1 < len(m.cells) {
				wallS = m.cells[r+1][c].wallN
			}
			if c+1 < len(m.cells[r]) {
				wallE = m.cells[r][c+1].wallW
			}
			partial.probeN.Wire(wiring{"cell": partial.cell, "next": partial.probeE, "wall": wallN})
			partial.probeE.Wire(wiring{"cell": partial.cell, "next": partial.probeS, "wall": wallE})
			partial.probeS.Wire(wiring{"cell": partial.cell, "next": partial.probeW, "wall": wallS})
			partial.probeW.Wire(wiring{"cell": partial.cell, "next": partial.probeN, "wall": wallW})
		}
	}
	return m
}

func (m *Maze) Run(ctx context.Context) {
	var objs []*Object
	for _, row := range m.cells {
		for _, partial := range row {
			objs = append(objs,
				partial.cell,
				partial.probeN,
				partial.probeE,
				partial.probeS,
				partial.probeW,
				partial.wallN,
				partial.wallW,
			)
		}
	}
	RunAll(ctx, objs...)
}

func (m *Maze) String() string {
	var b bytes.Buffer
	for r, row := range m.cells {
		for _, partial := range row {
			b.WriteString("┼")
			if r == 0 || !partial.openN {
				b.WriteString("─")
			} else {
				b.WriteString(" ")
			}
		}
		b.WriteString("\n")
		for c, partial := range row {
			if c == 0 || !partial.openW {
				b.WriteString("│ ")
			} else {
				b.WriteString("  ")
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

// +-+-+-+-+
// |.|.|.|.|
// +-+-+-+-+
// |.|.|.|.|
// +-+-+-+-+
