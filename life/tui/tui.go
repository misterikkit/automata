package tui

import (
	tcell "github.com/gdamore/tcell/v2"
	"github.com/pkg/errors"

	"github.com/misterikkit/automata/life/game"
)

type Event int

const (
	Escape Event = iota
	Left
	Right
	Enter
)

// TUI is a UI for the game.
type TUI struct {
	s tcell.Screen
	t []rune
	h func(Event)
}

// New allocates a TUI and takes over the terminal screen.
func New(text string, handler func(Event)) (*TUI, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create screen")
	}
	if err := s.Init(); err != nil {
		return nil, errors.Wrapf(err, "unable to initialize screen")
	}
	t := &TUI{s, []rune(text), handler}
	// TODO: make this joinable
	go func(t *TUI) {
		for {
			e := t.s.PollEvent()
			ek, ok := e.(*tcell.EventKey)
			if !ok {
				continue
			}
			switch ek.Key() {
			case tcell.KeyEsc:
				t.h(Escape)
			case tcell.KeyLeft:
				t.h(Left)
			case tcell.KeyRight:
				t.h(Right)
			case tcell.KeyEnter:
				t.h(Enter)
			}
		}
	}(t)
	return t, nil
}

// Close releases resources and returns the terminal to normal.
func (t *TUI) Close() { t.s.Fini() }

func (t *TUI) DrawGame(g game.Game) {
	box(t.s, 0, 0, g.Cols()+3, g.Rows()+1)
	for r := 0; r < g.Rows(); r++ {
		for c := 0; c < g.Cols(); c++ {
			val := tcell.RuneBullet
			if g.Get(r, c) {
				val = tcell.RuneBlock
			}
			t.s.SetContent(c+2, r+1, val, nil, tcell.StyleDefault)
		}
	}
	for i, r := range t.t {
		t.s.SetContent(i, g.Rows()+2, r, nil, tcell.StyleDefault)
	}
	t.s.Sync()
}

func box(s tcell.Screen, x, y, w, h int) {
	for i := 0; i < w; i++ {
		s.SetContent(x+i, y, tcell.RuneHLine, nil, tcell.StyleDefault)
		s.SetContent(x+i, y+h, tcell.RuneHLine, nil, tcell.StyleDefault)
	}
	for j := 0; j < h; j++ {
		s.SetContent(x, y+j, tcell.RuneVLine, nil, tcell.StyleDefault)
		s.SetContent(x+w, y+j, tcell.RuneVLine, nil, tcell.StyleDefault)
	}
	s.SetContent(x, y, tcell.RuneULCorner, nil, tcell.StyleDefault)
	s.SetContent(x, y+h, tcell.RuneLLCorner, nil, tcell.StyleDefault)
	s.SetContent(x+w, y, tcell.RuneURCorner, nil, tcell.StyleDefault)
	s.SetContent(x+w, y+h, tcell.RuneLRCorner, nil, tcell.StyleDefault)
	// Don't forget to Show() or Sync()
}
