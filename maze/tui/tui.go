package tui

import (
	"bytes"

	"github.com/fatih/color"
	tcell "github.com/gdamore/tcell/v2"
	"github.com/pkg/errors"

	"github.com/misterikkit/automata/maze/game"
)

type Game interface {
	Rows() int
	Cols() int
	Get(row, col int) game.Cell
}

type MappedGame interface {
	Game
	GetTag(row, col int) int
}

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

func (t *TUI) DrawGame(g Game) {
	box(t.s, 0, 0, g.Cols()+3, g.Rows()+1)
	for r := 0; r < g.Rows(); r++ {
		for c := 0; c < g.Cols(); c++ {
			val := tcell.RuneBullet
			if g.Get(r, c) {
				val = tcell.RuneBlock
			}
			style := styleFor(g, r, c)
			t.s.SetContent(c+2, r+1, val, nil, style)
		}
	}
	for i, r := range t.t {
		t.s.SetContent(i, g.Rows()+2, r, nil, tcell.StyleDefault)
	}
	t.s.Sync()
}

func box(s tcell.Screen, x, y, w, h int) {
	style := tcell.StyleDefault //.Foreground(tcell.ColorAliceBlue).Background(tcell.ColorOrchid)
	for i := 0; i < w; i++ {
		s.SetContent(x+i, y, tcell.RuneHLine, nil, style)
		s.SetContent(x+i, y+h, tcell.RuneHLine, nil, style)
	}
	for j := 0; j < h; j++ {
		s.SetContent(x, y+j, tcell.RuneVLine, nil, style)
		s.SetContent(x+w, y+j, tcell.RuneVLine, nil, style)
	}
	s.SetContent(x, y, tcell.RuneULCorner, nil, style)
	s.SetContent(x, y+h, tcell.RuneLLCorner, nil, style)
	s.SetContent(x+w, y, tcell.RuneURCorner, nil, style)
	s.SetContent(x+w, y+h, tcell.RuneLRCorner, nil, style)
	// Don't forget to Show() or Sync()
}

func Fmt(g MappedGame) string {
	var b bytes.Buffer
	// top row
	b.WriteRune(tcell.RuneULCorner)
	for i := 0; i < g.Cols(); i++ {
		b.WriteRune(tcell.RuneHLine)
	}
	b.WriteRune(tcell.RuneURCorner)
	b.WriteString("\n")
	// content
	for r := 0; r < g.Rows(); r++ {
		b.WriteRune(tcell.RuneVLine)
		for c := 0; c < g.Cols(); c++ {
			switch g.Get(r, c) {
			case true:
				b.WriteRune(tcell.RuneBlock)
			case false:
				color := fatihColors[g.GetTag(r, c)%len(fatihColors)]
				b.WriteString(color.Sprint("·"))
			}
		}
		b.WriteRune(tcell.RuneVLine)
		b.WriteString("\n")
	}
	// bottom row
	b.WriteRune(tcell.RuneLLCorner)
	for i := 0; i < g.Cols(); i++ {
		b.WriteRune(tcell.RuneHLine)
	}
	b.WriteRune(tcell.RuneLRCorner)
	b.WriteString("\n")
	return b.String()
}

func styleFor(g Game, row, col int) tcell.Style {
	gm, ok := g.(MappedGame)
	if !ok {
		return tcell.StyleDefault
	}
	tag := gm.GetTag(row, col)
	if tag < 0 {
		return tcell.StyleDefault
	}
	color := tcellColors[tag%len(tcellColors)]
	return tcell.StyleDefault.Background(color)
}

// lazy copy pasta
var tcellColors = []tcell.Color{
	// tcell.ColorMaroon,
	// tcell.ColorGreen,
	// tcell.ColorOlive,
	// tcell.ColorNavy,
	// tcell.ColorPurple,
	// tcell.ColorTeal,
	// tcell.ColorSilver,
	// tcell.ColorGray,
	// tcell.ColorRed,
	// tcell.ColorLime,
	// tcell.ColorYellow,
	// tcell.ColorBlue,
	// tcell.ColorFuchsia,
	// tcell.ColorAqua,

	tcell.ColorBlack,
	tcell.ColorDarkRed,
	tcell.ColorDarkGreen,
	tcell.ColorDarkGoldenrod,
	tcell.ColorDarkBlue,
	tcell.ColorDarkMagenta,
	tcell.ColorDarkCyan,
	tcell.ColorRed,
	tcell.ColorGreen,
	tcell.ColorYellow,
	tcell.ColorBlue,
	tcell.ColorLightPink,
	tcell.ColorLightCyan,
}

var fatihColors = []*color.Color{

	color.New(color.BgBlack),
	color.New(color.BgRed),
	color.New(color.BgGreen),
	color.New(color.BgYellow),
	color.New(color.BgBlue),
	color.New(color.BgMagenta),
	color.New(color.BgCyan),
	color.New(color.BgHiRed),
	color.New(color.BgHiGreen),
	color.New(color.BgHiYellow),
	color.New(color.BgHiBlue),
	color.New(color.BgHiMagenta),
	color.New(color.BgHiCyan),
}
