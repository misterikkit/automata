package gene

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/misterikkit/automata/maze/game"
)

type Gene struct {
	Alive [9]game.Cell
	Dead  [9]game.Cell
}

func FromString(s string) Gene {
	rs := []rune(s)
	g := Gene{}
	for i, r := range rs {
		sl := g.Alive[0:9]
		if i >= 9 {
			sl = g.Dead[0:9]
		}
		switch r {
		case '0':
			sl[i%9] = false
		case '1':
			sl[i%9] = true
		default:
			panic(fmt.Sprintf("Illegal rune: %v", r))
		}
	}
	return g
}

func (g Gene) String() string {
	bits := make([]string, 18)
	for i := 0; i < 9; i++ {
		switch g.Alive[i] {
		case true:
			bits[i] = "1"
		case false:
			bits[i] = "0"
		}
	}
	for i := 0; i < 9; i++ {
		switch g.Dead[i] {
		case true:
			bits[9+i] = "1"
		case false:
			bits[9+i] = "0"
		}
	}
	return strings.Join(bits, "")
}

func (g Gene) AsRule() game.Rule {
	return func(game game.Game, row, col int) game.Cell {
		n := game.CountNeighbors(row, col)
		switch game.Get(row, col) {
		case true:
			return g.Alive[n]
		case false:
			return g.Dead[n]
		}
		panic("ran out of booleans")
	}
}

func Random() Gene {
	g := Gene{}
	for i := 0; i < 9; i++ {
		g.Alive[i] = rand.Int()%2 == 0
		g.Dead[i] = rand.Int()%2 == 0
	}
	return g
}

func Clone() game.Rule {
	return func(g game.Game, row, col int) game.Cell { return g.Get(row, col) }
}
